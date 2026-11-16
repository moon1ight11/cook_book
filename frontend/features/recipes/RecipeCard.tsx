import Link from 'next/link';
import type { Recipe } from '@/types/api';
import { Badge } from '@/components/ui/Badge';
import { DIFFICULTY_LABEL, formatMinutes, tagEmoji } from '@/lib/format';

const DIFFICULTY_TONE = {
    easy: 'sage',
    medium: 'terracotta',
    hard: 'brown',
} as const;

export function RecipeCard({ recipe }: { recipe: Recipe }) {
    const totalTime = (recipe.prep_time ?? 0) + (recipe.cook_time ?? 0);
    const tone =
        recipe.difficulty && recipe.difficulty in DIFFICULTY_TONE
            ? DIFFICULTY_TONE[recipe.difficulty as keyof typeof DIFFICULTY_TONE]
            : 'neutral';
    const tags = recipe.tags ?? [];

    return (
        <Link
            href={`/recipes/${recipe.id}`}
            className="group flex flex-col overflow-hidden rounded-xl border border-warm-brown/10 bg-parchment shadow-sm transition-all hover:-translate-y-0.5 hover:shadow-md"
        >
            <div className="flex h-32 items-center justify-center bg-gradient-to-br from-warm-brown/10 to-terracotta/10 text-5xl">
                {tagEmoji(tags)}
            </div>
            <div className="flex flex-1 flex-col gap-2 p-4">
                <h3 className="line-clamp-2 font-serif text-lg font-semibold text-ink group-hover:text-terracotta">
                    {recipe.title}
                </h3>
                {recipe.description && (
                    <p className="line-clamp-2 text-sm text-warm-brown">{recipe.description}</p>
                )}
                <div className="mt-auto flex flex-wrap items-center gap-2 pt-2 text-xs text-warm-brown">
                    {recipe.difficulty && (
                        <Badge tone={tone}>{DIFFICULTY_LABEL[recipe.difficulty]}</Badge>
                    )}
                    {totalTime > 0 && <span>⏱ {formatMinutes(totalTime)}</span>}
                    {recipe.servings > 0 && <span>🍽 {recipe.servings} порц.</span>}
                </div>
                {tags.length > 0 && (
                    <div className="flex flex-wrap gap-1">
                        {tags.slice(0, 3).map((t) => (
                            <span
                                key={t}
                                className="rounded-full bg-warm-brown/10 px-2 py-0.5 text-[11px] text-warm-brown"
                            >
                                #{t}
                            </span>
                        ))}
                    </div>
                )}
            </div>
        </Link>
    );
}