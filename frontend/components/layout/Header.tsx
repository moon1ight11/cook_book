'use client';

import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { useState } from 'react';
import { useAuthStore } from '@/lib/store';
import { Button } from '@/components/ui/Button';
import { useSignOut } from '@/features/auth/hooks';
import Image from 'next/image';

const PUBLIC_NAV = [
    { href: '/', label: 'Рецепты' },
    { href: '/search', label: 'Поиск по тегам' },
];

export function Header() {
    const pathname = usePathname();
    const router = useRouter();
    const { user, status, logoutLocal } = useAuthStore();
    const signOut = useSignOut();
    const [menuOpen, setMenuOpen] = useState(false);

    const isAuthed = status === 'authenticated' && !!user;

    const isActive = (href: string) =>
        href === '/' ? pathname === '/' : pathname.startsWith(href);

    const handleSignOut = () => {
        signOut.mutate(undefined, {
            onSettled: () => {
                logoutLocal();
                router.push('/login');
            },
        });
    };

    return (
        <header className="sticky top-0 z-40 border-b border-warm-brown/10 bg-cream/95 backdrop-blur">
            <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3">
                <Link href="/" className="flex items-center gap-2">
                    <Image
                        src="/icon.png"
                        alt="CookBook"
                        width={32}
                        height={32}
                        className="h-8 w-8"
                        priority
                    />
                    <span className="font-serif text-2xl font-bold tracking-tight text-ink">
                        CookBook
                    </span>
                </Link>

                <nav className="hidden items-center gap-1 md:flex">
                    {PUBLIC_NAV.map((n) => (
                        <Link
                            key={n.href}
                            href={n.href}
                            className={`rounded-lg px-3 py-1.5 text-sm transition-colors ${isActive(n.href)
                                ? 'bg-warm-brown/10 font-medium text-ink'
                                : 'text-warm-brown hover:bg-warm-brown/5'
                                }`}
                        >
                            {n.label}
                        </Link>
                    ))}

                    {isAuthed && (
                        <>
                            <Link
                                href="/my-recipes"
                                className={`rounded-lg px-3 py-1.5 text-sm transition-colors ${isActive('/my-recipes')
                                    ? 'bg-warm-brown/10 font-medium text-ink'
                                    : 'text-warm-brown hover:bg-warm-brown/5'
                                    }`}
                            >
                                Мои рецепты
                            </Link>
                        </>
                    )}
                </nav>

                <div className="hidden items-center gap-2 md:flex">
                    {isAuthed ? (
                        <>
                            <Link href="/recipes/new">
                                <Button size="sm">+ Рецепт</Button>
                            </Link>
                            <Link
                                href="/profile"
                                className="ml-1 text-sm text-warm-brown hover:text-ink"
                            >
                                {user.name}
                            </Link>
                            {user.role === 'Admin' && (
                                <Link href="/admin">
                                    <Button variant="secondary" size="sm">Админка</Button>
                                </Link>
                            )}
                            <Button variant="ghost" size="sm" onClick={handleSignOut} loading={signOut.isPending}>
                                Выйти
                            </Button>
                        </>
                    ) : status === 'guest' ? (
                        <>
                            <Link href="/login"><Button variant="ghost" size="sm">Войти</Button></Link>
                            <Link href="/register"><Button size="sm">Регистрация</Button></Link>
                        </>
                    ) : null}
                </div>

                <button
                    className="rounded-lg p-2 text-ink md:hidden"
                    onClick={() => setMenuOpen((v) => !v)}
                    aria-label="Меню"
                >
                    ☰
                </button>
            </div>

            {menuOpen && (
                <div className="border-t border-warm-brown/10 bg-cream px-4 py-3 md:hidden">
                    <div className="flex flex-col gap-1">
                        {PUBLIC_NAV.map((n) => (
                            <Link
                                key={n.href}
                                href={n.href}
                                onClick={() => setMenuOpen(false)}
                                className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5"
                            >
                                {n.label}
                            </Link>
                        ))}
                        {isAuthed ? (
                            <>
                                <Link href="/my-recipes" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">Мои рецепты</Link>
                                <Link href="/recipes/new" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">+ Рецепт</Link>
                                <Link href="/profile" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">Профиль ({user.name})</Link>
                                {user.role === 'Admin' && (
                                    <Link href="/admin" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">Админка</Link>
                                )}
                                <button onClick={() => { setMenuOpen(false); handleSignOut(); }} className="rounded-lg px-3 py-2 text-left text-sm text-warm-brown hover:bg-warm-brown/5">Выйти</button>
                            </>
                        ) : status === 'guest' ? (
                            <>
                                <Link href="/login" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">Войти</Link>
                                <Link href="/register" onClick={() => setMenuOpen(false)} className="rounded-lg px-3 py-2 text-sm text-warm-brown hover:bg-warm-brown/5">Регистрация</Link>
                            </>
                        ) : null}
                    </div>
                </div>
            )}
        </header>
    );
}