'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useAuthStore } from '@/lib/store';

export function GuestGuard({ children }: { children: React.ReactNode }) {
    const router = useRouter();
    const status = useAuthStore((s) => s.status);

    useEffect(() => {
        if (status === 'authenticated') {
            router.replace('/');
        }
    }, [status, router]);

    return <>{children}</>;
}