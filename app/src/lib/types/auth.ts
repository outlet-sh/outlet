/**
 * Authentication and Authorization Types
 */

export type UserRole = 'admin' | 'manager' | 'agent';

export interface User {
	id: string;
	email: string;
	password_hash: string;
	role: UserRole;
	first_name: string;
	last_name: string;
	created_at: Date;
	updated_at: Date;
	last_login?: Date;
	is_active: boolean;
}

export interface Session {
	id: string;
	user_id: string;
	token: string;
	expires_at: Date;
	created_at: Date;
	ip_address?: string;
	user_agent?: string;
}

export interface SessionWithUser extends Session {
	user: User;
}

export interface LoginCredentials {
	email: string;
	password: string;
}

export interface AuthResult {
	success: boolean;
	user?: User;
	session?: Session;
	error?: string;
}

export interface SessionData {
	userId: string;
	email: string;
	role: UserRole;
	sessionId: string;
}
