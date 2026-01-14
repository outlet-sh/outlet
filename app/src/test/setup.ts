import { vi } from 'vitest';
import '@testing-library/jest-dom/vitest';

// Mock WebSocket for testing
global.WebSocket = vi.fn(() => ({
    send: vi.fn(),
    close: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    readyState: 1,
    CONNECTING: 0,
    OPEN: 1,
    CLOSING: 2,
    CLOSED: 3
})) as any;

// Mock window.dispatchEvent
Object.defineProperty(window, 'dispatchEvent', {
    value: vi.fn(),
    writable: true
});

// Mock window.addEventListener
Object.defineProperty(window, 'addEventListener', {
    value: vi.fn(),
    writable: true
});

// Mock localStorage
Object.defineProperty(window, 'localStorage', {
    value: {
        getItem: vi.fn(),
        setItem: vi.fn(),
        removeItem: vi.fn(),
        clear: vi.fn()
    },
    writable: true
});
