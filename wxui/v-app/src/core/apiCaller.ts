const API_BASE = import.meta.env.VITE_API_BASE_URL

// token lưu ở localStorage ví dụ
function getToken() {
    return sessionStorage.getItem('access_token') || ''
}

function setToken(token: string) {
    sessionStorage.setItem('access_token', token)
}

function removeToken() {
    sessionStorage.removeItem('access_token')
}
export const apiCaller = {
    getToken: () => {
        return sessionStorage.getItem('access_token') || ''
    },
    get: async (path: string) => {
        const res = await fetch(`${API_BASE}${path}`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${getToken()}`,
                'Content-Type': 'application/json'
            }
        })
        return res.json()
    },

    post: async (path: string, data: any) => {
        const res = await fetch(`${API_BASE}${path}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${getToken()}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        return res.json()
    },

    login: async (username: string, password: string) => {
        const formData = new URLSearchParams({
            grant_type: 'password',
            username,
            password
        });
        //http://localhost:8080/api/auth/login
        const res = await fetch(`${API_BASE}/api/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            credentials: "include",  // 🔥 IMPOTANT!
            body: formData.toString(),
        })

        const data = await res.json()

        return data
    }
}
