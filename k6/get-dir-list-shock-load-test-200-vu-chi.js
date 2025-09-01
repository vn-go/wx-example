import { check } from 'k6';
import http from 'k6/http';



export const options = {

    executor: 'constant-vus',
    vus: 200,
    duration: '30s',
};

export default function () {
    let res = http.get('http://localhost:8081/api/media/list');

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
// k6 run get-dir-list-shock-load-test-200-vu-chi.js