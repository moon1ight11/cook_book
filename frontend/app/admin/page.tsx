import { AuthGuard } from '@/components/layout/AuthGuard';
import { AdminPanels } from '@/features/admin/AdminPanels';

export default function AdminPage() {
    return (
        <AuthGuard requireAdmin>
            <div className="flex flex-col gap-6">
                <div>
                    <h1 className="font-serif text-3xl text-warm-brown">Админка</h1>
                    <p className="mt-1 text-sm text-warm-brown/60">
                        Точечные операции по id. Списка пользователей в API нет, поэтому
                        id нужно вводить вручную. Свой id можно посмотреть в{' '}
                        <a href="/profile" className="text-terracotta hover:underline">профиле</a>.
                    </p>
                </div>
                <AdminPanels />
            </div>
        </AuthGuard>
    );
}