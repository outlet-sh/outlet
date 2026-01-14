// WebSocket client for call sessions
import type {
	CallWSOptions,
	ClientMessage,
	ServerMessage,
	AudioChunkPayload,
	TranscriptPayload,
	BlockerPayload,
	SuggestionPayload,
	PhaseChangePayload,
	QualificationPayload,
	ErrorPayload,
	AgentNotePayload,
	StreamType
} from './types';
import { STREAM_TYPE_AGENT, STREAM_TYPE_MIXED } from './types';

export class CallWebSocket {
	private ws: WebSocket | null = null;
	private options: CallWSOptions;
	private reconnectAttempts = 0;
	private maxReconnectAttempts = Infinity; // Never give up - keep trying forever
	private baseReconnectDelay = 1000;
	private maxReconnectDelay = 30000;
	private heartbeatInterval: ReturnType<typeof setInterval> | null = null;
	private reconnectTimeout: ReturnType<typeof setTimeout> | null = null;
	private audioSequence = 0;
	private messageBuffer: Array<{ type: string; payload: unknown }> = [];
	private isReconnecting = false;
	private intentionalClose = false;

	constructor(options: CallWSOptions) {
		this.options = {
			role: 'agent',
			...options
		};
	}

	// Connect to the WebSocket server
	connect(): void {
		// Clear any pending reconnect
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}

		// Don't connect if already connected or connecting
		if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
			console.log('WebSocket already connected or connecting');
			return;
		}

		this.intentionalClose = false;
		const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const token = this.options.token || localStorage.getItem('auth_token') || 'anonymous';
		const wsUrl = `${protocol}//${window.location.host}/ws/call?session_id=${this.options.sessionId}&token=${token}&role=${this.options.role}`;

		console.log('Connecting to WebSocket:', wsUrl);
		this.ws = new WebSocket(wsUrl);

		this.ws.onopen = () => {
			console.log('WebSocket connected');
			this.reconnectAttempts = 0;
			this.isReconnecting = false;
			this.startHeartbeat();
			this.flushMessageBuffer();
			this.options.onConnected?.();
		};

		this.ws.onclose = (event) => {
			console.log('WebSocket disconnected:', event.code, event.reason);
			this.stopHeartbeat();
			this.options.onDisconnected?.();

			// Auto-reconnect if not intentional close - keep trying forever
			if (!this.intentionalClose) {
				this.scheduleReconnect();
			}
		};

		this.ws.onerror = (error) => {
			console.error('WebSocket error:', error);
			this.options.onError?.({ code: 'WS_ERROR', message: 'WebSocket connection error' });
		};

		this.ws.onmessage = (event) => {
			this.handleMessage(event.data);
		};
	}

	// Disconnect from the WebSocket server
	disconnect(): void {
		this.intentionalClose = true;
		this.stopHeartbeat();

		// Clear any pending reconnect
		if (this.reconnectTimeout) {
			clearTimeout(this.reconnectTimeout);
			this.reconnectTimeout = null;
		}

		if (this.ws) {
			this.ws.close(1000, 'Client disconnect');
			this.ws = null;
		}

		// Clear message buffer on intentional disconnect
		this.messageBuffer = [];
		this.reconnectAttempts = 0;
		this.isReconnecting = false;
	}

	// Check if connected
	isConnected(): boolean {
		return this.ws?.readyState === WebSocket.OPEN;
	}

	// Send audio chunk as binary WebSocket frame (raw PCM16 data)
	// DEPRECATED: Use sendAudioFrame with stream type for dual-stream diarization
	sendAudioChunk(data: Int16Array): void {
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
			// Don't buffer audio during reconnect (too much data)
			return;
		}

		try {
			// Send raw binary data directly - backend expects PCM16 @ 16kHz mono
			this.ws.send(data.buffer);
			this.audioSequence++;
		} catch (error) {
			console.error('Failed to send audio chunk:', error);
		}
	}

	// Send audio frame with stream type header for dual-stream diarization
	// Binary format: [1 byte stream_type] + [PCM16 audio data]
	// stream_type: 0 = agent mic, 1 = mixed Meet tab audio
	sendAudioFrame(streamType: StreamType, data: Int16Array): void {
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
			// Don't buffer audio during reconnect (too much data)
			return;
		}

		try {
			// Create combined buffer: [1 byte header] + [PCM16 data]
			const audioBytes = new Uint8Array(data.buffer);
			const combined = new Uint8Array(1 + audioBytes.byteLength);

			// Set stream type header byte
			combined[0] = streamType === 'agent' ? STREAM_TYPE_AGENT : STREAM_TYPE_MIXED;

			// Copy audio data after header
			combined.set(audioBytes, 1);

			// Send as binary WebSocket frame
			this.ws.send(combined.buffer);
			this.audioSequence++;
		} catch (error) {
			console.error('Failed to send audio frame:', error);
		}
	}

	// Send audio chunk as JSON (for compatibility/testing)
	sendAudioChunkJSON(data: Int16Array, format = 'pcm16', sampleRate = 16000): void {
		const payload: AudioChunkPayload = {
			data: Array.from(data).join(','), // Send as comma-separated values
			format,
			sample_rate: sampleRate,
			sequence: this.audioSequence++
		};
		this.send('audio_chunk', payload);
	}

	// Send transcript (for manual input or correction)
	sendTranscript(speaker: string, text: string, isFinal = true): void {
		const payload: TranscriptPayload = {
			speaker,
			text,
			confidence: 1.0,
			start_time: 0,
			end_time: 0,
			is_final: isFinal,
			timestamp: Date.now()
		};
		this.send('transcript', payload);
	}

	// Send agent note
	sendAgentNote(note: string): void {
		const payload: AgentNotePayload = {
			note,
			timestamp: new Date().toISOString()
		};
		this.send('agent_note', payload);
	}

	// Request report generation
	requestReport(format = 'pdf'): void {
		this.send('generate_report', { format });
	}

	// Send heartbeat
	private sendHeartbeat(): void {
		this.send('heartbeat', { timestamp: new Date().toISOString() });
	}

	// Generic send method
	private send(type: string, payload: unknown): void {
		const message: ClientMessage = {
			type: type as ClientMessage['type'],
			payload
		};

		// Buffer messages during reconnection (except heartbeat)
		if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
			if (this.isReconnecting && type !== 'heartbeat') {
				// Buffer important messages during reconnect, limit buffer size
				if (this.messageBuffer.length < 100) {
					this.messageBuffer.push({ type, payload });
					console.log(`Buffered message during reconnect: ${type}`);
				} else {
					console.warn('Message buffer full, dropping message');
				}
			} else {
				console.warn('WebSocket not connected, cannot send message');
			}
			return;
		}

		try {
			this.ws.send(JSON.stringify(message));
		} catch (error) {
			console.error('Failed to send WebSocket message:', error);
			// Buffer the message if send fails
			if (this.messageBuffer.length < 100) {
				this.messageBuffer.push({ type, payload });
			}
		}
	}

	// Flush buffered messages after reconnect
	private flushMessageBuffer(): void {
		if (this.messageBuffer.length === 0) return;

		console.log(`Flushing ${this.messageBuffer.length} buffered messages`);
		const buffer = [...this.messageBuffer];
		this.messageBuffer = [];

		for (const { type, payload } of buffer) {
			this.send(type, payload);
		}
	}

	// Handle incoming messages
	private handleMessage(data: string): void {
		try {
			const message: ServerMessage = JSON.parse(data);

			switch (message.type) {
				case 'transcript':
					this.options.onTranscript?.(message.payload as TranscriptPayload);
					break;

				case 'blocker':
					this.options.onBlocker?.(message.payload as BlockerPayload);
					break;

				case 'suggestion':
					this.options.onSuggestion?.(message.payload as SuggestionPayload);
					break;

				case 'phase_change':
					this.options.onPhaseChange?.(message.payload as PhaseChangePayload);
					break;

				case 'qualification_update':
					this.options.onQualificationUpdate?.(message.payload as QualificationPayload);
					break;

				case 'error':
					this.options.onError?.(message.payload as ErrorPayload);
					break;

				case 'pong':
					// Heartbeat response, no action needed
					break;

				default:
					console.log('Unknown message type:', message.type);
			}
		} catch (error) {
			console.error('Failed to parse WebSocket message:', error);
		}
	}

	// Start heartbeat interval
	private startHeartbeat(): void {
		this.heartbeatInterval = setInterval(() => {
			this.sendHeartbeat();
		}, 30000); // Every 30 seconds
	}

	// Stop heartbeat interval
	private stopHeartbeat(): void {
		if (this.heartbeatInterval) {
			clearInterval(this.heartbeatInterval);
			this.heartbeatInterval = null;
		}
	}

	// Schedule reconnection attempt with exponential backoff and jitter
	private scheduleReconnect(): void {
		this.reconnectAttempts++;
		this.isReconnecting = true;

		// Exponential backoff with jitter: base * 2^attempt + random(0-1000ms)
		const exponentialDelay = this.baseReconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
		const jitter = Math.random() * 1000; // Random 0-1000ms to prevent thundering herd
		const delay = Math.min(exponentialDelay + jitter, this.maxReconnectDelay);

		console.log(`Scheduling reconnect in ${Math.round(delay)}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);

		// Notify listeners that we're reconnecting
		this.options.onReconnecting?.(this.reconnectAttempts, this.maxReconnectAttempts);

		this.reconnectTimeout = setTimeout(() => {
			this.reconnectTimeout = null;
			console.log('Attempting to reconnect...');
			this.connect();
		}, delay);
	}

	// Get current connection state
	getState(): 'disconnected' | 'connecting' | 'connected' | 'reconnecting' {
		if (this.isReconnecting) return 'reconnecting';
		if (!this.ws) return 'disconnected';
		switch (this.ws.readyState) {
			case WebSocket.CONNECTING:
				return 'connecting';
			case WebSocket.OPEN:
				return 'connected';
			default:
				return 'disconnected';
		}
	}

	// Get reconnect info
	getReconnectInfo(): { attempts: number; maxAttempts: number; isReconnecting: boolean } {
		return {
			attempts: this.reconnectAttempts,
			maxAttempts: this.maxReconnectAttempts,
			isReconnecting: this.isReconnecting
		};
	}
}

// Factory function for convenience
export function createCallWebSocket(options: CallWSOptions): CallWebSocket {
	return new CallWebSocket(options);
}
