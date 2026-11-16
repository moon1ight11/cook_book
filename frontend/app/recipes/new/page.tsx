import { AuthGuard } from '@/components/layout/AuthGuard';
import { RecipeForm } from '@/features/recipes/RecipeForm';

export default function NewRecipePage() {
    return (
        <AuthGuard>
            <div className="mx-auto max-w-2xl">
                <h1 className="mb-6 font-serif text-3xl text-warm-brown">Новый рецепт</h1>
                <RecipeForm />
            </div>
        </AuthGuard>
    );
}