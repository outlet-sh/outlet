/**
 * Auth store for reactive authentication state
 * Wraps the auth utilities with a Svelte 5 reactive store
 */

import { browser } from '$app/environment';
import { goto } from '$app/navigation';
import {
	getToken,
	setToken,
	removeToken,
	getCurrentUser,
	isAuthenticated,
	isAdmin,
	decodeToken
} from '$lib/auth';
import { login as apiLogin, logout as apiLogout } from '$lib/api';

interface User {
	id: string;
	email: string;
	name?: string;
	role: string;
	exp: number;
}

interface AuthState {
	user: User | null;
	isAuthenticated: boolean;
	isAdmin: boolean;
	isLoading: boolean;
	error: string | null;
}

function createAuthStore() {
	let state = $state<AuthState>({
		user: null,
		isAuthenticated: false,
		isAdmin: false,
		isLoading: true,
		error: null
	});

	// Initialize from localStorage on creation
	if (browser) {
		const user = getCurrentUser();
		state = {
			user,
			isAuthenticated: isAuthenticated(),
			isAdmin: isAdmin(),
			isLoading: false,
			error: null
		};
	}

	return {
		get user() {
			return state.user;
		},
		get isAuthenticated() {
			return state.isAuthenticated;
		},
		get isAdmin() {
			return state.isAdmin;
		},
		get isLoading() {
			return state.isLoading;
		},
		get error() {
			return state.error;
		},

		async login(email: string, password: string) {
			state.isLoading = true;
			state.error = null;

			try {
				const response = await apiLogin({ email, password });
				if (response.token) {
					setToken(response.token);
					const user = decodeToken(response.token);
					state = {
						user,
						isAuthenticated: true,
						isAdmin: user?.role === 'admin' || user?.role === 'super_admin',
						isLoading: false,
						error: null
					};
					return true;
				}
				state.error = 'Login failed';
				state.isLoading = false;
				return false;
			} catch (err: any) {
				state.error = err.message || 'Login failed';
				state.isLoading = false;
				return false;
			}
		},

		async logout() {
			try {
				await apiLogout();
			} catch (err) {
				// Ignore errors on logout
			}
			removeToken();
			state = {
				user: null,
				isAuthenticated: false,
				isAdmin: false,
				isLoading: false,
				error: null
			};
			if (browser) {
				goto('/auth/login');
			}
		},

		async refresh() {
			const token = getToken();
			if (!token) {
				state.isAuthenticated = false;
				state.user = null;
				return false;
			}

			// Verify the token is still valid by decoding it
			const user = decodeToken(token);
			if (user && user.exp && user.exp * 1000 > Date.now()) {
				state = {
					user,
					isAuthenticated: true,
					isAdmin: user?.role === 'admin' || user?.role === 'super_admin',
					isLoading: false,
					error: null
				};
				return true;
			}

			// Token expired
			removeToken();
			state.isAuthenticated = false;
			state.user = null;
			return false;
		},

		initialize() {
			if (!browser) return;
			const user = getCurrentUser();
			state = {
				user,
				isAuthenticated: isAuthenticated(),
				isAdmin: isAdmin(),
				isLoading: false,
				error: null
			};
		},

		/**
		 * Set session directly from an API response
		 * Use this when you've already called the login API and have the response
		 */
		setSession(token: string, userInfo: { id: string; email: string; name?: string; role: string }) {
			setToken(token);
			const user = decodeToken(token) || {
				id: userInfo.id,
				email: userInfo.email,
				name: userInfo.name,
				role: userInfo.role,
				exp: Math.floor(Date.now() / 1000) + 3600 // Default 1 hour if decode fails
			};
			state = {
				user,
				isAuthenticated: true,
				isAdmin: user?.role === 'admin' || user?.role === 'super_admin',
				isLoading: false,
				error: null
			};
		}
	};
}

export const authStore = createAuthStore();
