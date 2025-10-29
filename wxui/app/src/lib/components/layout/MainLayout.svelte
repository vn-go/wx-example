<script>
	import Footer from './Footer.svelte';
	import Header from './Header.svelte';
	import Sidebar from './Sidebar.svelte';
</script>

<div class="flex h-screen bg-gray-50 text-gray-900">
	<Sidebar />
	<div class="flex flex-col flex-1">
		<Header />
		<main class="flex-1 overflow-y-auto p-6">
			<slot />
		</main>
		<Footer />
	</div>
</div>
