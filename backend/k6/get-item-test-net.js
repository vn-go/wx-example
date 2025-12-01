import { check } from "k6";
import http from "k6/http";

const TOKEN = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3NjQ1MDUzNjksImV4cCI6MTc2NDU5MTc2OSwicm9sZV9pZCI6IiIsImVtYWlsIjoiIiwidGVuYW50IjoibXktYXBwIiwidXNlcl9pZCI6IjRhMzhiYzA5LWMxZWMtNDkxZC1hY2FmLTc5YmM2OTcyOGQzZSIsInVzZXJuYW1lIjoiYWRtaW4iLCJpc19zeXNfYWRtaW4iOiJ0cnVlIiwibmJmIjoxNzY0NTA1MzY5LCJpc3MiOiJteS1hcHAifQ.z5wdPYC-4cYMrYdBXj2-XsnOwQW5ergSaEF_dZqtZL8`;

// Payload phải là JSON string nếu endpoint dùng [FromBody] string
const guidPayload = JSON.stringify("4a38bc09-c1ec-491d-acaf-79bc69728d3e");



// k6-test-users-list.js
import { sleep } from 'k6';
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

const params = {
    headers: {
        'authorization': TOKEN,
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    },
};

export default function () {
    let res = http.post(`${BASE_URL}/api/system/Users/get-item`, guidPayload, params);
    if (res.status === 200 && res.body) {
        let data = res.json(); // chỉ parse khi có body
    }
    const success = check(res, {
        'status is 200': (r) => r.status === 200,

        'returns datacontract': (r) => {
            const body = r.json();
            return body;
        },
    });

    errorRate.add(!success);

    // In ra log mỗi 100 VU để theo dõi realtime
    // if (__VU % 100 === 0) {
    //     console.log(`VU: ${__VU} | ITER: ${__ITER} | Status: ${res.status} | Time: ${res.timings.duration}ms | RPS hiện tại ≈ ${1000 / res.timings.duration * 500}`);
    // }

    sleep(0.1); // giảm sleep nếu muốn bơm tải mạnh hơn
}
// k6 run get-all-users-net.js