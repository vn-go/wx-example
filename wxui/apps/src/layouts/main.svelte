<script>
  export let title = "My App";
</script>

<div class="flex flex-col min-h-screen">
  <!-- Header -->
  <header class="bg-gray-100 p-4 shadow">
    <h1 class="tw-red-title">OK</h1>
    <slot name="header" {title}>
      <h1 class="text-2xl font-bold text-gray-800">{title}</h1>
    </slot>
  </header>

  <!-- Main: Sidebar + Content -->
  <div class="flex flex-1 overflow-hidden">
    <!-- Sidebar -->
    <aside class="w-56 bg-gray-50 border-r border-gray-200 p-4 overflow-y-auto hidden md:block">
      <slot name="sidebar">
        <nav>
          <ul class="space-y-2">
            <li><a href="/" class="text-gray-700 hover:text-gray-900">Home</a></li>
            <li><a href="/dashboard" class="text-gray-700 hover:text-gray-900">Dashboard</a></li>
          </ul>
        </nav>
      </slot>
    </aside>

    <!-- Content -->
    <main class="flex-1 p-6 overflow-y-auto">
      <slot>
       
        <p class="text-gray-500">Default content — nothing provided.</p>
      </slot>
    </main>
  </div>

  <!-- Footer -->
  <footer class="bg-gray-100 p-4 text-center text-sm text-gray-500 border-t border-gray-200">
    <slot name="footer">
      © 2025 My App
    </slot>
  </footer>
</div>
