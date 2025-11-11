<script lang="ts">
	import { ChevronDown, ChevronRight, FormInput, Group, User } from 'lucide-svelte';
	import { writable } from 'svelte/store';

	const openMenus = writable<Record<number, boolean>>({});

	function toggleMenu(id: number) {
		openMenus.update((m) => ({ ...m, [id]: !m[id] }));
	}

	type MenuItem = {
		id: number;
		title: string;
		icon: any;
		link: string;
		children?: MenuItem[];
	};

	const menuData: MenuItem[] = [
		{
			id: 1,
			title: 'System',
			icon: Group,
			link: '/system',
			children: [
				{ id: 11, title: 'Users', icon: User, link: '/system/users' },
				{ id: 12, title: 'Views', icon: FormInput, link: '/system/views' },
				{ id: 13, title: 'Roles', icon: Group, link: '/system/roles' }
			]
		},
		{
			id: 2,
			title: 'Settings',
			icon: FormInput,
			link: '/settings',
			children: [
				{ id: 21, title: 'General', icon: User, link: '/settings/general' },
				{ id: 22, title: 'Advanced', icon: FormInput, link: '/settings/advanced' }
			]
		}
	];
</script>

<nav class="flex-1 mt-2 space-y-1">
	<ul class="space-y-1">
		{#each menuData as item (item.id)}
			<li class="select-none">
				<!-- Menu chính -->
				<div class="sidebar-list-item debug" on:click={() => toggleMenu(item.id)}>
					<div class="flex items-center gap-2">
						<svelte:component this={item.icon} size={16} />
						<span>{item.title}</span>
					</div>
					{#if item.children}
						<svelte:component this={$openMenus[item.id] ? ChevronDown : ChevronRight} size={16} />
					{/if}
				</div>

				<!-- Submenu -->
				{#if item.children && $openMenus[item.id]}
					<ul class="ml-4 space-y-1 border-l-2 border-gray-950">
						{#each item.children as child (child.id)}
							<li>
								<a href={child.link} class="sidebar-list-item">
									<svelte:component this={child.icon} size={14} />
									<span>{child.title}</span>
								</a>
							</li>
						{/each}
					</ul>
				{/if}
			</li>
		{/each}
	</ul>
</nav>

<style>
	.submenu {
		transition: all 0.3s ease;
	}

	a {
		text-decoration: none;
		color: inherit;
	}
</style>
