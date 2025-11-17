<!-- src/lib/components/Layout.svelte -->
<script lang="ts">
    import { BaseUi } from "@core/core";
    class Layout extends BaseUi {
        headerEl: HTMLElement | undefined;
        sidebarEl: HTMLElement | undefined;
        bodyEl: HTMLElement | undefined;
        burgerEl: HTMLElement | undefined;
        isSidebarOpen: boolean | undefined;
        screen: HTMLElement | undefined;
        self: Layout | undefined;
        slotBody: HTMLElement | undefined;
        onBugerClick() {}
        async onInit() {
            if (this.screen) {
                this.slotBody = await this.findBySlotName(this.screen, "body");
                this.fixBodySize();
                const self = this; // Lưu lại this

                window.addEventListener("resize", () => {
                    self.fixBodySize();
                });
            }

            this.isntallBugerEvent();
        }
        /**
         * Adjusts the height of the body slot to fit the viewport minus the header height.
         * Uses `this.slotBody` and `this.headerEl` which are assumed to be bound HTMLElement references.
         */
        fixBodySize() {
            // Get references (already bound via bind:this in Svelte)
            const slotBody = this.slotBody;
            const headerEl = this.headerEl;

            // Skip if elements are not available
            if (!slotBody || !headerEl) {
                console.warn("[fixBodySize] Missing slotBody or headerEl");
                return;
            }

            // Get header height (including padding, border, margin)
            const headerHeight = headerEl.offsetHeight;

            // Get viewport height (full browser window height)
            const viewportHeight = window.innerHeight;

            // Calculate and set body height
            const bodyHeight = viewportHeight - headerHeight;
            slotBody.style.height = `${bodyHeight}px`;
            slotBody.style.width = `${window.innerWidth}px`;
            // Optional: Ensure content doesn't overflow
            slotBody.style.overflow = "auto";

            // Debug log
            console.log(
                `[fixBodySize] Header: ${headerHeight}px | Viewport: ${viewportHeight}px | Body: ${bodyHeight}px`,
            );
        }
        isntallBugerEvent() {
            // if (!this.burgerEl) return;
            // (this.burgerEl as any).addEventListener(
            //     "click",
            //     this.toggleSidebar,
            // );
        }

        toggleSidebar() {
            const self = this;
            console.log(this);
            function hideSideBar() {
                self.isSidebarOpen = false;
                self.sidebarEl?.classList.remove("translate-x-0");
                self.sidebarEl?.classList.add("-translate-x-full");
            }
            if (!this.sidebarEl) return;
            this.isSidebarOpen = !this.isSidebarOpen;
            if (this.isSidebarOpen) {
                this.sidebarEl.classList.remove("-translate-x-full");
                this.sidebarEl.classList.add("translate-x-0");
                this.slotBody?.addEventListener("mousemove", hideSideBar);
            } else {
                this.sidebarEl.classList.remove("translate-x-0");
                this.sidebarEl.classList.add("-translate-x-full");
                this.slotBody?.removeEventListener("mousemove", hideSideBar);
            }
        }
    }
    const layout = new Layout();
</script>

<div
    class="main min-h-screen flex flex-col bg-gray-50"
    bind:this={layout.screen}
>
    <!-- Header -->
    <header
        bind:this={layout.headerEl}
        class="h-11 bg-white border-b border-gray-200 flex items-center px-4 shadow-sm z-30"
    >
        <div class="flex items-center space-x-3">
            <!-- Burger icon (chỉ để hiển thị, không emit) -->
            <div
                bind:this={layout.burgerEl}
                on:click={() => {
                    layout.toggleSidebar();
                }}
                class="w-6 h-6 flex items-center justify-center cursor-pointer hover:bg-gray-100 rounded"
            >
                <svg
                    class="w-5 h-5 text-gray-600"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                >
                    <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M4 6h16M4 12h16M4 18h16"
                    />
                </svg>
            </div>
            <span class="text-sm font-medium hidden sm:inline text-gray-700"
                >My App</span
            >
        </div>
    </header>

    <!-- Main content area -->
    <div class="flex flex-1 relative shadow-lg rounded-xl">
        <!-- Sidebar -->
        <aside
            bind:this={layout.sidebarEl}
            style="position: fixed; top: 44px; height: calc(100vh - 44px);"
            class="sidebar w-64 bg-white border-r border-gray-200 z-20 transition-transform duration-300 ease-in-out -translate-x-full"
            class:open={false}
        >
            <div class="p-4">
                <slot name="sidebar">
                    <p class="text-gray-500 italic">Default Sidebar</p>
                </slot>
            </div>
        </aside>
        <slot name="body" class="debug">
            <div class="p-6">
                <p class="text-gray-500 italic">Default Body Content</p>
            </div>
        </slot>
        <!-- Body -->
        <!-- <main class="flex-1 pt-1 pl-0 md:pl-64 transition-all duration-300">
           
        </main> -->
    </div>
</div>

<style>
    /* Sidebar mở (nếu sau này bật) */
    .sidebar.open {
        transform: translateX(0);
    }
    .sidebar {
        transform: translateX(-100%);
    }
    @media (min-width: 768px) {
        .sidebar {
            transform: translateX(0) !important;
        }
    }
    .-translate-x-full {
        transform: translateX(-100%);
    }
</style>
