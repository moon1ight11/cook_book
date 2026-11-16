import { api } from '@/lib/api';
import type { SuccessResponse, UpdateUserRequest, User } from '@/types/api';

export async function getMe(): Promise<User> {
    const { data } = await api.get<User>('/private/users');
    return data;
}

export async function updateMe(payload: UpdateUserRequest): Promise<SuccessResponse> {
    const { data } = await api.patch<SuccessResponse>('/private/users', payload);
    return data;
}

export async function deleteMe(): Promise<SuccessResponse> {
    const { data } = await api.delete<SuccessResponse>('/private/users');
    return data;
}