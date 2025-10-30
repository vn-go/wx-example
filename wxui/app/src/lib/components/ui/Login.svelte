<script lang="ts">
	import Button from '@components/ui/Button.svelte';
	import { closeDialog } from '@store/dialogStore';
	import { login } from '@utils/login';

	let username = '';
	let password = '';
	let loginOk: boolean | null = null; // 👈 trạng thái ban đầu

	async function handleLogin() {
		try {
			const ok = await login(username, password);
			loginOk = ok;
			if (ok) {
				closeDialog();
			}
		} catch (err) {
			console.error('Đăng nhập thất bại:', err);
			loginOk = false;
		}
	}
</script>

<div class="bg-white p-6 rounded-xl">
	<h1 class="text-xl font-semibold mb-4">Đăng nhập</h1>

	<input
		type="text"
		placeholder="Tên đăng nhập"
		class="border rounded w-full p-2 mb-3"
		bind:value={username}
	/>

	<input
		type="password"
		placeholder="Mật khẩu"
		class="border rounded w-full p-2 mb-4"
		bind:value={password}
	/>

	<Button className="w-full" on:click={handleLogin}>Đăng nhập</Button>

	{#if loginOk === false}
		<p class="text-red-500 text-sm mt-3">Đăng nhập thất bại, vui lòng thử lại.</p>
	{/if}
</div>
