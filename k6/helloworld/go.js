// go.js
import { check } from 'k6';
import http from 'k6/http';



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
    let res = http.get('http://localhost:8081/api/app/hello');

    // Kiểm tra response code = 200
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    
}
