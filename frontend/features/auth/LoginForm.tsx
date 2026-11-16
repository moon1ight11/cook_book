'use client';

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import Link from 'next/link';
import { Input } from '@/components/ui/Input';
import { Button } from '@/components/ui/Button';
import { useSignIn } from './hooks';

const schema = z.object({
    email: z.string().email('Некорректный email'),
    password: z.string().min(1, 'Введите пароль'),
});

type FormValues = z.infer<typeof schema>;

export function LoginForm() {
    const signIn = useSignIn();
    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<FormValues>({ resolver: zodResolver(schema) });

    const onSubmit = (values: FormValues) => signIn.mutate(values);

    return (
        <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4">
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
                autoComplete="current-password"
                error={errors.password?.message}
                {...register('password')}
            />
            <Button type="submit" loading={signIn.isPending} className="mt-2">
                Войти
            </Button>
            <p className="text-center text-sm text-warm-brown/70">
                Нет аккаунта?{' '}
                <Link href="/register" className="font-medium text-terracotta hover:underline">
                    Зарегистрироваться
                </Link>
            </p>
        </form>
    );
}