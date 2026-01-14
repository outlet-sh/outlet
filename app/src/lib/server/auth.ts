import { createHash, randomBytes, pbkdf2 } from 'crypto';
import { promisify } from 'util';

const pbkdf2Async = promisify(pbkdf2);

// Password hashing configuration
const ITERATIONS = 100000;
const KEY_LENGTH = 64;
const DIGEST = 'sha512';

/**
 * Hash a password using PBKDF2
 */
export async function hashPassword(password: string): Promise<string> {
	const salt = randomBytes(16).toString('hex');
	const hash = await pbkdf2Async(password, salt, ITERATIONS, KEY_LENGTH, DIGEST);
	return `${salt}:${hash.toString('hex')}`;
}

/**
 * Verify a password against a hash
 */
export async function verifyPassword(password: string, hashedPassword: string): Promise<boolean> {
	const [salt, hash] = hashedPassword.split(':');
	const verifyHash = await pbkdf2Async(password, salt, ITERATIONS, KEY_LENGTH, DIGEST);
	return hash === verifyHash.toString('hex');
}

/**
 * Generate a random temporary password
 */
export function generateTemporaryPassword(length: number = 12): string {
	const chars = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*';
	let password = '';
	const randomValues = randomBytes(length);

	for (let i = 0; i < length; i++) {
		password += chars[randomValues[i] % chars.length];
	}

	return password;
}

/**
 * Check if user has required role
 */
export function hasRole(userRole: string, requiredRole: string | string[]): boolean {
	const roles = Array.isArray(requiredRole) ? requiredRole : [requiredRole];

	// Role hierarchy: admin > manager > agent
	const roleHierarchy: Record<string, number> = {
		admin: 3,
		manager: 2,
		agent: 1
	};

	const userRoleLevel = roleHierarchy[userRole] || 0;
	const requiredRoleLevel = Math.min(...roles.map(r => roleHierarchy[r] || 0));

	return userRoleLevel >= requiredRoleLevel;
}
