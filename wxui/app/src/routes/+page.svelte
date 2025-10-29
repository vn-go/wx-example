<script>
	import Button from '@components/ui/Button.svelte';
</script>

<h1 class="text-3xl font-bold mb-4">Trang chủ</h1>
<p>Chào mừng bạn đến với hệ thống quản trị!</p>

<div class="flex items-center justify-center h-screen bg-gray-100">
	<div class="bg-white p-6 rounded-xl shadow w-96">
		<h1 class="text-xl font-semibold mb-4">Đăng nhập</h1>
		<input type="text" placeholder="Tên đăng nhập" class="border rounded w-full p-2 mb-3" />
		<input type="password" placeholder="Mật khẩu" class="border rounded w-full p-2 mb-4" />
		<Button className="w-full">Đăng nhập</Button>
	</div>
</div>
