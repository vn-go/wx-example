<script lang="ts">
	import { Home, LogOut, Menu, Settings } from 'lucide-svelte';
	import { writable } from 'svelte/store';

	const collapsed = writable(false);

	function toggleSidebar() {
		collapsed.update((v) => !v);
	}
</script>

<div class="relative h-screen flex">
	<div
		class="fixed left-0 w-64 bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 shadow-xl z-50 flex flex-col transition-transform duration-300"
		style=" top: 0;     height: calc(100% - 0); transform: translateX({$collapsed
			? '-100%'
			: '0'});"
	>
		<!-- Header -->
		<div class="flex items-center justify-between px-4 py-3 border-b dark:border-gray-800">
			<span class="font-semibold text-lg" class:hidden={$collapsed}>My App</span>
			<button
				class="p-1 rounded-md hover:bg-gray-100 dark:hover:bg-gray-800"
				on:click={toggleSidebar}
				aria-label="Toggle sidebar"
			>
				<Menu size="20" />
			</button>
		</div>

		<!-- Menu -->
		<nav class="flex-1 mt-2 space-y-1">
			<a
				href="#"
				class="flex items-center gap-3 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors"
			>
				<Home size="18" />
				<span class:hidden={$collapsed}>Dashboard</span>
			</a>
			<a
				href="#"
				class="flex items-center gap-3 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors"
			>
				<Settings size="18" />
				<span class:hidden={$collapsed}>Settings</span>
			</a>
		</nav>

		<!-- Footer -->
		<div class="mt-auto border-t dark:border-gray-800 px-4 py-3 flex flex-col gap-2">
			<button
				class="flex items-center gap-3 hover:bg-gray-100 dark:hover:bg-gray-800 w-full px-2 py-2 rounded-md transition-colors"
			>
				<LogOut size="18" />
				<span class:hidden={$collapsed}>Sign out</span>
			</button>
		</div>
	</div>
</div>
