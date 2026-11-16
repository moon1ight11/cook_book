import { Suspense } from 'react';
import RecipeDetailClient from './RecipeDetailClient';
import { FullScreenSpinner } from '@/components/ui/Spinner';

async function RecipeDetailContent({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  return <RecipeDetailClient id={id} />;
}

export default function RecipeDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  return (
    <Suspense fallback={<FullScreenSpinner />}>
      <RecipeDetailContent params={params} />
    </Suspense>
  );
}