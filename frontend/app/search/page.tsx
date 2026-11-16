'use client';

import { Suspense, useState } from 'react';
import { useSearchParams, useRouter, usePathname } from 'next/navigation';
import { useRecipesByTags } from '@/features/recipes/hooks';
import { RecipeCard } from '@/features/recipes/RecipeCard';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { Spinner } from '@/components/ui/Spinner';
import { EmptyState } from '@/components/ui/EmptyState';
import { normalizeError } from '@/lib/errors';
import { normalizeTag } from '@/lib/format';

function SearchInner() {
    const searchParams = useSearchParams();
    const router = useRouter();
    const pathname = usePathname();

    const initial = (searchParams.getAll('tags') ?? []).map(normalizeTag).filter(Boolean);
    const [tags, setTags] = useState<string[]>(initial);
    const [input, setInput] = useState('');

    const { data, isLoading, isError, error, isFetching } = useRecipesByTags(tags, tags.length > 0);

    const syncUrl = (next: string[]) => {
        const qs = next.map((t) => `tags=${encodeURIComponent(t)}`).join('&');
        router.replace(qs ? `${pathname}?${qs}` : pathname);
    };

    const addTag = () => {
        const t = normalizeTag(input);
        if (!t || tags.includes(t)) return;
        const next = [...tags, t];
        setTags(next);
        setInput('');
        syncUrl(next);
    };

    const removeTag = (t: string) => {
        const next = tags.filter((x) => x !== t);
        setTags(next);
        syncUrl(next);
    };

    const results = data ?? [];

    return (
        <div className="flex flex-col gap-6">
            <div>
                <h1 className="font-serif text-3xl text-warm-brown">Поиск по тегам</h1>
                <p className="mt-1 text-sm text-warm-brown/60">
                    Найдём рецепты, у которых есть <b>хотя бы один</b> из выбранных тегов
                </p>
            </div>

            <div className="flex flex-col gap-3">
                <div className="flex gap-2">
                    <Input
                        placeholder="Например: суп"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter') { e.preventDefault(); addTag(); }
                        }}
                    />
                    <Button onClick={addTag}>Добавить</Button>
                </div>
                {tags.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                        {tags.map((t) => (
                            <button key={t} onClick={() => removeTag(t)} title="Удалить">
                                <Badge tone="terracotta">#{t} ✕</Badge>
                            </button>
                        ))}
                    </div>
                )}
            </div>

            {tags.length === 0 ? (
                <EmptyState emoji="🏷️" title="Добавьте хотя бы один тег" />
            ) : isLoading ? (
                <Spinner />
            ) : isError ? (
                <EmptyState emoji="⚠️" title="Ошибка поиска" description={normalizeError(error)} />
            ) : results.length === 0 ? (
                <EmptyState emoji="🔍" title="Ничего не найдено" description="Попробуйте другие теги" />
            ) : (
                <div className={`grid gap-4 sm:grid-cols-2 lg:grid-cols-3 ${isFetching ? 'opacity-60' : ''}`}>
                    {results.map((r) => <RecipeCard key={r.id} recipe={r} />)}
                </div>
            )}
        </div>
    );
}

export default function SearchPage() {
    return (
        <Suspense fallback={<Spinner />}>
            <SearchInner />
        </Suspense>
    );
}