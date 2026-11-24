'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { useRouter } from 'next/navigation';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { Card } from '@/components/ui/Card';
import { useUpdateMe } from './hooks';
import { useAuthStore } from '@/lib/store';
import { toast } from '@/lib/toast';
import type { User } from '@/types/api';

const schema = z
    .object({
        name: z.string().optional(),
        email: z.string().optional(),
        password: z.string().optional(),
        confirm: z.string().optional(),
    })
    .refine((v) => !v.password || v.password.length >= 6, {
        path: ['password'],
        message: 'Минимум 6 символов',
    })
    .refine((v) => v.password === v.confirm, {
        path: ['confirm'],
        message: 'Пароли не совпадают',
    })
    .refine((v) => v.name !== '' || v.email !== '' || v.password !== '', {
        path: ['name'],
        message: 'Измените хотя бы одно поле',
    });

type FormValues = z.infer<typeof schema>;

export function ProfileForm({ user }: { user: User }) {
    const router = useRouter();
    const update = useUpdateMe();
    const logoutLocal = useAuthStore((s) => s.logoutLocal);

    const {
        register,
        handleSubmit,
        formState: { errors, isDirty },
    } = useForm<FormValues>({
        resolver: zodResolver(schema),
        defaultValues: { name: user.name, email: user.email, password: '', confirm: '' },
    });

    const onSubmit = (values: FormValues) => {
        const payload: { name?: string; email?: string; password?: string } = {};
        if (values.name && values.name !== user.name) payload.name = values.name;
        if (values.email && values.email !== user.email) payload.email = values.email;
        if (values.password) payload.password = values.password;

        if (Object.keys(payload).length === 0) {
            toast.info('Нет изменений');
            return;
        }

        const changesCredentials = !!(payload.email || payload.password);

        update.mutate(payload, {
            onSuccess: () => {
                if (changesCredentials) {
                    toast.info('Данные входа изменены. Войдите заново.');
                    logoutLocal();
                    router.replace('/login');
                } else {
                    toast.success('Профиль обновлён');
                }
            },
        });
    };

    return (
        <Card>
            <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
                <Input label="Имя" error={errors.name?.message} {...register('name')} />
                <Input label="Email" type="email" error={errors.email?.message} {...register('email')} />

                <hr className="my-2 border-warm-brown/10" />
                <p className="text-xs text-warm-brown/60">
                    Смена email или пароля потребует повторного входа
                </p>

                <Input
                    label="Новый пароль (необязательно)"
                    type="password"
                    autoComplete="new-password"
                    error={errors.password?.message}
                    {...register('password')}
                />
                <Input
                    label="Повторите новый пароль"
                    type="password"
                    autoComplete="new-password"
                    error={errors.confirm?.message}
                    {...register('confirm')}
                />

                <div className="flex justify-end">
                    <Button type="submit" loading={update.isPending} disabled={!isDirty}>
                        Сохранить
                    </Button>
                </div>
            </form>
        </Card>
    );
}