<script lang="ts">
	export let items: any[] = []; // danh sách data
	export let loadMore: () => void; // function lazy load
	export let columnMinWidth = 250; // min-width mỗi cột (px)
	export let className = '';
	export let keyName = 'id'; // key của mỗi item

	let container: HTMLDivElement;

	function handleScroll() {
		if (!container) return;

		const { scrollTop, scrollHeight, clientHeight } = container;
		if (scrollTop + clientHeight >= scrollHeight - 50) {
			loadMore?.();
		}
	}
</script>

<div
	bind:this={container}
	class="overflow-y-auto h-full w-full p-1 {className}"
	on:scroll={handleScroll}
>
	<div
		class="grid gap-3"
		style="grid-template-columns: repeat(auto-fill, minmax({columnMinWidth}px, 1fr))"
	>
		{#each items as item ((item as any)[keyName])}
			<slot {item}></slot>
		{/each}
	</div>
</div>
