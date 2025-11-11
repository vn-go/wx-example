<script lang="ts">
    import Layout from '@layouts/main.svelte'
    let username="abc"
</script> 
<!-- Gán nội dung cho các slot bằng attribute slot="..." -->
<Layout title="Dashboard">
    <!-- header slot: nhận slot-prop `title` từ layout qua let:title -->
    <div slot="header" let:title class="main-header">
      <h2>{title} — Xin chào, {username}!</h2>
      <p>Trang quản trị</p>
    </div>
  
    <!-- sidebar slot -->
    <div slot="sidebar">
      <ul>
        <li><a href="/dashboard">Dashboard</a></li>
        <li><a href="/users">Users</a></li>
        <li><a href="/settings">Settings</a></li>
      </ul>
    </div>
  
    <!-- content (default slot) -->
    <div>
      <h3>Tổng quan</h3>
      <p>Đây là nội dung chính của trang.</p>
    </div>
  
    <!-- footer slot -->
    <div slot="footer">
      <small>Liên hệ: admin@myapp.test</small>
    </div>
  </Layout>