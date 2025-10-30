<script>
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();
	export let className = '';
</script>

<button
	on:click={(e) => dispatch('click', e)}
	class={`bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition ${className}`}
>
	<slot />
</button>
