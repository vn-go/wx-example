export async function POST({ request }) {
    const { username, password } = await request.json();

    if (username === 'admin' && password === '123456') {
        return new Response(JSON.stringify({ token: 'abc123' }), {
            status: 200,
            headers: { 'Content-Type': 'application/json' }
        });
    }

    return new Response(JSON.stringify({ error: 'Invalid credentials' }), {
        status: 401,
        headers: { 'Content-Type': 'application/json' }
    });
}
