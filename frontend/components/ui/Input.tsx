'use client';

import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
    label?: string;
    error?: string;
    hint?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
    ({ label, error, hint, id, className = '', ...rest }, ref) => {
        const inputId = id ?? rest.name;
        return (
            <div className="flex flex-col gap-1">
                {label && (
                    <label htmlFor={inputId} className="text-sm font-medium text-warm-brown">
                        {label}
                    </label>
                )}
                <input
                    ref={ref}
                    id={inputId}
                    className={`rounded-lg border bg-white px-3 py-2 text-warm-brown outline-none transition-colors placeholder:text-warm-brown/40 focus:border-terracotta focus:ring-2 focus:ring-terracotta/30 ${error ? 'border-red-500' : 'border-warm-brown/20'
                        } ${className}`}
                    {...rest}
                />
                {hint && !error && <span className="text-xs text-warm-brown/60">{hint}</span>}
                {error && <span className="text-xs text-red-600">{error}</span>}
            </div>
        );
    }
);
Input.displayName = 'Input';