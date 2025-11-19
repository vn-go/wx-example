const API_BASE = (import.meta as any).env.VITE_API_BASE_URL

const login = async (username: string, password: string) => {
    const formData = new URLSearchParams({
        grant_type: 'password',
        username,
        password
    });
    const res = await fetch(`${API_BASE}/api/auth/login`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        credentials: "include",  // 🔥 IMPOTANT!
        body: formData.toString(),
    });

    const data = await res.json()

    return data
}
export default login