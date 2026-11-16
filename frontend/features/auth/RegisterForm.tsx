'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import Link from 'next/link';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { useSignUp } from './hooks';

const schema = z
    .object({
        name: z.string().min(1, 'Введите имя').max(100),
        email: z.string().email('Некорректный email'),
        password: z.string().min(6, 'Минимум 6 символов'),
        confirm: z.string(),
    })
    .refine((v) => v.password === v.confirm, {
        path: ['confirm'],
        message: 'Пароли не совпадают',
    });

type FormValues = z.infer<typeof schema>;

export function RegisterForm() {
    const signUp = useSignUp();
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormValues>({ resolver: zodResolver(schema) });

    const onSubmit = ({ name, email, password }: FormValues) =>
        signUp.mutate({ name, email, password });

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
            <Input label="Имя" error={errors.name?.message} {...register('name')} />
            <Input
                label="Email"
                type="email"
                autoComplete="email"
                error={errors.email?.message}
                {...register('email')}
            />
            <Input
                label="Пароль"
                type="password"
                autoComplete="new-password"
                error={errors.password?.message}
                {...register('password')}
            />
            <Input
                label="Повторите пароль"
                type="password"
                autoComplete="new-password"
                error={errors.confirm?.message}
                {...register('confirm')}
            />
            <Button type="submit" loading={signUp.isPending} className="mt-2">
                Зарегистрироваться
            </Button>
            <p className="text-center text-sm text-warm-brown/70">
                Уже есть аккаунт?{' '}
                <Link href="/login" className="font-medium text-terracotta hover:underline">
                    Войти
                </Link>
            </p>
        </form>
    );
}