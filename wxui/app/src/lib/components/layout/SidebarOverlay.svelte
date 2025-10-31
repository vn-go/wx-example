<script lang="ts">
	import SidebarNav from '@components/layout/SidebarOverlay.nav.svelte';
	import { theme } from '@store/themeStore';
	import { sidebarCollapsed, toggleSidebar } from '@store/uiStore';
	import { LogOut, Moon, Sun } from 'lucide-svelte';
	import { onDestroy, onMount } from 'svelte';
	import { get, writable } from 'svelte/store';
	const collapsed = writable(true); // true = ẩn sidebar
	let showSidebar: boolean;

	collapsed.subscribe((v) => (showSidebar = !v));

	function toggleSidebarUpdate() {
		collapsed.update((v) => !v);
		toggleSidebar();
	}

	function toggleTheme() {
		theme.update((t) => (t === 'light' ? 'dark' : 'light'));
	}
	let sidebarRef: HTMLElement;

	function handleMouseMove(e: MouseEvent) {
		if (!sidebarRef) return;
		const rect = sidebarRef.getBoundingClientRect();
		// Nếu con trỏ nằm ngoài sidebar
		if (
			(e.clientX < rect.left ||
				e.clientX > rect.right ||
				e.clientY < rect.top ||
				e.clientY > rect.bottom) &&
			get(collapsed) == false
		) {
			//toggleSidebarUpdate(); // collapse sidebar / hide backdrop
			sidebarCollapsed.set(true); // collapse sidebar / hide backdrop
			collapsed.set(true); // collapse sidebar / hide backdrop
		}
	}
	onMount(() => {
		document.addEventListener('mousemove', handleMouseMove);
	});

	onDestroy(() => {
		document.removeEventListener('mousemove', handleMouseMove);
	});
</script>

<!-- Backdrop -->
<!-- {#if !collapsed} -->
<div
	class="fixed inset-0 bg-opacity-30 z-40"
	on:mousemove={() => {
		sidebarCollapsed.set(true); // collapse sidebar / hide backdrop
		collapsed.set(true); // collapse sidebar / hide backdrop
	}}
	style="display: {$sidebarCollapsed ? 'none' : 'block'};"
></div>
<!-- {/if} -->

<!-- Sidebar -->
<div
	bind:this={sidebarRef}
	class="fixed left-0 w-64 bg-white dark:bg-gray-900 text-gray-800 dark:text-gray-100 shadow-xl z-50 flex flex-col transition-transform duration-300"
	style="
    top: 3rem; /* thấp xuống dưới header */
    height: calc(100% - 3rem); /* chiều cao giảm đi phần header */
    transform: translateX({$sidebarCollapsed ? '-100%' : '0'});
  "
>
	<!-- Menu -->
	<SidebarNav />

	<!-- Footer -->
	<div class="mt-auto border-t dark:border-gray-800 px-4 py-3 flex flex-col gap-2">
		<button
			on:click={toggleTheme}
			class="flex items-center gap-3 hover:bg-gray-100 dark:hover:bg-gray-800 w-full px-2 py-2 rounded-md transition-colors"
		>
			{#if $theme === 'dark'}
				<Sun size="18" />
				<span>Light Mode</span>
			{:else}
				<Moon size="18" />
				<span>Dark Mode</span>
			{/if}
		</button>
		<button
			class="flex items-center gap-3 hover:bg-gray-100 dark:hover:bg-gray-800 w-full px-2 py-2 rounded-md transition-colors"
		>
			<LogOut size="18" />
			<span>Sign out</span>
		</button>
	</div>
</div>

<!-- Toggle button ghim -->
