// WebSocket types matching backend internal/websocket/types.go

// Session roles
export type SessionRole = 'agent' | 'prospect';

// Stream types for dual-stream diarization
// Matches backend internal/websocket/types.go StreamType
export type StreamType = 'agent' | 'mixed';

// Stream type byte values (for binary protocol)
export const STREAM_TYPE_AGENT = 0;
export const STREAM_TYPE_MIXED = 1;

// Client message types (from client to server)
export type ClientMessageType =
	| 'audio_chunk'
	| 'transcript'
	| 'agent_note'
	| 'generate_report'
	| 'heartbeat';

// Server message types (from server to client)
export type ServerMessageType =
	| 'transcript'
	| 'blocker'
	| 'suggestion'
	| 'analysis'
	| 'phase_change'
	| 'qualification_update'
	| 'pong'
	| 'error';

// Client message structure
export interface ClientMessage {
	type: ClientMessageType;
	payload: unknown;
}

// Server message structure
export interface ServerMessage {
	type: ServerMessageType;
	session_id: string;
	payload: unknown;
	timestamp: string;
}

// Audio chunk payload
export interface AudioChunkPayload {
	data: string; // Base64 encoded audio
	format: string; // webm, opus, etc.
	sample_rate: number;
	sequence: number;
}

// Transcript payload
export interface TranscriptPayload {
	speaker: string;
	text: string;
	confidence: number;
	start_time: number;
	end_time: number;
	is_final: boolean;
	timestamp: number; // Unix milliseconds
}

// Blocker payload
export interface BlockerPayload {
	id: string;
	type: string;
	severity: 'low' | 'medium' | 'high' | 'critical';
	description: string;
	evidence: string;
	solution?: string;
	recommended_solution?: string;
	estimated_timeline?: string;
	estimated_investment?: string;
	confidence?: number;
	timestamp: string;
}

// Suggestion payload
export interface SuggestionPayload {
	id: string;
	suggestion_type: string;
	type: string;
	title: string;
	description: string;
	content: string;
	script?: string;
	priority: 'low' | 'medium' | 'high';
	related_blocker_id?: string;
	for_agent: boolean;
	context: string;
	timestamp: string;
}

// Phase change payload
export interface PhaseChangePayload {
	current_phase: string;
	phase_start_time: number;
	time_in_phase: number;
	total_call_duration: number;
	phases_completed: string[];
	next_phase_recommendation?: string;
}

// Qualification payload
export interface QualificationPayload {
	budget_authority?: boolean;
	timeline_urgency?: boolean;
	decision_power?: boolean;
	notes?: string;
}

// Agent note payload
export interface AgentNotePayload {
	note: string;
	timestamp: string;
}

// Generate report payload
export interface GenerateReportPayload {
	format: string; // pdf, docx, etc.
}

// Heartbeat payload
export interface HeartbeatPayload {
	timestamp: string;
}

// Error payload
export interface ErrorPayload {
	code: string;
	message: string;
}

// Connection options
export interface CallWSOptions {
	sessionId: string;
	role?: SessionRole;
	token?: string;
	onTranscript?: (payload: TranscriptPayload) => void;
	onBlocker?: (payload: BlockerPayload) => void;
	onSuggestion?: (payload: SuggestionPayload) => void;
	onPhaseChange?: (payload: PhaseChangePayload) => void;
	onQualificationUpdate?: (payload: QualificationPayload) => void;
	onError?: (payload: ErrorPayload) => void;
	onConnected?: () => void;
	onDisconnected?: () => void;
	onReconnecting?: (attempt: number, maxAttempts: number) => void;
}
