
import { check, sleep } from 'k6';
import http from 'k6/http';

// Cấu hình test
export const options = {
    vus: 100,           // 50 người dùng ảo đồng thời
    duration: '30s',   // Chạy trong 30 giây
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% request < 200ms
        http_req_failed: ['rate<0.01'],   // Lỗi < 1%
    },
};

// URL API
const BASE_URL = 'http://localhost:8080/api/accounts/get-list-of-accounts'; // THAY BẰNG DOMAIN THẬT
const VIEW_PATH = '/system/user';

// JWT Token
const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3NmVjYTUxYy05MjdjLTRkMTMtYjRmNi1iYTVlNTc0NWMxNTAiLCJleHAiOjE3NjI0NTQxNjYsImlhdCI6MTc2MjQzOTc2Niwicm9sZSI6IiIsImVtYWlsIjoiIn0.kwYi3b0JzjRiCeyxeCUA5BwVPZKSitV3g7TT4BoeQ1w';

// Payload
const payload = JSON.stringify({
    index: 0,
    size: 20,
    orderBy: ["username"]

});

const params = {
    headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${TOKEN}`,
        'View-Path': VIEW_PATH,
    },
};

export default function () {
    const url = `${BASE_URL}`;

    const res = http.post(url, payload, params);

    // Kiểm tra response
    check(res, {
        'status is 200': (r) => r.status === 200,
        // 'response has data': (r) => {
        //     const body = JSON.parse(r.body);
        //     return Array.isArray(body.data) && body.data.length > 0;
        // },
        // 'total matches expected': (r) => {
        //     const body = JSON.parse(r.body);
        //     return body.total !== undefined;
        // },
    });

    // Log lỗi nếu có
    if (res.status !== 200) {
        console.error(`Request failed: ${res.status} ${res.body}`);
    }

    sleep(0.1); // Giảm tải, tránh spam quá nhanh
}

// Tạo HTML report khi chạy xong
// export function handleSummary(data) {
//     return {
//         'report.html': htmlReport(data),
//         stdout: textSummary(data, { indent: ' ', enableColors: true }),
//     };
// }