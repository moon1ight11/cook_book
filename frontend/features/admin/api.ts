import { api } from '@/lib/api';
import type {
    CreateAdminRequest,
    Role,
    SuccessResponse,
} from '@/types/api';

export async function createAdmin(payload: CreateAdminRequest): Promise<SuccessResponse> {
    const { data } = await api.post<SuccessResponse>('/admin/users', payload);
    return data;
}

export async function updateUserRole(userId: string, role: Role): Promise<SuccessResponse> {
    const { data } = await api.patch<SuccessResponse>(`/admin/users/${userId}`, { role });
    return data;
}

export async function deleteAnyUser(userId: string): Promise<SuccessResponse> {
    const { data } = await api.delete<SuccessResponse>(`/admin/users/${userId}`);
    return data;
}

export async function deleteAnyRecipe(recipeId: string): Promise<SuccessResponse> {
    const { data } = await api.delete<SuccessResponse>(`/admin/recipe/${recipeId}`);
    return data;
}