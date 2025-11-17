<!-- src/lib/components/Error.svelte -->
<script lang="ts">
    // Props
    export let code: string = "Error";
    export let title: string = "Load Failed";
    export let message: string =
        "Unable to load the component. Please try again.";
    export let retry: (() => void) | null = null;
    export let showRetry: boolean = true;
</script>

<div
    class="p-6 bg-red-50 border border-red-200 rounded-xl text-center h-full w-full"
>
    <!-- Icon -->
    <div
        class="mx-auto w-16 h-16 mb-4 flex items-center justify-center bg-red-100 rounded-full"
    >
        <svg
            class="w-10 h-10 text-red-600"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
            />
        </svg>
    </div>

    <!-- Title -->
    <h3 class="text-lg font-bold text-red-800 mb-1">
        {code}
    </h3>
    <h4 class="text-xl font-semibold text-red-700 mb-2">
        {title}
    </h4>

    <!-- Message -->
    <p class="text-red-600 text-sm mb-4 max-w-xs mx-auto">
        {message}
    </p>

    <!-- Retry Button -->
    {#if showRetry && retry}
        <button
            on:click={retry}
            class="px-5 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 hover:shadow-md transform hover:scale-105 transition-all duration-200"
        >
            Try Again
        </button>
    {/if}
</div>

<style>
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(8px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
    .animate-fade-in {
        animation: fadeIn 0.4s ease-out;
    }
</style>
