<script lang="ts">
	import { FormInput, Group, User } from 'lucide-svelte';
	import { writable } from 'svelte/store';
	const openMenus = writable<Record<number, boolean>>({});

	function toggleMenu(id: number) {
		openMenus.update((m) => ({ ...m, [id]: !m[id] }));
	}
	type MenuItem = {
		id: number;
		title: string;
		icon: any; // component icon
		link: string;
		children?: MenuItem[];
	};

	// Giả lập dữ liệu menu từ API
	const menuData: MenuItem[] = [
		{
			id: 0,
			title: 'System',
			icon: Group,
			link: '/system',
			children: [
				{ id: 1, title: 'Users', icon: User, link: '/system/users' },
				{ id: 1, title: 'views', icon: FormInput, link: '/system/views' },
				{ id: 2, title: 'Roles', icon: Group, link: '/system/roles' }
			]
		}
	];
</script>

<nav class="flex-1 mt-2 space-y-1">
	<ul>
		{#each menuData as item (item.id)}
			<li>
				<div on:click={() => item.children && toggleMenu(item.id)}>
					<svelte:component this={item.icon} />
					{item.title}
					{#if item.children}
						<span
							>{#if $openMenus[item.id]}▼{:else}▶{/if}</span
						>
					{/if}
				</div>

				{#if item.children && $openMenus[item.id]}
					<ul class="submenu">
						{#each item.children as child (child.id)}
							<li>
								<svelte:component this={child.icon} />
								<a href={child.link}>{child.title}</a>
							</li>
						{/each}
					</ul>
				{/if}
			</li>
		{/each}
	</ul>
</nav>
