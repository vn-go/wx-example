import { textSummary } from 'https://jslib.k6.io/k6-summary/0.0.4/index.js';
import { check } from 'k6';
import http from 'k6/http';

export function handleSummary(data) {
    return {
        'summary.json': JSON.stringify(data), // lưu gọn summary ra file JSON
        stdout: textSummary(data, { indent: ' ', enableColors: true })
    };
}

export const options = {
    thresholds: {},
    summaryTrendStats: ["avg", "p(95)", "p(99)"], // bớt percentiles
    stages: [
        { duration: '30s', target: 50 },    // tăng lên 50 VUs trong 30s
        { duration: '30s', target: 200 },   // tăng lên 200 VUs trong 30s
        { duration: '30s', target: 500 },   // tăng lên 500 VUs trong 30s
        { duration: '30s', target: 1000 },  // tăng lên 1000 VUs trong 30s
        { duration: '30s', target: 0 },     // hạ dần về 0
    ],
};

export default function () {
    let res = http.get('http://localhost:3000/hello');

    // Kiểm tra response code = 200
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    // Nghỉ 1s giữa các request (tùy chỉnh)
    //sleep(1);
}
// k6 run hello-chi.js
//k6 run --out influxdb=http://localhost:8086/k6 --out-options pushInterval=3s hello-chi.js
//$env:K6_INFLUXDB_PUSH_INTERVAL="3s"
//$env:K6_INFLUXDB_CONCURRENT_WRITES="4"
//k6 run --out influxdb=http://localhost:8086/k6 hello-chi.js
/*
$env:K6_INFLUXDB_PUSH_INTERVAL="1s"   # gửi dữ liệu nhỏ hơn, thường xuyên hơn
$env:K6_INFLUXDB_CONCURRENT_WRITES="10"

.\influxd.exe --config .\influxdb.conf

*/