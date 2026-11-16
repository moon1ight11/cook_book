import { api } from '@/lib/api';
import type {
    CreateRecipeRequest,
    DetailedRecipe,
    ListRecipesParams,
    Paginated,
    Recipe,
    SuccessResponse,
    UpdateRecipeRequest,
} from '@/types/api';

export async function listRecipes(params: ListRecipesParams): Promise<Paginated<Recipe>> {
    const { data } = await api.get<Paginated<Recipe>>('/public/recipes', { params });
    return { ...data, data: data.data ?? [] };
}

export async function getRecipe(id: string): Promise<DetailedRecipe> {
    const { data } = await api.get<DetailedRecipe>(`/public/recipes/${id}`);
    return data;
}

export async function recipesByTags(tags: string[]): Promise<Recipe[]> {
    const { data } = await api.get<Recipe[] | null>('/public/recipes-by-tags', {
        params: { tags },
        paramsSerializer: { indexes: null },
    });
    return data ?? [];
}

export async function recipesByUser(userId: string): Promise<Recipe[]> {
    const { data } = await api.get<Recipe[] | null>(`/public/recipes-by-user/${userId}`);
    return data ?? [];
}

export async function createRecipe(payload: CreateRecipeRequest): Promise<Recipe> {
    const { data } = await api.post<Recipe>('/private/recipes', payload);
    return data;
}

export async function updateRecipe(
    id: string,
    payload: UpdateRecipeRequest
): Promise<SuccessResponse> {
    const { data } = await api.patch<SuccessResponse>(`/private/recipes/${id}`, payload);
    return data;
}

export async function deleteRecipe(id: string): Promise<SuccessResponse> {
    const { data } = await api.delete<SuccessResponse>(`/private/recipes/${id}`);
    return data;
}