'use client';

import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import type { Ingredient } from '@/types/api';

export function IngredientEditor({
    value,
    onChange,
}: {
    value: Ingredient[];
    onChange: (v: Ingredient[]) => void;
}) {
    const update = (idx: number, patch: Partial<Ingredient>) => {
        const next = value.map((ing, i) => (i === idx ? { ...ing, ...patch } : ing));
        onChange(next);
    };

    const add = () => onChange([...value, { name: '', quantity: '', unit: '' }]);
    const remove = (idx: number) => onChange(value.filter((_, i) => i !== idx));

    return (
        <div className="flex flex-col gap-2">
            {value.map((ing, idx) => (
                <div key={idx} className="grid grid-cols-[1fr_100px_100px_auto] gap-2">
                    <Input
                        placeholder="Ингредиент"
                        value={ing.name}
                        onChange={(e) => update(idx, { name: e.target.value })}
                    />
                    <Input
                        placeholder="Кол-во"
                        value={ing.quantity}
                        onChange={(e) => update(idx, { quantity: e.target.value })}
                    />
                    <Input
                        placeholder="Ед."
                        value={ing.unit}
                        onChange={(e) => update(idx, { unit: e.target.value })}
                    />
                    <Button type="button" variant="ghost" size="sm" onClick={() => remove(idx)}>
                        ✕
                    </Button>
                </div>
            ))}
            <Button type="button" variant="secondary" size="sm" onClick={add}>
                + Ингредиент
            </Button>
        </div>
    );
}