// src/lib/utils/login.js
import { setAccessToken } from '@store/auth';

// Lấy URL cơ sở từ biến môi trường
const API_BASE_URL = import.meta.env.VITE_PUBLIC_API_BASE_URL || 'http://localhost:5000/api';

/**
 * Hàm đăng nhập, gửi POST với form-urlencoded.
 * @param {string} username 
 * @param {string} password 
 * @returns {Promise<boolean>} true nếu thành công, false nếu thất bại
 */
export async function login(username, password) {
    const formData = new URLSearchParams({
        grant_type: 'password',
        username,
        password
    });

    const LOGIN_ENDPOINT = `${API_BASE_URL}/auth/login`;

    try {
        const response = await fetch(LOGIN_ENDPOINT, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: formData.toString(),
        });

        // Kiểm tra phản hồi HTTP
        if (!response.ok) {
            const errorBody = await response.json().catch(() => ({}));
            const message =
                errorBody.error_description ||
                errorBody.message ||
                `Đăng nhập thất bại (${response.status})`;
            throw new Error(message);
        }

        // Parse JSON response
        const res = await response.json();

        // Kiểm tra có token không
        if (res.access_token) {
            setAccessToken(res.access_token);
            console.log('✅ Token saved to store:', res.access_token);
        } else {
            console.warn('⚠️ Không tìm thấy access_token trong phản hồi:', res);
            throw new Error('Phản hồi không chứa access_token');
        }

        return true;
    } catch (error) {
        console.error('Lỗi khi thực hiện đăng nhập:', error);
        throw error;
    }
}
