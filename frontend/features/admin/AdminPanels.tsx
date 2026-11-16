'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { Card } from '@/components/ui/Card';
import { Input } from '@/components/ui/Input';
import { Select } from '@/components/ui/Select';
import { Button } from '@/components/ui/Button';
import { Modal } from '@/components/ui/Modal';
import {
    useCreateAdmin,
    useUpdateUserRole,
    useDeleteAnyUser,
    useDeleteAnyRecipe,
} from './hooks';
import { useAuthStore } from '@/lib/store';
import type { Role } from '@/types/api';

// ---------- создать админа ----------

const createAdminSchema = z.object({
    name: z.string().min(1, 'Введите имя'),
    email: z.string().email('Некорректный email'),
    password: z.string().min(6, 'Минимум 6 символов'),
});
type CreateAdminValues = z.infer<typeof createAdminSchema>;

function CreateAdminPanel() {
    const create = useCreateAdmin();
    const {
        register,
        handleSubmit,
        reset,
        formState: { errors },
    } = useForm<CreateAdminValues>({ resolver: zodResolver(createAdminSchema) });

    const onSubmit = (v: CreateAdminValues) =>
        create.mutate(v, { onSuccess: () => reset() });

    return (
        <Card>
            <h2 className="font-serif text-lg text-warm-brown">Создать администратора</h2>
            <p className="mt-1 text-sm text-warm-brown/60">
                Новый пользователь сразу получит роль Admin.
            </p>
            <form onSubmit={handleSubmit(onSubmit)} className="mt-4 flex flex-col gap-3">
                <Input label="Имя" error={errors.name?.message} {...register('name')} />
                <Input label="Email" type="email" error={errors.email?.message} {...register('email')} />
                <Input
                    label="Пароль"
                    type="password"
                    autoComplete="new-password"
                    error={errors.password?.message}
                    {...register('password')}
                />
                <Button type="submit" loading={create.isPending}>Создать</Button>
            </form>
        </Card>
    );
}

// ---------- сменить роль ----------

const roleSchema = z.object({
    userId: z.string().min(1, 'Введите user_id'),
    role: z.enum(['Basic', 'Admin']),
});
type RoleValues = z.infer<typeof roleSchema>;

function UpdateRolePanel() {
    const me = useAuthStore((s) => s.user);
    const update = useUpdateUserRole();
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<RoleValues>({
        resolver: zodResolver(roleSchema),
        defaultValues: { userId: me?.id ?? '', role: 'Basic' },
    });

    const onSubmit = (v: RoleValues) => update.mutate({ userId: v.userId, role: v.role as Role });

    return (
        <Card>
            <h2 className="font-serif text-lg text-warm-brown">Сменить роль пользователя</h2>
            <p className="mt-1 text-sm text-warm-brown/60">
                Введите Mongo ObjectId пользователя (24 hex-символа).
            </p>
            <form onSubmit={handleSubmit(onSubmit)} className="mt-4 flex flex-col gap-3">
                <Input
                    label="user_id"
                    placeholder="507f1f77bcf86cd799439011"
                    error={errors.userId?.message}
                    {...register('userId')}
                />
                <Select
                    label="Новая роль"
                    options={[
                        { value: 'Basic', label: 'Basic' },
                        { value: 'Admin', label: 'Admin' },
                    ]}
                    {...register('role')}
                />
                <Button type="submit" loading={update.isPending}>Обновить роль</Button>
            </form>
        </Card>
    );
}

// ---------- удалить пользователя ----------

function DeleteUserPanel() {
    const [open, setOpen] = useState(false);
    const del = useDeleteAnyUser();
    const {
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm<{ userId: string }>({
        resolver: zodResolver(z.object({ userId: z.string().min(1, 'Введите user_id') })),
    });
    const userId = watch('userId');

    const submit = handleSubmit(() => setOpen(true));

    const confirm = () => {
        del.mutate(userId, { onSuccess: () => setOpen(false) });
    };

    return (
        <Card className="border-red-500/20">
            <h2 className="font-serif text-lg text-red-700">Удалить пользователя</h2>
            <p className="mt-1 text-sm text-warm-brown/60">
                Удалит любого пользователя по id, включая администраторов.
            </p>
            <form onSubmit={submit} className="mt-4 flex flex-col gap-3">
                <Input label="user_id" error={errors.userId?.message} {...register('userId')} />
                <Button type="submit" variant="danger">Удалить</Button>
            </form>

            <Modal open={open} onClose={() => setOpen(false)} title="Удалить пользователя?">
                <p className="text-sm text-warm-brown/70">
                    Пользователь <code className="rounded bg-warm-brown/10 px-1">{userId}</code> будет удалён безвозвратно.
                </p>
                <div className="mt-4 flex justify-end gap-2">
                    <Button variant="ghost" onClick={() => setOpen(false)}>Отмена</Button>
                    <Button variant="danger" loading={del.isPending} onClick={confirm}>Удалить</Button>
                </div>
            </Modal>
        </Card>
    );
}

// ---------- удалить рецепт ----------

function DeleteRecipePanel() {
    const [open, setOpen] = useState(false);
    const del = useDeleteAnyRecipe();
    const {
        register,
        handleSubmit,
        watch,
        formState: { errors },
    } = useForm<{ recipeId: string }>({
        resolver: zodResolver(z.object({ recipeId: z.string().min(1, 'Введите recipe_id') })),
    });
    const recipeId = watch('recipeId');

    const submit = handleSubmit(() => setOpen(true));
    const confirm = () => del.mutate(recipeId, { onSuccess: () => setOpen(false) });

    return (
        <Card className="border-red-500/20">
            <h2 className="font-serif text-lg text-red-700">Удалить рецепт</h2>
            <p className="mt-1 text-sm text-warm-brown/60">
                Удалит любой рецепт по id, независимо от автора.
            </p>
            <form onSubmit={submit} className="mt-4 flex flex-col gap-3">
                <Input label="recipe_id" error={errors.recipeId?.message} {...register('recipeId')} />
                <Button type="submit" variant="danger">Удалить</Button>
            </form>

            <Modal open={open} onClose={() => setOpen(false)} title="Удалить рецепт?">
                <p className="text-sm text-warm-brown/70">
                    Рецепт <code className="rounded bg-warm-brown/10 px-1">{recipeId}</code> будет удалён безвозвратно.
                </p>
                <div className="mt-4 flex justify-end gap-2">
                    <Button variant="ghost" onClick={() => setOpen(false)}>Отмена</Button>
                    <Button variant="danger" loading={del.isPending} onClick={confirm}>Удалить</Button>
                </div>
            </Modal>
        </Card>
    );
}

export function AdminPanels() {
    return (
        <div className="grid gap-6 md:grid-cols-2">
            <CreateAdminPanel />
            <UpdateRolePanel />
            <DeleteUserPanel />
            <DeleteRecipePanel />
        </div>
    );
}