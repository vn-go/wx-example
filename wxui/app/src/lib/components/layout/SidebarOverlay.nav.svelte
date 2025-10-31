<script lang="ts">
	import { Group, User } from 'lucide-svelte';

	type MenuItem = {
		id: number;
		title: string;
		icon: any; // component icon
		link: string;
	};

	// Giả lập dữ liệu menu từ API
	const menuData: MenuItem[] = [
		{ id: 1, title: 'Users', icon: User, link: '/users' },
		{ id: 2, title: 'Roles', icon: Group, link: '/roles' }
	];
</script>

<nav class="flex-1 mt-2 space-y-1">
	{#each menuData as item}
		<a
			href={item.link}
			class="flex items-center gap-3 px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-md transition-colors"
		>
			<svelte:component this={item.icon} size="18" />
			<span>{item.title}</span>
		</a>
	{/each}
</nav>
