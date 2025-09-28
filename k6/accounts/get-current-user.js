//http://localhost:8080/api/accounts/test

// import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.4/index.js';
import { check } from 'k6'; // 👈 thêm dòng này
import http from 'k6/http';
// export function handleSummary(data) {
//     return {
//         'summary.json': JSON.stringify(data), // lưu gọn summary ra file JSON
//         stdout: textSummary(data, { indent: ' ', enableColors: true })
//     };
// }

export const options = {
    //executor: 'constant-vus',
    vus: 200,
    duration: '30s',

};
let token = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhOTQ3YWQzYy1lNzQ3LTQ3YzctODc4ZS02YWIyMjlmNjQyNjIiLCJleHAiOjE3NTg5NjkxNzUsImlhdCI6MTc1ODk1NDc3NSwicm9sZSI6IiIsImVtYWlsIjoiIn0.Z7oFHA06wtiRVVuEOjRNDmM2FKvLIbiub6L2UeV7z_s`
export default function () {

    let url = 'http://localhost:8080/api/accounts/me';

    let headers = {
        'Authorization': token,
        'Content-Type': 'application/json',
    };

    let res = http.post(url, null, { headers: headers });
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
// k6 run get-current-user.js
