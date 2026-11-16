'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter, useSearchParams } from 'next/navigation';
import { qk } from '@/lib/queryKeys';
import { useAuthStore } from '@/lib/store';
import { normalizeError } from '@/lib/errors';
import { toast } from '@/lib/toast';
import * as authApi from './api';
import type { SignInRequest, SignUpRequest } from '@/types/api';

export function useMe(enabled = true) {
    const setUser = useAuthStore((s) => s.setUser);
    const setStatus = useAuthStore((s) => s.setStatus);

    return useQuery({
        queryKey: qk.me,
        queryFn: authApi.fetchMe,
        enabled,
        retry: false,
        staleTime: 60_000,
        select: (data) => {
            setUser({ id: data.id, name: data.name, email: data.email, role: data.role });
            setStatus('authenticated');
            return data;
        },
    });
}

export function useSignIn() {
    const qc = useQueryClient();
    const router = useRouter();
    const searchParams = useSearchParams();
    const setUser = useAuthStore((s) => s.setUser);
    const setStatus = useAuthStore((s) => s.setStatus);

    return useMutation({
        mutationFn: (payload: SignInRequest) => authApi.signIn(payload),
        onSuccess: async (data) => {
            localStorage.setItem('refresh_token', data.refresh_token);
            try {
                const me = await authApi.fetchMe();
                setUser({ id: me.id, name: me.name, email: me.email, role: me.role });
                setStatus('authenticated');
                qc.setQueryData(qk.me, me);
                toast.success(`С возвращением, ${me.name}!`);
                const next = searchParams.get('next');
                router.replace(next ? decodeURIComponent(next) : '/');
            } catch {
                toast.error('Не удалось загрузить профиль');
                router.replace('/');
            }
        },
        onError: (err) => {
            toast.error(normalizeError(err));
        },
    });
}

export function useSignUp() {
    const qc = useQueryClient();
    const router = useRouter();
    const setUser = useAuthStore((s) => s.setUser);
    const setStatus = useAuthStore((s) => s.setStatus);

    return useMutation({
        mutationFn: (payload: SignUpRequest) => authApi.signUp(payload),
        onSuccess: async (data) => {
            localStorage.setItem('refresh_token', data.refresh_token);
            try {
                const me = await authApi.fetchMe();
                setUser({ id: me.id, name: me.name, email: me.email, role: me.role });
                setStatus('authenticated');
                qc.setQueryData(qk.me, me);
                toast.success('Аккаунт создан!');
                router.replace('/');
            } catch {
                toast.success('Аккаунт создан, войдите');
                router.replace('/login');
            }
        },
        onError: (err) => {
            toast.error(normalizeError(err));
        },
    });
}

export function useSignOut() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: authApi.signOut,
        onSettled: () => {
            qc.clear();
        },
    });
}