//http://localhost:8080/api/media/upload
//play load field la file (binary)
// read file from D:\code\go\wx-example\wx-example\media\uploads\test.txt 

import { check } from 'k6';
import http from 'k6/http';

export const options = {
    vus: 200,          // 200 virtual users
    duration: '30s',   // chạy 30 giây
};

const fileData = open('D:\\code\\go\\wx-example\\wx-example\\media\\test.txt', 'b'); // 'b' = binary

export default function () {
    const url = 'http://localhost:8081/api/media/upload';

    const payload = {
        file: http.file(fileData, 'test.txt', 'text/plain'), // tạo multipart file
    };

    const res = http.post(url, payload);
    //console.log(res.status)
    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
// k6 run upload-file-200-vu-chi.js