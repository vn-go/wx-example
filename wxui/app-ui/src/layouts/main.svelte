<script lang="ts">
  import '../app.css';
  export let title = "My App";
</script>

<div class="flex flex-col min-h-screen debug">
  <!-- Header -->
  <header class="debug">
    
    <slot name="header" {title}>
     {title}
    </slot>
  </header>

  <!-- Main: Sidebar + Content -->
  <div class="flex flex-1 overflow-hidden">
    <!-- Sidebar -->
    <aside class="debug">
      <slot name="sidebar">
       
      </slot>
    </aside>

    <!-- Content -->
    <main class="debug flex-1">
      <slot>
       
        <p class="text-gray-500">Default content — nothing provided.</p>
      </slot>
    </main>
  </div>

 
</div>
