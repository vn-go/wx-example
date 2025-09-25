import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 100,       // số lượng virtual users
    iterations: 100 // chỉ chạy 1 lần
};

export default function () {
    const url = 'http://localhost:8080/api/tenant/create';
    const payload = JSON.stringify({
        name: "" // gửi rỗng
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 401': (r) => r.status === 401,
    });
}
//k6 run tenant_create_test.js