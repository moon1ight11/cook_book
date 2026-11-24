'use client';

import {
    useMutation,
    useQuery,
    useQueryClient,
    keepPreviousData,
} from '@tanstack/react-query';
import { qk } from '@/lib/queryKeys';
import { normalizeError } from '@/lib/errors';
import { toast } from '@/lib/toast';
import * as recipesApi from './api';
import type {
    CreateRecipeRequest,
    ListRecipesParams,
    UpdateRecipeRequest,
} from '@/types/api';

export function useRecipes(params: ListRecipesParams) {
    return useQuery({
        queryKey: qk.recipesList(params),
        queryFn: () => recipesApi.listRecipes(params),
        placeholderData: keepPreviousData, // не мигаем при смене страницы
    });
}

export function useRecipe(id: string | undefined) {
    return useQuery({
        queryKey: qk.recipe(id ?? ''),
        queryFn: () => recipesApi.getRecipe(id as string),
        enabled: !!id, 
        retry: false, 
    });
}

export function useRecipesByTags(tags: string[], enabled = true) {
    return useQuery({
        queryKey: qk.recipesByTags(tags),
        queryFn: () => recipesApi.recipesByTags(tags),
        enabled: enabled && tags.length > 0,
    });
}

export function useRecipesByUser(userId: string | undefined) {
    return useQuery({
        queryKey: qk.recipesByUser(userId ?? ''),
        queryFn: () => recipesApi.recipesByUser(userId as string),
        enabled: !!userId,
    });
}

export function useCreateRecipe() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (payload: CreateRecipeRequest) => recipesApi.createRecipe(payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: qk.recipes });
            toast.success('Рецепт опубликован');
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useUpdateRecipe(id: string) {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (payload: UpdateRecipeRequest) => recipesApi.updateRecipe(id, payload),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: qk.recipes });
            toast.success('Рецепт обновлён');
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}

export function useDeleteRecipe() {
    const qc = useQueryClient();
    return useMutation({
        mutationFn: (id: string) => recipesApi.deleteRecipe(id),
        onSuccess: () => {
            qc.invalidateQueries({ queryKey: qk.recipes });
            toast.success('Рецепт удалён');
        },
        onError: (err) => toast.error(normalizeError(err)),
    });
}