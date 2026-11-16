'use client';

import { ButtonHTMLAttributes, forwardRef } from 'react';

type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';
type Size = 'sm' | 'md' | 'lg';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: Variant;
    size?: Size;
    loading?: boolean;
}

const VARIANTS: Record<Variant, string> = {
    primary:
        'bg-terracotta text-white hover:bg-terracotta/90 disabled:bg-terracotta/50',
    secondary:
        'bg-parchment text-warm-brown border border-warm-brown/20 hover:bg-warm-brown/10',
    ghost: 'bg-transparent text-warm-brown hover:bg-warm-brown/10',
    danger: 'bg-red-600 text-white hover:bg-red-700 disabled:bg-red-600/50',
};

const SIZES: Record<Size, string> = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
    ({ variant = 'primary', size = 'md', loading, className = '', children, disabled, ...rest }, ref) => {
        return (
            <button
                ref={ref}
                disabled={disabled || loading}
                className={`inline-flex items-center justify-center gap-2 rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-terracotta/40 disabled:cursor-not-allowed ${VARIANTS[variant]} ${SIZES[size]} ${className}`}
                {...rest}
            >
                {loading && (
                    <span className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                )}
                {children}
            </button>
        );
    }
);
Button.displayName = 'Button';