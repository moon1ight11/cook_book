import { Suspense } from 'react';
import { GuestGuard } from '@/components/layout/GuestGuard';
import { LoginForm } from '@/features/auth/LoginForm';
import { Card } from '@/components/ui/Card';

export default function LoginPage() {
    return (
        <GuestGuard>
            <div className="mx-auto max-w-md py-10">
                <div className="mb-6 text-center">
                    <h1 className="font-serif text-3xl text-warm-brown">Вход</h1>
                    <p className="mt-1 text-sm text-warm-brown/60">
                        Рады видеть вас снова на кухне
                    </p>
                </div>
                <Card>
                    <Suspense fallback={null}>
                        <LoginForm />
                    </Suspense>
                </Card>
            </div>
        </GuestGuard>
    );
}