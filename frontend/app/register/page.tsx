import { GuestGuard } from '@/components/layout/GuestGuard';
import { RegisterForm } from '@/features/auth/RegisterForm';
import { Card } from '@/components/ui/Card';

export default function RegisterPage() {
    return (
        <GuestGuard>
            <div className="mx-auto max-w-md py-10">
                <div className="mb-6 text-center">
                    <h1 className="font-serif text-3xl text-warm-brown">Регистрация</h1>
                    <p className="mt-1 text-sm text-warm-brown/60">
                        Создайте аккаунт и делитесь рецептами
                    </p>
                </div>
                <Card>
                    <RegisterForm />
                </Card>
            </div>
        </GuestGuard>
    );
}