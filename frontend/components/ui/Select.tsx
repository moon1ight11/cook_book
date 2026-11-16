'use client';

import { SelectHTMLAttributes, forwardRef } from 'react';

interface Option {
    value: string;
    label: string;
}

interface SelectProps extends SelectHTMLAttributes<HTMLSelectElement> {
    label?: string;
    error?: string;
    options: Option[];
    placeholder?: string;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
    ({ label, error, options, placeholder, id, className = '', ...rest }, ref) => {
        const fieldId = id ?? rest.name;
        return (
            <div className="flex flex-col gap-1">
                {label && (
                    <label htmlFor={fieldId} className="text-sm font-medium text-warm-brown">
                        {label}
                    </label>
                )}
                <select
                    ref={ref}
                    id={fieldId}
                    className={`rounded-lg border bg-white px-3 py-2 text-warm-brown outline-none transition-colors focus:border-terracotta focus:ring-2 focus:ring-terracotta/30 ${error ? 'border-red-500' : 'border-warm-brown/20'
                        } ${className}`}
                    {...rest}
                >
                    {placeholder && <option value="">{placeholder}</option>}
                    {options.map((o) => (
                        <option key={o.value} value={o.value}>
                            {o.label}
                        </option>
                    ))}
                </select>
                {error && <span className="text-xs text-red-600">{error}</span>}
            </div>
        );
    }
);
Select.displayName = 'Select';