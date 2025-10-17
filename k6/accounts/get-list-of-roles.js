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
let token = `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhOTQ3YWQzYy1lNzQ3LTQ3YzctODc4ZS02YWIyMjlmNjQyNjIiLCJleHAiOjE3NTk2NDU2MzUsImlhdCI6MTc1OTYzMTIzNSwicm9sZSI6IiIsImVtYWlsIjoiIn0.aFm-8oQ8Gdc7eA0xpMcr9sr3Z5s010yHLqrMFwsLl08`
export default function () {

  // let url = 'http://localhost:8080/api/accounts/get-list-of-roles';
  //let url = 'http://localhost:8080/api/accounts/get-list-of-roles-s-q-l'
  //let url = 'http://localhost:8080/api/pure/get'
  //let url = 'http://localhost:8080/hello'
  //let url = 'http://localhost:8080/api/hz/check'
  let url = "http://localhost:8080/api/roles/GetListOfRoles"
  //let url = "http://localhost:8080/api/accounts/get-list-of-roles"
  let headers = {
    'Authorization': token,
    'Content-Type': 'application/json',
  };

  let res = http.post(url, JSON.stringify({
    index: 0,
    size: 1000,
    orderBy: null,
  }), { headers: headers });
  //let res = http.post(url)
  check(res, {
    'status is 200': (r) => r.status === 200,
  });

}

// k6 run get-list-of-roles.js
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
/*
test load list of roles in database 
381 rows of 
{
    "id": 1121,
    "roleId": "27f590ee-de6a-45a3-930f-3fb602954dae",
    "code": "X-0300",
    "name": "Role-0300",
    "description": "",
    "createdOn": "2025-09-27T11:12:53Z",
    "modifiedOm": null,
    "createdBy": "admin",
    "isActive": true
}
-- version 2------
Build ID: D:\code\go\wx-example\wx-example\wxapi\wxapi.exe2025-10-04 21:51:25.7431819 +0700 +07
Type: cpu
Time: 2025-10-04 22:02:01 +07
Duration: 30.16s, Total samples = 323.66s (1073.22%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 209.36s, 64.69% of 323.66s total
Dropped 1036 nodes (cum <= 1.62s)
Showing top 10 nodes out of 228
      flat  flat%   sum%        cum   cum%
   111.28s 34.38% 34.38%    111.28s 34.38%  runtime.stdcall2
    38.78s 11.98% 46.36%     38.88s 12.01%  runtime.stdcall0
    22.70s  7.01% 53.38%     24.18s  7.47%  runtime.cgocall
     8.90s  2.75% 56.13%     12.86s  3.97%  runtime.findObject
     6.65s  2.05% 58.18%      6.65s  2.05%  runtime.stdcall1
     5.68s  1.75% 59.94%      5.68s  1.75%  runtime.(*mspan).base (inline)
     4.48s  1.38% 61.32%      9.48s  2.93%  github.com/json-iterator/go.(*Stream).WriteString
     3.79s  1.17% 62.49%      3.79s  1.17%  runtime.nextFreeFast (inline)
     3.66s  1.13% 63.62%      3.66s  1.13%  runtime.procyield
     3.44s  1.06% 64.69%     16.30s  5.04%  runtime.scanobject
(pprof)
--- version 1-------
PS D:\code\go\wx-example\wx-example\k6\accounts> go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in C:\Users\MSI CYBORG\pprof\pprof.wxapi.exe.samples.cpu.002.pb.gz
File: wxapi.exe
Build ID: D:\code\go\wx-example\wx-example\wxapi\wxapi.exe2025-10-04 21:51:25.7431819 +0700 +07
Type: cpu
Time: 2025-10-04 22:05:38 +07
Duration: 30.20s, Total samples = 322.47s (1067.83%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 202.32s, 62.74% of 322.47s total
Dropped 1048 nodes (cum <= 1.61s)
Showing top 10 nodes out of 238
      flat  flat%   sum%        cum   cum%
   113.54s 35.21% 35.21%    113.55s 35.21%  runtime.stdcall2
    38.99s 12.09% 47.30%     39.04s 12.11%  runtime.stdcall0
    18.34s  5.69% 52.99%     19.52s  6.05%  runtime.cgocall
     7.22s  2.24% 55.23%     10.47s  3.25%  runtime.findObject
     5.58s  1.73% 56.96%      5.58s  1.73%  runtime.stdcall1
     4.75s  1.47% 58.43%      4.75s  1.47%  runtime.(*mspan).base (inline)
     4.35s  1.35% 59.78%      4.35s  1.35%  runtime.procyield
     3.44s  1.07% 60.85%      3.44s  1.07%  runtime.nextFreeFast
     3.12s  0.97% 61.81%     13.21s  4.10%  runtime.scanobject
     2.99s  0.93% 62.74%      2.99s  0.93%  runtime.memmove
(pprof)
*/
/*
PS D:\code\go\wx-example\wx-example\k6\accounts> go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in C:\Users\MSI CYBORG\pprof\pprof.wxapi.exe.samples.cpu.003.pb.gz
File: wxapi.exe
Build ID: D:\code\go\wx-example\wx-example\wxapi\wxapi.exe2025-10-04 21:51:25.7431819 +0700 +07
Type: cpu
Time: 2025-10-04 22:10:52 +07
Duration: 30.20s, Total samples = 321.50s (1064.74%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 202.95s, 63.13% of 321.50s total
Dropped 1042 nodes (cum <= 1.61s)
Showing top 10 nodes out of 237
      flat  flat%   sum%        cum   cum%
   116.52s 36.24% 36.24%    116.52s 36.24%  runtime.stdcall2
    37.40s 11.63% 47.88%     37.45s 11.65%  runtime.stdcall0
    17.75s  5.52% 53.40%     18.98s  5.90%  runtime.cgocall
     7.51s  2.34% 55.73%     10.78s  3.35%  runtime.findObject
     5.98s  1.86% 57.59%      5.98s  1.86%  runtime.stdcall1
     4.45s  1.38% 58.98%      4.45s  1.38%  runtime.(*mspan).base (inline)
     4.11s  1.28% 60.26%      4.11s  1.28%  runtime.procyield
     3.15s  0.98% 61.23%     25.39s  7.90%  github.com/go-sql-driver/mysql.(*textRows).readRow
     3.04s  0.95% 62.18%     18.19s  5.66%  database/sql.convertAssignRows
     3.04s  0.95% 63.13%      3.04s  0.95%  runtime.nextFreeFast (inline)
(pprof) 

*/