'use client';

import { useRecipesByUser } from '@/features/recipes/hooks';
import { RecipeCard } from '@/features/recipes/RecipeCard';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/ui/EmptyState';
import { normalizeError } from '@/lib/errors';

export default function UserRecipesClient({ id }: { id: string }) {
    const { data, isLoading, isError, error } = useRecipesByUser(id);
    const items = data ?? [];

    return (
        <div className="flex flex-col gap-6">
            <h1 className="font-serif text-3xl text-ink">Рецепты пользователя</h1>
            {isLoading ? (
                <Spinner />
            ) : isError ? (
                <EmptyState emoji="⚠️" title="Ошибка" description={normalizeError(error)} />
            ) : items.length === 0 ? (
                <EmptyState emoji="🍽️" title="Пока нет рецептов" />
            ) : (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {items.map((r) => <RecipeCard key={r.id} recipe={r} />)}
                </div>
            )}
        </div>
    );
}