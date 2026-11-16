'use client';

import { TextareaHTMLAttributes, forwardRef } from 'react';

interface TextareaProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
    label?: string;
    error?: string;
}

export const Textarea = forwardRef<HTMLTextAreaElement, TextareaProps>(
    ({ label, error, id, className = '', ...rest }, ref) => {
        const fieldId = id ?? rest.name;
        return (
            <div className="flex flex-col gap-1">
                {label && (
                    <label htmlFor={fieldId} className="text-sm font-medium text-warm-brown">
                        {label}
                    </label>
                )}
                <textarea
                    ref={ref}
                    id={fieldId}
                    className={`min-h-[100px] rounded-lg border bg-white px-3 py-2 text-warm-brown outline-none transition-colors placeholder:text-warm-brown/40 focus:border-terracotta focus:ring-2 focus:ring-terracotta/30 ${error ? 'border-red-500' : 'border-warm-brown/20'
                        } ${className}`}
                    {...rest}
                />
                {error && <span className="text-xs text-red-600">{error}</span>}
            </div>
        );
    }
);
Textarea.displayName = 'Textarea';