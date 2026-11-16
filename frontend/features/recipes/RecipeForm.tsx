'use client';

import { Controller, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useRouter } from 'next/navigation';
import { Input } from '@/components/ui/Input';
import { Textarea } from '@/components/ui/Textarea';
import { Select } from '@/components/ui/Select';
import { Button } from '@/components/ui/Button';
import { IngredientEditor } from './IngredientEditor';
import { useCreateRecipe, useUpdateRecipe } from './hooks';
import { normalizeTag } from '@/lib/format';
import type { DetailedRecipe, Ingredient } from '@/types/api';

const ingredientSchema = z.object({
    name: z.string().min(1, 'Название'),
    quantity: z.string(),
    unit: z.string(),
});

const numberFromInput = (fallback: number) =>
    z.preprocess(
        (v) => (v === '' || v === null || v === undefined ? fallback : Number(v)),
        z.number().min(fallback === 1 ? 1 : 0, fallback === 1 ? 'Минимум 1' : 'Не может быть отрицательным')
    );

const schema = z.object({
    title: z.string().min(3, 'Минимум 3 символа').max(100, 'Максимум 100'),
    description: z.string().max(500, 'Максимум 500').optional(),
    ingredients: z.array(ingredientSchema).min(1, 'Добавьте хотя бы один ингредиент'),
    steps: z.array(z.string().min(1, 'Шаг не может быть пустым')).min(1, 'Добавьте шаги'),
    prep_time: numberFromInput(0),
    cook_time: numberFromInput(0),
    servings: numberFromInput(1),
    difficulty: z.enum(['easy', 'medium', 'hard']).optional(),
    tagsInput: z.string().optional(),
});

type FormValues = z.input<typeof schema>;
type ParsedValues = z.output<typeof schema>;

export function RecipeForm({ initial }: { initial?: DetailedRecipe }) {
    const router = useRouter();
    const create = useCreateRecipe();
    const update = useUpdateRecipe(initial?.id ?? '');
    const isEdit = !!initial;

    const {
        register,
        handleSubmit,
        control,
        formState: { errors },
    } = useForm<FormValues, unknown, ParsedValues>({
        resolver: zodResolver(schema),
        defaultValues: initial
            ? {
                title: initial.title,
                description: initial.description,
                ingredients: initial.ingredients?.length
                    ? initial.ingredients
                    : [{ name: '', quantity: '', unit: '' }],
                steps: initial.steps?.length ? initial.steps : [''],
                prep_time: String(initial.prep_time ?? '') as unknown as number,
                cook_time: String(initial.cook_time ?? '') as unknown as number,
                servings: String(initial.servings ?? '') as unknown as number,
                difficulty: (initial.difficulty || undefined) as 'easy' | 'medium' | 'hard' | undefined,
                tagsInput: initial.tags.join(', '),
            }
            : {
                title: '',
                description: '',
                ingredients: [{ name: '', quantity: '', unit: '' }],
                steps: [''],
                prep_time: '' as unknown as number,
                cook_time: '' as unknown as number,
                servings: '' as unknown as number,
                difficulty: 'medium',
                tagsInput: '',
            },
    });

    const onSubmit = (values: ParsedValues) => {
        const tags = (values.tagsInput ?? '')
            .split(',')
            .map(normalizeTag)
            .filter(Boolean);

        const payload = {
            title: values.title,
            description: values.description,
            ingredients: values.ingredients,
            steps: values.steps,
            prep_time: values.prep_time,
            cook_time: values.cook_time,
            servings: values.servings,
            difficulty: values.difficulty,
            tags,
        };

        if (isEdit) {
            update.mutate(payload, {
                onSuccess: () => router.push(`/recipes/${initial!.id}`),
            });
        } else {
            create.mutate(payload, {
                onSuccess: (recipe) => router.push(`/recipes/${recipe.id}`),
            });
        }
    };

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-5">
            <Input label="Название" error={errors.title?.message} {...register('title')} />

            <Textarea
                label="Описание"
                error={errors.description?.message}
                {...register('description')}
            />

            <div>
                <label className="mb-1 block text-sm font-medium text-warm-brown">
                    Ингредиенты
                </label>
                <Controller
                    control={control}
                    name="ingredients"
                    render={({ field }) => (
                        <IngredientEditor
                            value={field.value as Ingredient[]}
                            onChange={field.onChange}
                        />
                    )}
                />
                {errors.ingredients && (
                    <p className="mt-1 text-xs text-red-600">{errors.ingredients.message}</p>
                )}
            </div>

            <div>
                <label className="mb-1 block text-sm font-medium text-warm-brown">Шаги</label>
                <Controller
                    control={control}
                    name="steps"
                    render={({ field }) => (
                        <div className="flex flex-col gap-2">
                            {(field.value as string[]).map((step, idx) => (
                                <div key={idx} className="flex items-start gap-2">
                                    <span className="mt-3 w-6 shrink-0 text-sm text-warm-brown">{idx + 1}.</span>
                                    <div className="min-w-0 flex-1">
                                        <Textarea
                                            value={step}
                                            onChange={(e) => {
                                                const next = [...(field.value as string[])];
                                                next[idx] = e.target.value;
                                                field.onChange(next);
                                            }}
                                            className="min-h-[60px] w-full"
                                        />
                                    </div>
                                    <Button
                                        type="button"
                                        variant="ghost"
                                        size="sm"
                                        className="mt-1 shrink-0"
                                        onClick={() =>
                                            field.onChange((field.value as string[]).filter((_, i) => i !== idx))
                                        }
                                    >
                                        ✕
                                    </Button>
                                </div>
                            ))}
                            <Button
                                type="button"
                                variant="secondary"
                                size="sm"
                                onClick={() => field.onChange([...(field.value as string[]), ''])}
                            >
                                + Шаг
                            </Button>
                        </div>
                    )}
                />
                {errors.steps && (
                    <p className="mt-1 text-xs text-red-600">{errors.steps.message as string}</p>
                )}
            </div>

            <div className="grid gap-3 sm:grid-cols-3">
                <Input
                    label="Подготовка, мин"
                    type="number"
                    min={0}
                    placeholder="0"
                    error={errors.prep_time?.message}
                    {...register('prep_time')}
                />
                <Input
                    label="Готовка, мин"
                    type="number"
                    min={0}
                    placeholder="0"
                    error={errors.cook_time?.message}
                    {...register('cook_time')}
                />
                <Input
                    label="Порций"
                    type="number"
                    min={1}
                    placeholder="1"
                    error={errors.servings?.message}
                    {...register('servings')}
                />
            </div>

            <Select
                label="Сложность"
                error={errors.difficulty?.message}
                options={[
                    { value: 'easy', label: 'Легко' },
                    { value: 'medium', label: 'Средне' },
                    { value: 'hard', label: 'Сложно' },
                ]}
                {...register('difficulty')}
            />

            <Input
                label="Теги"
                placeholder="суп, мясо, быстро (через запятую)"
                error={errors.tagsInput?.message}
                {...register('tagsInput')}
            />

            <div className="flex gap-2">
                <Button type="submit" loading={create.isPending || update.isPending}>
                    {isEdit ? 'Сохранить' : 'Опубликовать'}
                </Button>
                <Button type="button" variant="ghost" onClick={() => router.back()}>
                    Отмена
                </Button>
            </div>
        </form>
    );
}