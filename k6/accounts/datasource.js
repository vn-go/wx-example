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
    vus: 10,
    duration: '60s',

};
let token = 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3NmVjYTUxYy05MjdjLTRkMTMtYjRmNi1iYTVlNTc0NWMxNTAiLCJleHAiOjE3NjA3MjIzMTQsImlhdCI6MTc2MDcwNzkxNCwicm9sZSI6IiIsImVtYWlsIjoiIn0.jIxUnXmfsw2-6sklgFWhwOiZoraZznJTV0eT4TOmUwA'
export default function () {

    // let url = 'http://localhost:8080/api/accounts/get-list-of-roles';
    //let url = 'http://localhost:8080/api/accounts/get-list-of-roles-s-q-l'
    //let url = 'http://localhost:8080/api/pure/get'
    //let url = 'http://localhost:8080/hello'
    //let url = 'http://localhost:8080/api/hz/check'
    //let url = "http://localhost:8080/api/data-source/get"
    //let url = "http://localhost:8080/api/accounts/get-list-of-roles"
    let url = "https://bcb9f8917298.ngrok-free.app/api/data-source/get"
    let headers = {
        'Authorization': token,
        'Content-Type': 'application/json',
    };

    let res = http.post(url, JSON.stringify({
        "name": "user",
        "fields": "",
        "filter": ""
    }), { headers: headers });
    //let res = http.post(url)
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
//k6 run datasource.js
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=300