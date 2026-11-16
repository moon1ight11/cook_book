import type { Difficulty } from '@/types/api';

export const DIFFICULTY_LABEL: Record<Difficulty, string> = {
    easy: 'Легко',
    medium: 'Средне',
    hard: 'Сложно',
};

export function formatMinutes(total: number): string {
    if (!total || total <= 0) return '—';
    const h = Math.floor(total / 60);
    const m = total % 60;
    if (h && m) return `${h} ч ${m} мин`;
    if (h) return `${h} ч`;
    return `${m} мин`;
}

export function formatDate(iso: string): string {
    try {
        return new Date(iso).toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: 'long',
            year: 'numeric',
        });
    } catch {
        return '';
    }
}

const TAG_EMOJI: Record<string, string> = {
    суп: '🍲',
    десерт: '🍰',
    мясо: '🥩',
    салат: '🥗',
    выпечка: '🥐',
    напиток: '🥤',
    рыба: '🐟',
    завтрак: '🍳',
    паста: '🍝',
    овощи: '🥦',
};

export function tagEmoji(tags: string[]): string {
    for (const t of tags) {
        const key = t.toLowerCase().trim();
        if (TAG_EMOJI[key]) return TAG_EMOJI[key];
    }
    return '🍽️';
}

export function normalizeTag(tag: string): string {
    return tag.toLowerCase().trim();
}