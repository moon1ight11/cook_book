import { create } from 'zustand';
import type { Role } from '@/types/api';

export interface AuthUser {
    id: string;
    name: string;
    email: string;
    role: Role;
}

export type AuthStatus = 'idle' | 'loading' | 'authenticated' | 'guest';

interface AuthState {
    user: AuthUser | null;
    status: AuthStatus;
    setUser: (u: AuthUser | null) => void;
    setStatus: (s: AuthStatus) => void;
    logoutLocal: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
    user: null,
    status: 'idle',
    setUser: (user) => set({ user }),
    setStatus: (status) => set({ status }),
    logoutLocal: () => {
        localStorage.removeItem('refresh_token');
        set({ user: null, status: 'guest' });
    },
}));