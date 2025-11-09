<script lang="ts">
	import Layout from '@layouts/form.svelte';
</script>

<Layout title="My Page">
	<h2 slot="header">Custom Header Content</h2>

	<p slot="main">This is the main content passed from App.svelte</p>

	<ul slot="sidebar">
		<li>Item 1</li>
		<li>Item 2</li>
		<li>Item 3</li>
	</ul>

	<small slot="footer">© 2025 My Company</small>
</Layout>
