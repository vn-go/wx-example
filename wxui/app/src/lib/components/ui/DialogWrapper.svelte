<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	export let component: any = null;
	export let props: Record<string, any> = {};

	const dispatch = createEventDispatcher();

	function close() {
		dispatch('close');
	}
</script>

{#if component}
	<div class="fixed inset-0 flex items-center justify-center bg-black/50 z-50">
		<svelte:component this={component} {...props} on:close={close} />
	</div>
{/if}
