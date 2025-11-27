import { check, sleep } from 'k6';
import http from 'k6/http';

// ----------------------------------------------------
// Dữ liệu tĩnh (Static Data)
// ----------------------------------------------------

// URL API cần kiểm tra
const API_URL = 'http://localhost:8080/api/users/update';

// JWT dùng trong Header Authorization
const AUTH_HEADER_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjQwODc4ODQsImlhdCI6MTc2NDA4MDY4NCwiRGF0YSI6eyJyb2xlSWQiOm51bGwsImVtYWlsIjoiIiwidGVuYW50IjoiIiwidXNlcklkIjoiNGEzOGJjMDktYzFlYy00OTFkLWFjYWYtNzliYzY5NzI4ZDNlIiwidXNlcm5hbWUiOiJhZG1pbiJ9fQ.w3BfX-4RYs0nAqY7jwyJeJO1m7xshQcIKrdWGcg9gLw';

// JWT dùng trong JSON Body
const BODY_TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7IklkIjoiN2MzNzgwOTctYWU0ZS00MDEyLWFlNGItMmZiMGU2MDQ2ZmY4IiwiVXNlcm5hbWUiOiJVLTAwOTIiLCJDcmVhdGVkQnkiOiJhZG1pbiIsIkNyZWF0ZWRPbiI6IjIwMjUtMTEtMjVUMTE6MjU6MjAuMTYzMTY5KzA3OjAwIiwiTW9kaWZpZWRCeSI6bnVsbCwiTW9kaWZpZWRPbiI6bnVsbH19.Mskj2ymVHF_JohyaP9bW5nAj1JRSFN8O5ZqbMX7yAAE';

// ----------------------------------------------------
// Cấu hình k6 (k6 Options)
// ----------------------------------------------------
export const options = {
    // 20 người dùng ảo chạy trong 1 phút
    vus: 200,
    duration: '1m',
    thresholds: {
        // Đảm bảo 95% yêu cầu hoàn thành dưới 300ms
        http_req_duration: ['p(95)<300'],
        // Tỷ lệ giao dịch thành công (status 200) phải trên 99%
        checks: ['rate>0.99'],
    },
};

// ----------------------------------------------------
// Hàm chính của VU (Virtual User)
// ----------------------------------------------------
export default function () {
    // 1. Chuẩn bị Header
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${AUTH_HEADER_TOKEN}`, // Gán JWT vào header
    };

    // 2. Chuẩn bị Body (Payload)
    // Cập nhật 'modifiedOn' thành thời gian hiện tại
    const currentISOTime = new Date().toISOString();

    const payload = JSON.stringify({
        "data": {
            "id": "4adf81c7-8a90-4fd8-baa1-65065a9d600d",
            "createdOn": "2025-11-25T11:25:20.163169+07:00",
            "createdBy": "admin",
            "modifiedOn": currentISOTime, // Sử dụng thời gian hiện tại
            "ModifiedBy": "admin",
            "description": "Updated via k6 load test",
            "username": "U-0092",
            "displayName": "Lương Văn Trọng",
            "email": "tronglv@qtsc.com.vn",
            "roleId": null,
            "isSysAdmin": false
        },
        "token": BODY_TOKEN // Gán JWT vào body
    });

    // 3. Thực hiện yêu cầu POST
    const res = http.post(API_URL, payload, { headers: headers });

    // 4. Kiểm tra (Checks)
    check(res, {
        'Status is 200': (r) => r.status === 200,
        'Body is not empty': (r) => r.body.length > 0,
        // Có thể thêm kiểm tra nội dung phản hồi nếu cần, ví dụ:
        // 'Response success field is true': (r) => r.json().success === true,
    });

    // Tạm dừng 1 giây giữa các lần lặp của mỗi VU
    sleep(1);
}