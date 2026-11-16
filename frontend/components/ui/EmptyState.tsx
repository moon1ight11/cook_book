import { ReactNode } from 'react';

export function EmptyState({
    emoji = '🍽️',
    title,
    description,
    action,
}: {
    emoji?: string;
    title: string;
    description?: string;
    action?: ReactNode;
}) {
    return (
        <div className="flex flex-col items-center justify-center gap-3 rounded-xl border border-dashed border-warm-brown/20 bg-parchment/50 px-6 py-14 text-center">
            <span className="text-5xl">{emoji}</span>
            <h3 className="font-serif text-xl text-warm-brown">{title}</h3>
            {description && <p className="max-w-md text-sm text-warm-brown/70">{description}</p>}
            {action && <div className="mt-2">{action}</div>}
        </div>
    );
}