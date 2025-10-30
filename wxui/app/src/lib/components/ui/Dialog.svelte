<script lang="ts">
	import type { SvelteComponentTyped } from 'svelte';
	import { createEventDispatcher } from 'svelte';

	export let component: typeof SvelteComponentTyped | null = null;
	export let props: Record<string, any> = {};

	const dispatch = createEventDispatcher();

	function close() {
		dispatch('close');
	}
</script>

{#if component}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
		<div class="bg-white w-[90vw] max-w-md rounded-xl shadow-xl p-4 relative">
			<button class="absolute top-2 right-2 text-gray-500 text-xl" on:click={close}>
				&times;
			</button>
			<svelte:component this={component} {...props} />
		</div>
	</div>
{:else}
	<div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
		<div class="bg-white w-[90vw] max-w-md rounded-xl shadow-xl p-4 relative">
			<button class="absolute top-2 right-2 text-gray-500 text-xl" on:click={close}>
				&times;
			</button>
			Lỗi load component
		</div>
	</div>
{/if}
