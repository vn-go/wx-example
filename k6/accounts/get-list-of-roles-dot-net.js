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
  duration: '60s',

};
let token = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhOTQ3YWQzYy1lNzQ3LTQ3YzctODc4ZS02YWIyMjlmNjQyNjIiLCJleHAiOjE3NTk2MDI4MDcsImlhdCI6MTc1OTU4ODQwNywicm9sZSI6IiIsImVtYWlsIjoiIn0.NZtq2cYjlCr4ugVSE8UKCnF1FnxtamTdjbAcZcunRP0`
export default function () {

  // let url = 'http://localhost:8080/api/accounts/get-list-of-roles';
  //let url = 'http://localhost:8080/api/accounts/get-list-of-roles-s-q-l'
  //let url = 'http://localhost:8080/api/pure/get'
  //let url = 'http://localhost:8080/hello'
  //let url = 'http://localhost:8080/api/hz/check'
  let url = "http://localhost:5000/api/roles"
  //let url = "http://localhost:8080/api/accounts/get-list-of-roles"
  let headers = {
    'Authorization': token,
    'Content-Type': 'application/json',
  };

  let res = http.get(url);
  //let res = http.post(url)
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

}

// k6 run get-list-of-roles-dot-net.js
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
/*
PS D:\code\go\wx-example\wx-example\k6\accounts> k6 run get-list-of-roles-dot-net.js

         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: get-list-of-roles-dot-net.js
        output: -

     scenarios: (100.00%) 1 scenario, 200 max VUs, 1m0s max duration (incl. graceful stop) 
              * default: 200 looping VUs for 30s (gracefulStop: 30s)



  █ TOTAL RESULTS

    checks_total.......................: 27202   904.041997/s
    checks_succeeded...................: 100.00% 27202 out of 27202
    checks_failed......................: 0.00%   0 out of 27202

    ✓ status is 200

    HTTP
    http_req_duration.......................................................: avg=220.14ms min=11.86ms med=217.56ms max=583.23ms p(90)=272.88ms p(95)=301.98ms
      { expected_response:true }............................................: avg=220.14ms min=11.86ms med=217.56ms max=583.23ms p(90)=272.88ms p(95)=301.98ms
    http_req_failed.........................................................: 0.00%  0 out of 27202
    http_reqs...............................................................: 27202  904.041997/s

    EXECUTION
    iteration_duration......................................................: avg=220.96ms min=11.86ms med=218.04ms max=596.23ms p(90)=274.21ms p(95)=304.75ms
    iterations..............................................................: 27202  904.041997/s
    vus.....................................................................: 200    min=200        max=200
    vus_max.................................................................: 200    min=200        max=200

    NETWORK
    data_received...........................................................: 2.1 GB 69 MB/s
    data_sent...............................................................: 2.5 MB 83 kB/s




running (0m30.1s), 000/200 VUs, 27202 complete and 0 interrupted iterations
default ✓ [======================================] 200 VUs  30s
-----
PS D:\code\go\wx-example\wx-example\k6\accounts> k6 run get-list-of-roles.js

         /\      Grafana   /‾‾/  
    /\  /  \     |\  __   /  /   
   /  \/    \    | |/ /  /   ‾‾\ 
  /          \   |   (  |  (‾)  |
 / __________ \  |_|\_\  \_____/ 

     execution: local
        script: get-list-of-roles.js
  █ TOTAL RESULTS

    checks_total.......................: 96466   3209.644652/s
    checks_succeeded...................: 100.00% 96466 out of 96466
    checks_failed......................: 0.00%   0 out of 96466

    ✓ status is 200

    HTTP
    http_req_duration.......................................................: avg=62.09ms min=7.01ms med=58.69ms max=244.21ms p(90)=92.49ms p(95)=104.36ms
      { expected_response:true }............................................: avg=62.09ms min=7.01ms med=58.69ms max=244.21ms p(90)=92.49ms p(95)=104.36ms
    http_req_failed.........................................................: 0.00%  0 out of 96466
    http_reqs...............................................................: 96466  3209.644652/s

    EXECUTION
    iteration_duration......................................................: avg=62.22ms min=7.01ms med=58.84ms max=244.19ms p(90)=92.64ms p(95)=104.5ms
    iterations..............................................................: 96466  3209.644652/s
    vus.....................................................................: 200    min=200        max=200
    vus_max.................................................................: 200    min=200        max=200

    NETWORK
    data_received...........................................................: 7.4 GB 247 MB/s
    data_sent...............................................................: 42 MB  1.4 MB/s




running (0m30.1s), 000/200 VUs, 96466 complete and 0 interrupted iterations
default ✓ [======================================] 200 VUs  30s
*/
/*

dotnet tool install -g dotnet-trace
Get-Process dotnet, hrm | Select-Object -First 1 Id
----
$env:PATH += ";C:\Users\MSI CYBORG\.dotnet\tools"
k6 run get-list-of-roles-dot-net.js
dotnet-trace collect -p 25940 --duration 00:00:30 --profiler speedscope --output ./hrm_cpu_profile.json
dotnet-trace collect -p 25940 --duration 00:00:30 -c cpu-sampling -o ./hrm_cpu_profile.nettrace

dotnet-trace report ./hrm_cpu_profile.nettrace topN -n 10



*/