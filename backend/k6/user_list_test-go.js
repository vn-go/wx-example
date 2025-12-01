import { check, sleep } from 'k6';
import http from 'k6/http';

// 1. Cấu hình kịch bản tải (Load Options)
export const options = {
    // 10 người dùng ảo (VUs)
    vus: 200,
    // Chạy trong 30 giây
    duration: '60s',
    // Định ngưỡng kiểm tra: 95% request phải hoàn thành dưới 500ms
    thresholds: {
        http_req_duration: ['p(95)<500'],
    },
};

// 2. Dữ liệu POST (Nếu endpoint yêu cầu một body)
const payload = JSON.stringify({}); // Sử dụng body rỗng
// Hoặc sử dụng body phức tạp hơn nếu API yêu cầu:
/*
const payload = JSON.stringify({
    page: 1,
    pageSize: 10
});
*/

// 3. Header yêu cầu (Bắt buộc phải là JSON)
const params = {
    headers: {
        'Content-Type': 'application/json',
        // Thêm Authorization Header nếu cần (ví dụ: token JWT)
        // 'Authorization': 'Bearer YOUR_JWT_TOKEN', 
    },
};

// 4. Hàm Main (Hàm được chạy bởi mỗi VU)
export default function () {
    const url = 'http://localhost:8080/api/system/user-portal/get-list';

    // Thực hiện yêu cầu POST
    const res = http.post(url, payload, params);

    // Kiểm tra (Checks)
    check(res, {
        'status is 200': (r) => r.status === 200,
        'response body is not empty': (r) => r.body.length > 0,
    });

    // Đợi 0.5 giây trước khi thực hiện request tiếp theo
    sleep(0.5);
}