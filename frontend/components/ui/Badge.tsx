import { ReactNode } from 'react';

type Tone = 'neutral' | 'terracotta' | 'sage' | 'brown';

const TONES: Record<Tone, string> = {
    neutral: 'bg-warm-brown/10 text-warm-brown',
    terracotta: 'bg-terracotta/15 text-terracotta',
    sage: 'bg-sage/20 text-sage',
    brown: 'bg-warm-brown/20 text-warm-brown',
};

export function Badge({
    children,
    tone = 'neutral',
    className = '',
}: {
    children: ReactNode;
    tone?: Tone;
    className?: string;
}) {
    return (
        <span
            className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${TONES[tone]} ${className}`}
        >
            {children}
        </span>
    );
}