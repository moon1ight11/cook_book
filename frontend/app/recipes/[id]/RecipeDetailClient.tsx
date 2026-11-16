'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useRecipe, useDeleteRecipe } from '@/features/recipes/hooks';
import { useAuthStore } from '@/lib/store';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { Modal } from '@/components/ui/Modal';
import { DIFFICULTY_LABEL, formatDate, formatMinutes, tagEmoji } from '@/lib/format';

export default function RecipeDetailClient({ id }: { id: string }) {
    const router = useRouter();
    const { user } = useAuthStore();
    const { data: recipe, isLoading, isError } = useRecipe(id);
    const del = useDeleteRecipe();
    const [confirmOpen, setConfirmOpen] = useState(false);

    if (isLoading) return <Spinner />;

    if (isError || !recipe) {
        return (
            <EmptyState
                emoji="🍳"
                title="Рецепт не найден"
                description="Возможно, он был удалён или ссылка неверна"
                action={<Link href="/"><Button>К списку рецептов</Button></Link>}
            />
        );
    }

    const isOwner = user?.id === recipe.owner_id;
    const totalTime = (recipe.prep_time ?? 0) + (recipe.cook_time ?? 0);
    const ingredients = recipe.ingredients ?? [];
    const steps = recipe.steps ?? [];
    const tags = recipe.tags ?? [];

    const handleDelete = () => {
        del.mutate(recipe.id, { onSuccess: () => router.push('/') });
    };

    return (
        <article className="mx-auto max-w-3xl">
            <div className="mb-6 flex h-48 items-center justify-center rounded-2xl bg-gradient-to-br from-warm-brown/10 to-terracotta/15 text-7xl">
                {tagEmoji(tags)}
            </div>

            <div className="flex flex-wrap items-start justify-between gap-4">
                <div>
                    <h1 className="font-serif text-3xl text-ink">{recipe.title}</h1>
                    <p className="mt-1 text-sm text-warm-brown">
                        от{' '}
                        <Link href={`/users/${recipe.owner_id}`} className="hover:text-terracotta hover:underline">
                            автора
                        </Link>{' '}
                        · {formatDate(recipe.created_at)}
                    </p>
                </div>
                {isOwner && (
                    <div className="flex gap-2">
                        <Link href={`/recipes/${recipe.id}/edit`}>
                            <Button variant="secondary" size="sm">Редактировать</Button>
                        </Link>
                        <Button variant="danger" size="sm" onClick={() => setConfirmOpen(true)}>
                            Удалить
                        </Button>
                    </div>
                )}
            </div>

            {recipe.description && (
                <p className="mt-4 text-warm-brown">{recipe.description}</p>
            )}

            <div className="mt-4 flex flex-wrap items-center gap-3 text-sm text-warm-brown">
                {recipe.difficulty && <Badge tone="terracotta">{DIFFICULTY_LABEL[recipe.difficulty]}</Badge>}
                {totalTime > 0 && <span>⏱ {formatMinutes(totalTime)}</span>}
                {recipe.prep_time > 0 && <span>Подготовка: {formatMinutes(recipe.prep_time)}</span>}
                {recipe.cook_time > 0 && <span>Готовка: {formatMinutes(recipe.cook_time)}</span>}
                {recipe.servings > 0 && <span>🍽 {recipe.servings} порц.</span>}
            </div>

            {tags.length > 0 && (
                <div className="mt-3 flex flex-wrap gap-1">
                    {tags.map((t) => (
                        <Link key={t} href={`/search?tags=${encodeURIComponent(t)}`}>
                            <Badge tone="neutral">#{t}</Badge>
                        </Link>
                    ))}
                </div>
            )}

            <section className="mt-8">
                <h2 className="font-serif text-xl text-ink">Ингредиенты</h2>
                {ingredients.length === 0 ? (
                    <p className="mt-3 rounded-xl border border-dashed border-warm-brown/20 bg-parchment/50 px-4 py-6 text-sm text-warm-brown">
                        Ингредиенты не указаны
                    </p>
                ) : (
                    <ul className="mt-3 divide-y divide-warm-brown/10 rounded-xl border border-warm-brown/10 bg-parchment">
                        {ingredients.map((ing, i) => (
                            <li key={i} className="flex items-center justify-between px-4 py-2 text-sm">
                                <span className="text-ink">{ing.name}</span>
                                <span className="text-warm-brown">
                                    {ing.quantity} {ing.unit}
                                </span>
                            </li>
                        ))}
                    </ul>
                )}
            </section>

            <section className="mt-8">
                <h2 className="font-serif text-xl text-ink">Приготовление</h2>
                {steps.length === 0 ? (
                    <p className="mt-3 rounded-xl border border-dashed border-warm-brown/20 bg-parchment/50 px-4 py-6 text-sm text-warm-brown">
                        Шаги не указаны
                    </p>
                ) : (
                    <ol className="mt-3 flex flex-col gap-3">
                        {steps.map((step, i) => (
                            <li key={i} className="flex gap-3 rounded-xl border border-warm-brown/10 bg-parchment p-4">
                                <span className="flex h-7 w-7 flex-shrink-0 items-center justify-center rounded-full bg-terracotta text-sm font-semibold text-white">
                                    {i + 1}
                                </span>
                                <p className="text-sm text-warm-brown">{step}</p>
                            </li>
                        ))}
                    </ol>
                )}
            </section>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)} title="Удалить рецепт?">
                <p className="text-sm text-warm-brown">
                    Действие необратимо. Рецепт «{recipe.title}» будет удалён.
                </p>
                <div className="mt-4 flex justify-end gap-2">
                    <Button variant="ghost" onClick={() => setConfirmOpen(false)}>Отмена</Button>
                    <Button variant="danger" loading={del.isPending} onClick={handleDelete}>
                        Удалить
                    </Button>
                </div>
            </Modal>
        </article>
    );
}