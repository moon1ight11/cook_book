'use client';

import { useMutation, useQueryClient } from '@tanstack/react-query';
import { qk } from '@/lib/queryKeys';
import { normalizeError } from '@/lib/errors';
import { toast } from '@/lib/toast';
import * as adminApi from './api';
import type { CreateAdminRequest, Role } from '@/types/api';

export function useCreateAdmin() {
    return useMutation({
        mutationFn: (payload: CreateAdminRequest) => adminApi.createAdmin(payload),
        onSuccess: () => toast.success('Администратор создан'),
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useUpdateUserRole() {
    return useMutation({
        mutationFn: ({ userId, role }: { userId: string; role: Role }) =>
            adminApi.updateUserRole(userId, role),
        onSuccess: () => toast.success('Роль обновлена'),
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useDeleteAnyUser() {
    return useMutation({
        mutationFn: (userId: string) => adminApi.deleteAnyUser(userId),
        onSuccess: () => toast.success('Пользователь удалён'),
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useDeleteAnyRecipe() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (recipeId: string) => adminApi.deleteAnyRecipe(recipeId),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: qk.recipes });
            toast.success('Рецепт удалён');
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}