'use client';

import { useToastStore } from '@/lib/toast';

const KIND_STYLES: Record<string, string> = {
    success: 'bg-sage text-white',
    error: 'bg-red-600 text-white',
    info: 'bg-warm-brown text-white',
};

export function Toaster() {
    const items = useToastStore((s) => s.items);
    const remove = useToastStore((s) => s.remove);

    return (
        <div className="pointer-events-none fixed inset-x-0 top-4 z-[100] flex flex-col items-center gap-2">
            {items.map((t) => (
                <div
                    key={t.id}
                    onClick={() => remove(t.id)}
                    className={`pointer-events-auto cursor-pointer rounded-lg px-4 py-2 text-sm shadow-lg ${KIND_STYLES[t.kind]}`}
                >
                    {t.message}
                </div>
            ))}
        </div>
    );
}