'use client';

import { useState } from 'react';
import Link from 'next/link';
import { AuthGuard } from '@/components/layout/AuthGuard';
import { useAuthStore } from '@/lib/store';
import { useMe, useDeleteMe } from '@/features/profile/hooks';
import { ProfileForm } from '@/features/profile/ProfileForm';
import { Card } from '@/components/ui/Card';
import { Button } from '@/components/ui/Button';
import { Badge } from '@/components/ui/Badge';
import { Modal } from '@/components/ui/Modal';
import { Spinner } from '@/components/ui/Spinner';
import { formatDate } from '@/lib/format';

function ProfileInner() {
    const storedUser = useAuthStore((s) => s.user);
    const { data: user, isLoading } = useMe();
    const del = useDeleteMe();
    const [confirmOpen, setConfirmOpen] = useState(false);

    if (isLoading || !user) return <Spinner />;

    return (
        <div className="mx-auto flex max-w-2xl flex-col gap-6">
            <h1 className="font-serif text-3xl text-warm-brown">Профиль</h1>

            <Card>
                <div className="flex flex-wrap items-center justify-between gap-3">
                    <div>
                        <p className="font-serif text-xl text-warm-brown">{user.name}</p>
                        <p className="text-sm text-warm-brown/60">{user.email}</p>
                    </div>
                    <div className="flex items-center gap-2">
                        <Badge tone={user.role === 'Admin' ? 'terracotta' : 'sage'}>
                            {user.role === 'Admin' ? 'Администратор' : 'Пользователь'}
                        </Badge>
                        <Link href={`/users/${user.id}`}>
                            <Button variant="secondary" size="sm">Мои рецепты</Button>
                        </Link>
                    </div>
                </div>
                <p className="mt-3 text-xs text-warm-brown/50">
                    На сайте с {formatDate(user.created_at)}
                </p>
            </Card>

            <ProfileForm user={user} />

            <Card className="border-red-500/20">
                <h2 className="font-serif text-lg text-red-700">Опасная зона</h2>
                <p className="mt-1 text-sm text-warm-brown/60">
                    Удаление аккаунта необратимо. Рецепты останутся, но потеряют автора.
                </p>
                <div className="mt-3">
                    <Button variant="danger" size="sm" onClick={() => setConfirmOpen(true)}>
                        Удалить аккаунт
                    </Button>
                </div>
            </Card>

            <Modal open={confirmOpen} onClose={() => setConfirmOpen(false)} title="Удалить аккаунт?">
                <p className="text-sm text-warm-brown/70">
                    Это действие нельзя отменить. Вы выйдете из системы, а аккаунт будет удалён.
                </p>
                <div className="mt-4 flex justify-end gap-2">
                    <Button variant="ghost" onClick={() => setConfirmOpen(false)}>Отмена</Button>
                    <Button variant="danger" loading={del.isPending} onClick={() => del.mutate()}>
                        Удалить
                    </Button>
                </div>
            </Modal>
        </div>
    );
}

export default function ProfilePage() {
    return (
        <AuthGuard>
            <ProfileInner />
        </AuthGuard>
    );
}