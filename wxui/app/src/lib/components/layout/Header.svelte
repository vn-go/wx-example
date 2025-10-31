<script lang="ts">
	import { accessToken, clearAccessToken } from '@store/auth';
	import { toggleSidebar } from '@store/uiStore';
	import { Menu } from 'lucide-svelte';
	import { onMount } from 'svelte';
	import { get, writable } from 'svelte/store';

	let currentTheme: 'light' | 'dark';

	let showMenu = false;
	let userName = 'User';

	onMount(() => {
		const token = get(accessToken);
		if (token) {
			// nếu token chứa tên người dùng hoặc jwt, có thể decode ở đây
			userName = 'Admin'; // ví dụ hardcode, có thể thay bằng decode(token)
		}
	});

	function toggleMenu() {
		showMenu = !showMenu;
	}

	function signOut() {
		clearAccessToken();
		showMenu = false;
		location.reload(); // hoặc emit event "logout"
	}
	const collapsed = writable(true);
</script>

<header class="flex items-center justify-between px-1 py-1 bg-gray-100 shadow">
	<h1 class="text-xl font-semibold">
		<button class="z-60 p-2 bg-white" on:click={toggleSidebar} aria-label="Toggle sidebar">
			<Menu size="20" />
		</button>
		hello
	</h1>

	<div class="relative">
		<button class="flex items-center gap-2 py-1 hover:bg-gray-100">
			<img
				src="https://ui-avatars.com/api/?name={userName}&background=random"
				alt="Avatar"
				class="w-8 h-8 rounded-full"
			/>
			<span class="font-medium">{userName}</span>
			<svg
				class="w-4 h-4"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				xmlns="http://www.w3.org/2000/svg"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"
				></path>
			</svg>
		</button>

		{#if showMenu}
			<div
				class="absolute right-0 mt-2 w-40 bg-white border border-gray-200 rounded-lg shadow-lg z-50"
			>
				<button
					class="block w-full text-left px-4 py-2 text-sm hover:bg-gray-100"
					on:click={signOut}
				>
					Sign out
				</button>
			</div>
		{/if}
	</div>
</header>

<style>
	header {
		z-index: 100;
	}
</style>
