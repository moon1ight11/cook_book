'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRecipes } from '@/features/recipes/hooks';
import { RecipeCard } from '@/features/recipes/RecipeCard';
import { RecipeFilters, FiltersValue } from '@/features/recipes/RecipeFilters';
import { Pagination } from '@/components/ui/Pagination';
import { EmptyState } from '@/components/ui/EmptyState';
import { Button } from '@/components/ui/Button';
import { normalizeError } from '@/lib/errors';

const PAGE_SIZE = 12;

function RecipeCardSkeleton() {
  return (
    <div className="overflow-hidden rounded-xl border border-warm-brown/10 bg-parchment">
      <div className="h-32 animate-pulse bg-warm-brown/10" />
      <div className="space-y-2 p-4">
        <div className="h-4 w-3/4 animate-pulse rounded bg-warm-brown/10" />
        <div className="h-3 w-full animate-pulse rounded bg-warm-brown/10" />
        <div className="h-3 w-2/3 animate-pulse rounded bg-warm-brown/10" />
      </div>
    </div>
  );
}

export default function HomePage() {
  const [page, setPage] = useState(1);
  const [filters, setFilters] = useState<FiltersValue>({ title: '', difficulty: '' });

  const params = {
    page,
    page_size: PAGE_SIZE,
    ...(filters.title ? { title: filters.title } : {}),
    ...(filters.difficulty ? { difficulty: filters.difficulty } : {}),
  };

  const { data, isLoading, isError, error, isFetching } = useRecipes(params);
  const recipes = data?.data ?? [];
  const totalPages = data?.total_pages ?? 1;

  const onFiltersChange = (v: FiltersValue) => {
    setFilters(v);
    setPage(1);
  };

  const dim = isFetching && !isLoading;

  return (
    <div className="flex flex-col gap-6">
      <section className="rounded-2xl bg-gradient-to-br from-parchment to-cream p-6 sm:p-10">
        <h1 className="font-serif text-3xl text-ink sm:text-4xl">CookBook</h1>
        <p className="mt-2 max-w-xl text-warm-brown">
          Рецепты, которые хочется готовить. Ищите по названию, сложности или тегам.
        </p>
        <div className="mt-4 flex flex-wrap gap-2">
          <Link href="/search"><Button variant="secondary">Поиск по тегам</Button></Link>
          <Link href="/recipes/new"><Button>+ Добавить рецепт</Button></Link>
        </div>
      </section>

      <RecipeFilters value={filters} onChange={onFiltersChange} />

      {isLoading ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: PAGE_SIZE }).map((_, i) => (
            <RecipeCardSkeleton key={i} />
          ))}
        </div>
      ) : isError ? (
        <EmptyState
          emoji="⚠️"
          title="Не удалось загрузить рецепты"
          description={normalizeError(error)}
        />
      ) : recipes.length === 0 ? (
        <EmptyState
          title="Рецептов пока нет"
          description={
            filters.title || filters.difficulty
              ? 'Попробуйте изменить фильтры'
              : 'Станьте первым, кто поделится рецептом'
          }
          action={<Link href="/recipes/new"><Button>Добавить рецепт</Button></Link>}
        />
      ) : (
        <>
          <div
            className={`grid gap-4 sm:grid-cols-2 lg:grid-cols-3 transition-opacity ${
              dim ? 'opacity-60' : 'opacity-100'
            }`}
          >
            {recipes.map((r) => (
              <RecipeCard key={r.id} recipe={r} />
            ))}
          </div>
          <Pagination page={page} totalPages={totalPages} onChange={setPage} />
        </>
      )}
    </div>
  );
}