'use client';

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { qk } from '@/lib/queryKeys';
import { normalizeError } from '@/lib/errors';
import { toast } from '@/lib/toast';
import { useAuthStore } from '@/lib/store';
import * as profileApi from './api';
import type { UpdateUserRequest } from '@/types/api';

export function useMe() {
    return useQuery({
        queryKey: qk.me,
        queryFn: profileApi.getMe,
        staleTime: 60_000,
    });
}

export function useUpdateMe() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (payload: UpdateUserRequest) => profileApi.updateMe(payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: qk.me });
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useDeleteMe() {
    const router = useRouter();
    const logoutLocal = useAuthStore((s) => s.logoutLocal);
    return useMutation({
        mutationFn: profileApi.deleteMe,
        onSuccess: () => {
            toast.success('Аккаунт удалён');
            logoutLocal();
            router.replace('/');
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}