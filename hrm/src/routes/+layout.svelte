<script>
	import Sidebar from '$lib/components/Sidebar.svelte';
	import '../app.css';
	import { loadLocale } from '$lib/i18n';
  import { locale, t } from 'svelte-i18n';
</script>

<div class="min-h-screen flex">
	<Sidebar />
	<main class="flex-1 bg-gray-100 p-6">
		<slot />
	</main>
</div>
