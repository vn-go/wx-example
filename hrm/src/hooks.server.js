import { redirect } from '@sveltejs/kit';

export async function handle({ event, resolve }) {
    const token = event.cookies.get('session'); // lấy token từ cookie

    // Nếu đang không ở /login mà lại chưa có token -> redirect
    if (!token && !event.url.pathname.startsWith('/login')) {
        throw redirect(303, '/login');
    }

    // Nếu đã login hoặc đang ở /login -> cho qua
    return resolve(event);
}
