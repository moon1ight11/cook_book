import type { ListRecipesParams } from '@/types/api';

export const qk = {
    me: ['me'] as const,

    recipes: ['recipes'] as const,
    recipesList: (params: ListRecipesParams) => ['recipes', 'list', params] as const,
    recipe: (id: string) => ['recipes', 'detail', id] as const,
    recipesByUser: (userId: string) => ['recipes', 'byUser', userId] as const,
    recipesByTags: (tags: string[]) => ['recipes', 'byTags', tags] as const,
};