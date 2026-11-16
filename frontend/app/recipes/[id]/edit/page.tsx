import { Suspense } from 'react';
import EditRecipeClient from './EditRecipeClient';
import { FullScreenSpinner } from '@/components/ui/Spinner';

async function EditContent({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  return <EditRecipeClient id={id} />;
}

export default function EditRecipePage({ params }: { params: Promise<{ id: string }> }) {
  return (
    <Suspense fallback={<FullScreenSpinner />}>
      <EditContent params={params} />
    </Suspense>
  );
}