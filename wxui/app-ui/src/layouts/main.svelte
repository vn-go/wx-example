<script lang="ts">
  import '../app.css';
  export let title = "My App";
</script>

<div class="flex flex-col min-h-screen debug">
  <!-- Header -->
  <slot name="header" {title}>
    {title}
   </slot>

  <!-- Main: Sidebar + Content -->
  <div class="flex flex-1 overflow-hidden">
    <!-- Sidebar -->
    <aside class="debug" style="position: relative; width: 200px">
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
