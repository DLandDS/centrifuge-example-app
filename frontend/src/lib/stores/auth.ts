import { writable } from 'svelte/store';
import type { User } from '$lib/types';

interface AuthState {
    isAuthenticated: boolean;
    user: User | null;
    token: string | null;
}

const initialState: AuthState = {
    isAuthenticated: false,
    user: null,
    token: null
};

export const authStore = writable<AuthState>(initialState);

export const login = (token: string, user: User) => {
    authStore.set({
        isAuthenticated: true,
        user,
        token
    });
    
    // Store in localStorage for persistence
    localStorage.setItem('auth_token', token);
    localStorage.setItem('user', JSON.stringify(user));
};

export const logout = () => {
    authStore.set(initialState);
    localStorage.removeItem('auth_token');
    localStorage.removeItem('user');
};

export const initializeAuth = () => {
    if (typeof window !== 'undefined') {
        const token = localStorage.getItem('auth_token');
        const userJson = localStorage.getItem('user');
        
        if (token && userJson) {
            try {
                const user = JSON.parse(userJson);
                authStore.set({
                    isAuthenticated: true,
                    user,
                    token
                });
            } catch (e) {
                logout();
            }
        }
    }
};