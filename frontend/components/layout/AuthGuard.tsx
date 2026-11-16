'use client';

import { useRouter } from 'next/navigation';
import { useEffect } from 'react';
import { useAuthStore } from '@/lib/store';
import { FullScreenSpinner } from '@/components/ui/Spinner';

export function AuthGuard({
    children,
    requireAdmin = false,
}: {
    children: React.ReactNode;
    requireAdmin?: boolean;
}) {
    const router = useRouter();
    const { status, user } = useAuthStore();

    useEffect(() => {
        if (status === 'guest') {
            const next = encodeURIComponent(window.location.pathname + window.location.search);
            router.replace(`/login?next=${next}`);
        } else if (status === 'authenticated' && requireAdmin && user?.role !== 'Admin') {
            router.replace('/');
        }
    }, [status, user, requireAdmin, router]);

    if (status === 'idle' || status === 'loading') {
        return <FullScreenSpinner />;
    }
    if (status === 'guest') {
        return <FullScreenSpinner />;
    }
    if (requireAdmin && user?.role !== 'Admin') {
        return <FullScreenSpinner />;
    }
    return <>{children}</>;
}