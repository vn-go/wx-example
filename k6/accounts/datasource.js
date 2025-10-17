/*
{
    "name": "role",
    "fields": "left(code,2) C,count(id) total,day(createdOn) d",
    "filter": "total>0"
}
*/
import { check } from 'k6'; // 👈 thêm dòng này
import http from 'k6/http';


export const options = {
    //executor: 'constant-vus',
    vus: 200,
    duration: '60s',

};
let token = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJhOTQ3YWQzYy1lNzQ3LTQ3YzctODc4ZS02YWIyMjlmNjQyNjIiLCJleHAiOjE3NTk2NTU0NzMsImlhdCI6MTc1OTY0MTA3Mywicm9sZSI6IiIsImVtYWlsIjoiIn0.RwP9vbM-BQrdoLjyL-2HYYAbJ7K0I-AKksFgPiB6BrA'
export default function () {

    // let url = 'http://localhost:8080/api/accounts/get-list-of-roles';
    //let url = 'http://localhost:8080/api/accounts/get-list-of-roles-s-q-l'
    //let url = 'http://localhost:8080/api/pure/get'
    //let url = 'http://localhost:8080/hello'
    //let url = 'http://localhost:8080/api/hz/check'
    let url = "http://localhost:8080/api/data-source/get"
    //let url = "http://localhost:8080/api/accounts/get-list-of-roles"
    let headers = {
        'Authorization': token,
        'Content-Type': 'application/json',
    };

    let res = http.post(url, JSON.stringify({
        "name": "role",
        "fields": "left(code,2) C,count(id) total,day(createdOn) d",
        "filter": "total>0"
    }), { headers: headers });
    //let res = http.post(url)
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
//k6 run datasource.js
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=300