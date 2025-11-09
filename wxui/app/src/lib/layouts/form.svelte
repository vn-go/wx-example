<script>
	export let title = 'Default Title';
</script>

<div class="layout">
	<header>
		<slot name="header">
			<!-- default header content -->
			<h1>{title}</h1>
		</slot>
	</header>

	<main>
		<slot name="main">
			<!-- default main content -->
			<p>Main content goes here</p>
		</slot>
	</main>

	<aside>
		<slot name="sidebar">
			<!-- default sidebar content -->
			<p>Sidebar content</p>
		</slot>
	</aside>

	<footer>
		<slot name="footer">
			<!-- default footer content -->
			<p>Footer</p>
		</slot>
	</footer>
</div>

<style>
	.layout {
		display: grid;
		grid-template-areas:
			'header header'
			'main sidebar'
			'footer footer';
		grid-template-columns: 3fr 1fr;
		gap: 1rem;
	}
	header {
		grid-area: header;
		background: #eee;
		padding: 1rem;
	}
	main {
		grid-area: main;
		background: #fafafa;
		padding: 1rem;
	}
	aside {
		grid-area: sidebar;
		background: #f0f0f0;
		padding: 1rem;
	}
	footer {
		grid-area: footer;
		background: #ddd;
		padding: 1rem;
	}
</style>
