import { AxiosError } from 'axios';
import type { ErrorResponse } from '@/types/api';

const HUMAN: Record<string, string> = {
    'bad request': 'Некорректные данные',
    unauthorized: 'Требуется вход в аккаунт',
    'internal server error': 'Ошибка на сервере, попробуйте позже',
    'name, email and password are required': 'Заполните имя, email и пароль',
    'access violation': 'Недостаточно прав',
    'missing access token': 'Требуется вход в аккаунт',
};

export function normalizeError(err: unknown): string {
    if (err instanceof AxiosError) {
        if (!err.response) return 'Нет соединения с сервером';
        const data = err.response.data as ErrorResponse | undefined;
        if (data?.error) return HUMAN[data.error] ?? data.error;
        if (err.response.status === 401) return 'Требуется вход в аккаунт';
        if (err.response.status === 403) return 'Недостаточно прав';
        return `Ошибка ${err.response.status}`;
    }
    if (err instanceof Error) return err.message;
    return 'Что-то пошло не так';
}