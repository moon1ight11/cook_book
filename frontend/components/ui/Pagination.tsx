'use client';

import { Button } from './Button';

export function Pagination({
    page,
    totalPages,
    onChange,
}: {
    page: number;
    totalPages: number;
    onChange: (page: number) => void;
}) {
    if (totalPages <= 1) return null;

    return (
        <div className="flex items-center justify-center gap-4 py-6">
            <Button
                variant="secondary"
                size="sm"
                disabled={page <= 1}
                onClick={() => onChange(page - 1)}
            >
                ← Назад
            </Button>
            <span className="text-sm text-warm-brown/80">
                стр. {page} из {totalPages}
            </span>
            <Button
                variant="secondary"
                size="sm"
                disabled={page >= totalPages}
                onClick={() => onChange(page + 1)}
            >
                Вперёд →
            </Button>
        </div>
    );
}