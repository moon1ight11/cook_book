import { create } from 'zustand';

export type ToastKind = 'success' | 'error' | 'info';

export interface ToastItem {
    id: string;
    kind: ToastKind;
    message: string;
}

interface ToastState {
    items: ToastItem[];
    push: (kind: ToastKind, message: string) => void;
    remove: (id: string) => void;
}

export const useToastStore = create<ToastState>((set) => ({
    items: [],
    push: (kind, message) => {
        const id = Math.random().toString(36).slice(2);
        set((s) => ({ items: [...s.items, { id, kind, message }] }));
        setTimeout(() => {
            set((s) => ({ items: s.items.filter((t) => t.id !== id) }));
        }, 4000);
    },
    remove: (id) => set((s) => ({ items: s.items.filter((t) => t.id !== id) })),
}));

export const toast = {
    success: (m: string) => useToastStore.getState().push('success', m),
    error: (m: string) => useToastStore.getState().push('error', m),
    info: (m: string) => useToastStore.getState().push('info', m),
};