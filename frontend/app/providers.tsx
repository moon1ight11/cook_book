'use client';

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useEffect, useState } from 'react';
import { api } from '@/lib/api';
import { useAuthStore } from '@/lib/store';
import { Toaster } from '@/components/ui/Toaster';

function AuthBootstrap({ children }: { children: React.ReactNode }) {
    const setUser = useAuthStore((s) => s.setUser);
    const setStatus = useAuthStore((s) => s.setStatus);

    useEffect(() => {
        const refresh = localStorage.getItem('refresh_token');
        if (!refresh) {
            setStatus('guest');
            return;
        }
        setStatus('loading');
        api
            .get('/private/users')
            .then(({ data }) => {
                setUser({ id: data.id, name: data.name, email: data.email, role: data.role });
                setStatus('authenticated');
            })
            .catch(() => {
                setUser(null);
                setStatus('guest');
            });
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    return <>{children}</>;
}

export function Providers({ children }: { children: React.ReactNode }) {
    const [client] = useState(
        () =>
            new QueryClient({
                defaultOptions: {
                    queries: { staleTime: 30_000, retry: 1, refetchOnWindowFocus: false },
                },
            })
    );

    return (
        <QueryClientProvider client={client}>
            <AuthBootstrap>{children}</AuthBootstrap>
            <Toaster />
        </QueryClientProvider>
    );
}