import { ReactNode } from 'react';

export function Card({
    children,
    className = '',
}: {
    children: ReactNode;
    className?: string;
}) {
    return (
        <div
            className={`rounded-xl border border-warm-brown/10 bg-parchment p-5 shadow-sm ${className}`}
        >
            {children}
        </div>
    );
}