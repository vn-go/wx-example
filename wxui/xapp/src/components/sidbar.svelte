<script lang="ts">
    import { BaseUi } from "@core/core";
    import { writable } from "svelte/store";

    class SideBar extends BaseUi {
        menuData: any;

        constructor() {
            super();
            this.menuData = [
                {
                    title: "system",
                    children: [
                        {
                            path: "/system/users",
                            title: "Users",
                        },
                        {
                            path: "/system/roles",
                            title: "Roles",
                        },
                    ],
                },
            ];
        }
    }
    let sideBar = new SideBar();
    const openGroups = writable<Record<string, boolean>>({});
    function toggleGroup(title: string) {
        openGroups.update((curr) => ({ ...curr, [title]: !curr[title] }));
    }
</script>

<div>
    <h3 class="font-semibold text-lg mb-4">Menu</h3>
    <ul class="space-y-2">
        {#each sideBar.menuData as group}
            <li>
                <button
                    on:click={() => toggleGroup(group.title)}
                    class="flex justify-between items-center w-full p-2 rounded hover:bg-gray-100 transition"
                >
                    <span>{group.title}</span>
                    <svg
                        class="w-4 h-4 transform transition-transform duration-200"
                        class:rotate-90={$openGroups[group.title]}
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M9 5l7 7-7 7"
                        ></path>
                    </svg>
                </button>

                {#if $openGroups[group.title]}
                    <ul class="pl-4 mt-1 space-y-1">
                        {#each group.children as item}
                            <li>
                                <a
                                    href="javascript:void(0);"
                                    on:click={() =>
                                        sideBar.jumpToUrl(item.path)}
                                    class="block p-2 rounded hover:bg-gray-200 transition"
                                >
                                    {item.title}
                                </a>
                            </li>
                        {/each}
                    </ul>
                {/if}
            </li>
        {/each}
    </ul>
</div>
