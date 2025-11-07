<script lang="ts">
	import Container from '@components/ui/Container.svelte';
	import { apiCall } from '@lib/utils/apis';
	import { onMount } from 'svelte';
	class UIForm {
		viewPath: string;
		OnMounth(callback: () => void) {
			onMount(callback);
		}
		async PostData(apiEnpoint: string, data: any): Promise<any> {
			let response = await apiCall.post(apiEnpoint, data, {
				'View-Path': this.viewPath
			});
			return response.data as any;
		}
		constructor(viewPath: string) {
			this.viewPath = viewPath;
		}
	}
	class User extends UIForm {}
	const user = new User('system/user');

	let items: any[] = [];

	async function loadMore() {
		const start = items.length;
		var newItems = await user.PostData('accounts/get-list-of-accounts', {
			index: start,
			size: 20,
			orderBy: ['username']
		});

		items = [...items, ...(newItems || [])];
	}
	onMount(async () => {
		var newItems = await user.PostData('accounts/get-list-of-accounts', {
			index: 0,
			size: 20,
			orderBy: ['username']
		});

		items = [...items, ...(newItems || [])];
	});
	function formatKey(key: string): string {
		return key
			.replace(/([A-Z])/g, ' $1')
			.replace(/[_-]+/g, ' ')
			.replace(/\b\w/g, (c) => c.toUpperCase())
			.trim();
	}
	// import Container from '@components/ui/Container.svelte';
</script>

<div class="flex h-full w-full">
	<!-- <div class="w-full h-full bg-blue-500" style="overflow-y: auto;">
		<div style="height: 100000px;">dasdsad</div>
	</div> -->
	<!-- <Container /> -->
	<Container {items} {loadMore} columnMinWidth={350} keyName="userId">
		<!-- Dùng svelte:fragment để nhận slot props -->
		<svelte:fragment let:item>
			<div
				class="bg-white p-4 rounded-lg shadow-sm border border-transparent
           transform transition duration-300 ease-in-out
           hover:shadow-lg hover:-translate-y-1 hover:scale-[1.01]
           focus:outline-none focus:ring-2 focus:ring-indigo-300
           active:translate-y-0
           overflow-hidden"
				tabindex="0"
				role="article"
				aria-label={`User ${item.username ?? ''}`}
			>
				{#each Object.entries(item) as [key, value]}
					<div class="grid grid-cols-2 gap-2 py-1">
						<div class="text-sm font-medium text-gray-600 capitalize">{formatKey(key)}:</div>
						<div class="text-sm text-gray-800 break-words">{value ?? '—'}</div>
					</div>
				{/each}
			</div>
		</svelte:fragment>
	</Container>
</div>
