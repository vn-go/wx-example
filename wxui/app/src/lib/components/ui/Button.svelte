<script>
	export let className = '';
</script>

<button
	class={`bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition ${className}`}
>
	<slot />
</button>
