<script>
	let username = '';
	let password = '';
	let error = '';

	function handleLogin() {
		if (username === 'admin' && password === '123456') {
			alert('Đăng nhập thành công ✅');
			// ở thực tế thì redirect về dashboard
			window.location.href = '/dashboard';
		} else {
			error = 'Sai email hoặc mật khẩu!';
		}
	}
</script>

<div class="flex items-center justify-center min-h-screen bg-gray-100">
	<div class="w-full max-w-md bg-white rounded-2xl shadow-lg p-8">
		<h1 class="text-2xl font-bold text-center text-gray-800 mb-6">Đăng nhập</h1>

		{#if error}
			<p class="text-red-600 text-sm mb-4">{error}</p>
		{/if}

		<form on:submit|preventDefault={handleLogin} class="space-y-4">
			<div>
				<label class="block text-sm font-medium text-gray-700">Username</label>
				<input
					type="text"
					bind:value={username}
					placeholder="Your username"
					class="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
					required
				/>
			</div>

			<div>
				<label class="block text-sm font-medium text-gray-700">Mật khẩu</label>
				<input
					type="password"
					bind:value={password}
					placeholder="••••••••"
					class="mt-1 block w-full rounded-lg border border-gray-300 px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
					required
				/>
			</div>

			<button
				type="submit"
				class="w-full bg-blue-600 text-white font-semibold py-2 rounded-lg hover:bg-blue-700 transition"
			>
				Đăng nhập
			</button>
		</form>

		<p class="mt-6 text-center text-sm text-gray-600">
			Chưa có tài khoản?
			<a href="/register" class="text-blue-600 hover:underline">Đăng ký</a>
		</p>
	</div>
</div>
