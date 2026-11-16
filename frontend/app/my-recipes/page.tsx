'use client';

import Link from 'next/link';
import { AuthGuard } from '@/components/layout/AuthGuard';
import { useRecipesByUser } from '@/features/recipes/hooks';
import { RecipeCard } from '@/features/recipes/RecipeCard';
import { useAuthStore } from '@/lib/store';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import { normalizeError } from '@/lib/errors';

function MyRecipesInner() {
    const user = useAuthStore((s) => s.user);
    const { data, isLoading, isError, error } = useRecipesByUser(user?.id);

    // бэк отдаёт без сортировки — сортируем по created_at DESC на клиенте
    const sorted = data
        ? [...data].sort((a, b) => +new Date(b.created_at) - +new Date(a.created_at))
        : [];

    return (
        <div className="flex flex-col gap-6">
            <div className="flex items-center justify-between">
                <h1 className="font-serif text-3xl text-warm-brown">Мои рецепты</h1>
                <Link href="/recipes/new"><Button>+ Новый рецепт</Button></Link>
            </div>
            {isLoading ? (
                <Spinner />
            ) : isError ? (
                <EmptyState emoji="⚠️" title="Ошибка" description={normalizeError(error)} />
            ) : sorted.length === 0 ? (
                <EmptyState
                    emoji="📖"
                    title="У вас пока нет рецептов"
                    action={<Link href="/recipes/new"><Button>Добавить первый рецепт</Button></Link>}
                />
            ) : (
                <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {sorted.map((r) => <RecipeCard key={r.id} recipe={r} />)}
                </div>
            )}
        </div>
    );
}

export default function MyRecipesPage() {
    return (
        <AuthGuard>
            <MyRecipesInner />
        </AuthGuard>
    );
}