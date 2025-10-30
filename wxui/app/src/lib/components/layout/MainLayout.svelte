<!--wxui\app\src\lib\components\layout\MainLayout.svelte-->
<script>
	$: console.log('📊 Layout render:', $dialogVisible, $dialogComponent);

	import DialogWrapper from '@components/ui/Dialog.svelte';
	import {
		closeDialog,
		dialogComponent,
		dialogProps,
		dialogVisible,
		showDialog
	} from '@store/dialogStore';

	import { accessToken } from '@store/auth';
	//<--loi co nay
	// import { showDialog } from '@store/dialogStore';
	import { onMount } from 'svelte';
	import { get } from 'svelte/store';
	import Footer from './Footer.svelte';
	import Header from './Header.svelte';
	import Sidebar from './Sidebar.svelte';
	// store popup login mà bạn có
	onMount(() => {
		console.log('✅ App just mounted or browser reloaded');
		const token = get(accessToken); // 👈 đọc giá trị hiện tại
		if (!token) {
			console.log('⚠️ Chưa có access token → hiển thị form đăng nhập');
			showDialog('Login');
		} else {
			console.log('✅ Đã có access token');
			// có thể gọi API xác thực hoặc load user info ở đây
		}
	});
</script>

{#if $dialogVisible}
	<DialogWrapper component={$dialogComponent} props={$dialogProps} on:close={closeDialog} />
{/if}

<div class="flex h-screen bg-gray-50 text-gray-900">
	<Sidebar />
	<div class="flex flex-col flex-1">
		<Header />
		<main class="flex-1 p-6 debug">
			<button on:click={() => showDialog('Login')}>Show Login Dialog</button>
			<slot />
		</main>
		<Footer />
	</div>
</div>
