'use client';

import { useEffect, useState } from 'react';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import type { Difficulty } from '@/types/api';

export interface FiltersValue {
    title: string;
    difficulty: '' | Difficulty;
}

export function RecipeFilters({
    value,
    onChange,
}: {
    value: FiltersValue;
    onChange: (v: FiltersValue) => void;
}) {
    const [title, setTitle] = useState(value.title);

    // debounce поиска по title — бэк ищет regex-подстрокой
    useEffect(() => {
        const t = setTimeout(() => {
            if (title !== value.title) onChange({ ...value, title });
        }, 400);
        return () => clearTimeout(t);
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [title]);

    return (
        <div className="grid gap-3 sm:grid-cols-[1fr_200px]">
            <Input
                placeholder="Поиск по названию…"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
            />
            <Select
                value={value.difficulty}
                onChange={(e) =>
                    onChange({ ...value, difficulty: e.target.value as FiltersValue['difficulty'] })
                }
                placeholder="Любая сложность"
                options={[
                    { value: 'easy', label: 'Легко' },
                    { value: 'medium', label: 'Средне' },
                    { value: 'hard', label: 'Сложно' },
                ]}
            />
        </div>
    );
}