export function Spinner({ className = '' }: { className?: string }) {
    return (
        <div className={`flex items-center justify-center py-10 ${className}`}>
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-terracotta/30 border-t-terracotta" />
        </div>
    );
}

export function FullScreenSpinner() {
    return (
        <div className="flex min-h-[60vh] items-center justify-center">
            <div className="h-10 w-10 animate-spin rounded-full border-4 border-terracotta/30 border-t-terracotta" />
        </div>
    );
}