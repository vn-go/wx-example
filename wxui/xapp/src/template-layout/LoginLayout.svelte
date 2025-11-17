<!-- src/lib/components/LoginLayout.svelte -->
<script lang="ts">
    // Không cần props – chỉ layout
</script>

<div
    class="h-full bg-gradient-to-br from-indigo-50 via-white to-purple-50 flex items-center justify-center p-4"
>
    <div class="w-full max-w-md bg-white rounded-2xl shadow-xl p-8 space-y-6">
        <!-- Logo / Title (optional) -->
        <div class="text-center mb-6">
            <h1 class="text-3xl font-bold text-gray-800">Welcome Back</h1>
            <p class="text-gray-500 text-sm mt-1">Sign in to your account</p>
        </div>

        <!-- Username Slot -->
        <div class="space-y-1">
            <label class="block text-sm font-medium text-gray-700"
                >Username</label
            >
            <slot name="username-input">
                <div class="text-gray-400 italic text-sm">
                    No username input
                </div>
            </slot>
        </div>

        <!-- Password Slot -->
        <div class="space-y-1">
            <label class="block text-sm font-medium text-gray-700"
                >Password</label
            >
            <slot name="password-input">
                <div class="text-gray-400 italic text-sm">
                    No password input
                </div>
            </slot>
        </div>

        <!-- Error Message Slot -->
        <div class="min-h-6">
            <slot name="error-message"></slot>
        </div>

        <!-- Submit Button Slot -->
        <div class="pt-4">
            <slot name="submit-button">
                <button
                    class="w-full py-3 bg-indigo-600 text-white rounded-xl font-medium hover:bg-indigo-700 transition"
                >
                    Sign In
                </button>
            </slot>
        </div>

        <!-- Footer -->
        <div class="text-center text-xs text-gray-500">
            <slot name="footer">
                <p>© 2025 Your App. All rights reserved.</p>
            </slot>
        </div>
    </div>
</div>

<style>
    /* Optional: focus ring */
    input:focus {
        @apply ring-2 ring-indigo-500 ring-offset-2 outline-none;
    }
</style>
