
<script lang="ts">
    import BaseComponent from "@core/base.component";
    
    class Header extends BaseComponent {
        burgerOpen = $state(false);
        dropdownOpen=$state(false);
        constructor() {
            super();
        }
        toggleBurger(){
            this.burgerOpen = !this.burgerOpen;
            this.emit('toggle-burger');
        }
        onDestroy(callback: () => void): void {
            this.off('toggle-burger');
        }
        toggleDropdown(){
            this.dropdownOpen = !this.dropdownOpen;
        }

    }
    let header = $state(new Header()) ; // how to use header in another component?
    export { header };
   
    
</script>
<header class="main-header-menu">
    <div class="flex items-center gap-2">
    <button
    
    aria-label="Toggle sidebar"
    onclick={()=>{
        
        header.toggleBurger();
    }}
    aria-expanded={header.burgerOpen}
    >
    <div class="space-y-1 burger {header.burgerOpen ? 'open' : ''}">
    <span class="block w-5 h-0.5 bg-current transition-transform {header.burgerOpen ? 'translate-y-1.5 rotate-45' : ''}"></span>
    <span class="block w-5 h-0.5 bg-current transition-opacity {header.burgerOpen ? 'opacity-0' : ''}"></span>
    <span class="block w-5 h-0.5 bg-current transition-transform {header.burgerOpen ? '-translate-y-1.5 -rotate-45' : ''}"></span>
    </div>
    </button>
    </div>
    
    
    <div class="flex-1 flex justify-center">
    <h1 class="font-semibold tracking-wide"></h1>
    </div>
    
    
    <div class="relative header-dropdown">
    <button
    class="p-2 rounded-md hover:bg-gray-800 focus:outline-none"
    aria-haspopup="true"
    aria-expanded={header.burgerOpen}
    onclick={()=>{
        debugger;
        header.toggleDropdown();
    }}
    >
    Menu ▾
    </button>
    
    
    {#if header.dropdownOpen}
    <div
    class="absolute right-0 mt-2 w-40 bg-white text-gray-800 rounded-lg shadow-lg overflow-hidden z-40"
    role="menu"
    >
    <a href="#" class="block px-4 py-2 text-sm hover:bg-gray-100" role="menuitem">Profile</a>
    <a href="#" class="block px-4 py-2 text-sm hover:bg-gray-100" role="menuitem">Settings</a>
    <button
    class="w-full text-left px-4 py-2 text-sm hover:bg-gray-100"
    role="menuitem"
    onclick={() => alert('Signed out')}
    >
    Sign out
    </button>
    </div>
    {/if}
    </div>
    </header>