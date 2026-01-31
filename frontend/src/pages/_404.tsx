import { useLocation } from "preact-iso";

export function NotFound() {
    const { route } = useLocation();

    return (
        <div class="min-h-screen flex items-center justify-center bg-gray-50">
            <div class="text-center">
                <h1 class="text-6xl font-bold text-gray-800 mb-4">404</h1>
                <p class="text-xl text-gray-600 mb-8">Page not found</p>
                <button
                    onClick={() => route("/")}
                    class="text-blue-600 hover:underline cursor-pointer bg-transparent border-none text-base">
                    Go back home
                </button>
            </div>
        </div>
    );
}
