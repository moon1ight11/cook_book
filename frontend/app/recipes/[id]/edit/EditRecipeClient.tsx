'use client';

import Link from 'next/link';
import { AuthGuard } from '@/components/layout/AuthGuard';
import { RecipeForm } from '@/features/recipes/RecipeForm';
import { useRecipe } from '@/features/recipes/hooks';
import { useAuthStore } from '@/lib/store';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';

function EditInner({ id }: { id: string }) {
    const { user } = useAuthStore();
    const { data: recipe, isLoading, isError } = useRecipe(id);

    if (isLoading) return <Spinner />;
    if (isError || !recipe) {
        return (
            <EmptyState
                emoji="🍳"
                title="Рецепт не найден"
                action={<Link href="/"><Button>К списку</Button></Link>}
            />
        );
    }
    if (user?.id !== recipe.owner_id) {
        return (
            <EmptyState
                emoji="🔒"
                title="Нет доступа"
                description="Редактировать рецепт может только его автор"
                action={<Link href={`/recipes/${recipe.id}`}><Button>К рецепту</Button></Link>}
            />
        );
    }
    return (
        <div className="mx-auto max-w-2xl">
            <h1 className="mb-6 font-serif text-3xl text-ink">Редактирование</h1>
            <RecipeForm initial={recipe} />
        </div>
    );
}

export default function EditRecipeClient({ id }: { id: string }) {
    return (
        <AuthGuard>
            <EditInner id={id} />
        </AuthGuard>
    );
}