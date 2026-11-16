'use client';

import { ReactNode, useEffect } from 'react';

export function Modal({
    open,
    onClose,
    title,
    children,
}: {
    open: boolean;
    onClose: () => void;
    title?: string;
    children: ReactNode;
}) {
    useEffect(() => {
        if (!open) return;
        const onKey = (e: KeyboardEvent) => {
            if (e.key === 'Escape') onClose();
        };
        document.addEventListener('keydown', onKey);
        document.body.style.overflow = 'hidden';
        return () => {
            document.removeEventListener('keydown', onKey);
            document.body.style.overflow = '';
        };
    }, [open, onClose]);

    if (!open) return null;

    return (
        <div
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
            onClick={onClose}
        >
            <div
                className="w-full max-w-md rounded-xl bg-cream p-6 shadow-xl"
                onClick={(e) => e.stopPropagation()}
            >
                {title && (
                    <h3 className="mb-4 font-serif text-xl text-warm-brown">{title}</h3>
                )}
                {children}
            </div>
        </div>
    );
}