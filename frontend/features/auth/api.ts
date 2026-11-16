import { api } from '@/lib/api';
import type {
    RefreshTokenResponse,
    SignInRequest,
    SignUpRequest,
    SuccessResponse,
    User,
} from '@/types/api';

export async function signIn(payload: SignInRequest): Promise<RefreshTokenResponse> {
    const { data } = await api.post<RefreshTokenResponse>('/auth/sign-in', payload);
    return data;
}

export async function signUp(payload: SignUpRequest): Promise<RefreshTokenResponse> {
    const { data } = await api.post<RefreshTokenResponse>('/auth/sign-up', payload);
    return data;
}

export async function signOut(): Promise<SuccessResponse> {
    const { data } = await api.post<SuccessResponse>('/private/sign-out');
    return data;
}

export async function fetchMe(): Promise<User> {
    const { data } = await api.get<User>('/private/users');
    return data;
}