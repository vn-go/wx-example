<script>
	import { page } from '$app/stores';

	let isOpen = true; // desktop mặc định mở
	let isMobileOpen = false; // mobile menu

	$: currentPath = $page.url.pathname;

	function toggleSidebar() {
		isOpen = !isOpen;
	}

	function toggleMobile() {
		isMobileOpen = !isMobileOpen;
	}
</script>

<!-- Wrapper -->
<div class="flex">
	<!-- Mobile topbar -->
	<div class="md:hidden flex items-center justify-between bg-gray-800 text-white p-4">
		<span class="font-bold">Admin</span>
		<button on:click={toggleMobile} class="p-2 rounded hover:bg-gray-700"> ☰ </button>
	</div>

	<!-- Sidebar -->
	<aside
		class="
    transition-all duration-300 bg-gray-800 text-white h-screen fixed md:static z-50
    md:translate-x-0
  "
		class:w-64={isOpen}
		class:w-16={!isOpen}
		class:translate-x-0={isMobileOpen}
		class:-translate-x-full={!isMobileOpen}
	>
		<!-- Header -->
		<div class="hidden md:flex items-center justify-between p-4 border-b border-gray-700">
			{#if isOpen}
				<span class="text-xl font-bold">Admin</span>
			{/if}
			<button class="p-1 rounded hover:bg-gray-700" on:click={toggleSidebar} title="Toggle">
				{#if isOpen}
					⏪
				{:else}
					⏩
				{/if}
			</button>
		</div>

		<!-- Navigation -->
		<nav class="p-2 space-y-2">
			<a
				href="/"
				class="flex items-center px-3 py-2 rounded hover:bg-gray-700"
				class:bg-blue-600={currentPath === '/'}
			>
				<span>🏠</span>
				{#if isOpen}<span class="ml-2">Dashboard</span>{/if}
			</a>
			<a
				href="/users"
				class="flex items-center px-3 py-2 rounded hover:bg-gray-700"
				class:bg-blue-600={currentPath.startsWith('/users')}
			>
				<span>👤</span>
				{#if isOpen}<span class="ml-2">Users</span>{/if}
			</a>
			<a
				href="/settings"
				class="flex items-center px-3 py-2 rounded hover:bg-gray-700"
				class:bg-blue-600={currentPath.startsWith('/settings')}
			>
				<span>⚙️</span>
				{#if isOpen}<span class="ml-2">Settings</span>{/if}
			</a>
		</nav>
	</aside>

	<!-- Overlay khi mở mobile menu -->
	{#if isMobileOpen}
		<div class="fixed inset-0 bg-black/50 md:hidden" on:click={toggleMobile}></div>
	{/if}
</div>
