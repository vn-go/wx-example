import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.1/index.js';
import { htmlReport } from 'https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js';
import { check, sleep } from 'k6';
import http from 'k6/http';

// Cấu hình test
export const options = {
    vus: 100,              // 100 người dùng đồng thời
    duration: '30s',       // chạy 30 giây
    // Hoặc muốn load nặng hơn:
    // stages: [
    //   { duration: '2m', target: 200 },  // ramp-up
    //   { duration: '5m', target: 500 },  // sustain
    //   { duration: '2m', target: 0 },    // ramp-down
    // ],

    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% request dưới 500ms
        http_req_failed: ['rate<0.01'],   // lỗi < 1%
    },
};

const BASE_URL = 'http://localhost:5106';
const LOGIN_URL = `${BASE_URL}/api/auth/login`;

// Dữ liệu đăng nhập
const payload = 'username=admin&password=%2F%5Cdmin123451212';
// password = '/\dmin123451212' → đã URL-encoded đúng
// Nếu muốn k6 tự encode thì dùng object (xem cách 2 bên dưới)

export default function () {
    // Cách 1: Dùng params + form-urlencoded string (nhanh, gọn)
    const params = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
            // Nếu cần cookie, auth thì thêm ở đây
            // 'Cookie': 'refresh_token=xxx',
        },
    };

    const res = http.post(LOGIN_URL, payload, params);

    // Kiểm tra kết quả
    const success = check(res, {
        'status là 200': (r) => r.status === 200,
        'có access_token trong body': (r) => {
            try {
                const json = r.json();
                return json.access_token && json.access_token.length > 50;
            } catch (err) {
                return false;
            }
        },
        'không bị 415': (r) => r.status !== 415,
        'không bị 401/403': (r) => r.status !== 401 && r.status !== 403,
    });

    if (!success) {
        console.log(`Failed request: ${res.status} ${res.body}`);
    }

    sleep(1); // nghỉ 1 giây giữa các request (giảm tải)
}


// Tạo báo cáo HTML đẹp khi chạy xong
export function handleSummary(data) {
    return {
        'k6-report.html': htmlReport(data),
        stdout: textSummary(data, { indent: '→', enableColors: true }),
    };
}

// k6 run login-test-net.js