// k6-test-users-list.js
import { check, sleep } from 'k6';
import http from 'k6/http';
import { Rate } from 'k6/metrics';

// Đo tỷ lệ lỗi
export let errorRate = new Rate('errors');

export let options = {
    stages: [
        { duration: '30s', target: 100 },   // warm-up
        { duration: '2m', target: 500 },   // tải chính (500 VU ≈ 3000–5000 RPS tùy server)
        { duration: '30s', target: 1000 },  // peak
        { duration: '1m', target: 0 },     // ramp-down
    ],
    thresholds: {
        http_req_duration: ['p(95)<300', 'p(99)<600'], // 95% request < 300ms, 99% < 600ms
        errors: ['rate<0.01'], // lỗi < 1%
    },
};

const BASE_URL = 'http://localhost:5106'; // ĐỔI THÀNH 5107 hoặc domain khi test Go
const TOKEN = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjQ0ODg3OTIsImV4cCI6MTc2NDU3NTE5Miwicm9sZV9pZCI6IiIsImVtYWlsIjoiIiwidGVuYW50IjoibXktYXBwIiwidXNlcl9pZCI6IjRhMzhiYzA5LWMxZWMtNDkxZC1hY2FmLTc5YmM2OTcyOGQzZSIsInVzZXJuYW1lIjoiYWRtaW4iLCJpc19zeXNfYWRtaW4iOiJ0cnVlIiwibmJmIjoxNzY0NDg4NzkyLCJpc3MiOiJteS1hcHAifQ.53Sw1Qrvh0W-5B9Ivh5s_BUpYEbmZ2Nhu-LUHTGkWl0';

const params = {
    headers: {
        'authorization': TOKEN,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    },
};

export default function () {
    let res = http.post(`${BASE_URL}/api/system/Users/list`, {}, params);
    if (res.status === 200 && res.body) {
        let data = res.json(); // chỉ parse khi có body
    } else {
        console.log(`Status: ${res.status}, body: "${res.body}"`);
    }
    const success = check(res, {
        'status is 200': (r) => r.status === 200,
        'response has data array': (r) => {
            const body = r.json();
            return body && body && Array.isArray(body);
        },
        'returns ~300 users': (r) => {
            const body = r.json();
            return body && body.length >= 200; // ít nhất 200 để chắc có dữ liệu
        },
    });

    errorRate.add(!success);

    // In ra log mỗi 100 VU để theo dõi realtime
    if (__VU % 100 === 0) {
        console.log(`VU: ${__VU} | ITER: ${__ITER} | Status: ${res.status} | Time: ${res.timings.duration}ms | RPS hiện tại ≈ ${1000 / res.timings.duration * 500}`);
    }

    sleep(0.1); // giảm sleep nếu muốn bơm tải mạnh hơn
}
// k6 run get-all-users-net.js