import Link from 'next/link';
import { Button } from '@/components/ui/Button';

export default function NotFound() {
    return (
        <div className="flex flex-col items-center justify-center gap-4 py-24 text-center">
            <span className="text-6xl">🍽️</span>
            <h1 className="font-serif text-3xl text-warm-brown">Страница не найдена</h1>
            <p className="max-w-md text-sm text-warm-brown/60">
                Возможно, ссылка устарела или рецепт был удалён.
            </p>
            <Link href="/"><Button>На главную</Button></Link>
        </div>
    );
}