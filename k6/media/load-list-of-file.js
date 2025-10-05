//http://localhost:8080/api/media/list



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

export default function () {

    //let url = 'http://localhost:8080/api/media/list';
    let url = "http://localhost:8080/api/roles/GetListOfRoles"

    // let headers = {
    //     'Authorization': token,
    //     'Content-Type': 'application/json',
    // };

    let res = http.get(url);
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

}
// k6 run load-list-of-file.js
//-- version 1
//go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
/*
--- version 1---
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in C:\Users\MSI CYBORG\pprof\pprof.wx.exe.samples.cpu.001.pb.gz
07
Type: cpu
Time: 2025-10-04 19:18:10 +07
Duration: 30.18s, Total samples = 3035.77s (10059.64%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 2888.25s, 95.14% of 3035.77s total
Dropped 825 nodes (cum <= 15.18s)
Showing top 10 nodes out of 66
      flat  flat%   sum%        cum   cum%
  2840.24s 93.56% 93.56%   2840.89s 93.58%  runtime.cgocall
    18.49s  0.61% 94.17%     18.49s  0.61%  runtime.stdcall2
     7.01s  0.23% 94.40%     33.32s  1.10%  internal/filepathlite.Clean
     5.72s  0.19% 94.59%     20.06s  0.66%  syscall.UTF16ToString
     4.63s  0.15% 94.74%     31.73s  1.05%  runtime.mallocgcSmallNoscan
     3.61s  0.12% 94.86%     21.84s  0.72%  runtime.scanobject
     3.44s  0.11% 94.97%     49.84s  1.64%  runtime.mallocgc
     1.86s 0.061% 95.03%    163.37s  5.38%  os.(*File).readdir
     1.68s 0.055% 95.09%     17.33s  0.57%  strings.(*Builder).WriteString
     1.57s 0.052% 95.14%     36.06s  1.19%  media/controllers/media.(*Media).ListAllFolderAndFiles.func1    
(pprof)
--- version 2 lan 1----
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 19:37:54.301224 +0700 +07
Type: cpu
Time: 2025-10-04 19:38:18 +07
Duration: 30.13s, Total samples = 2057.69s (6829.75%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 1915.42s, 93.09% of 2057.69s total
Dropped 787 nodes (cum <= 10.29s)
Showing top 10 nodes out of 77
      flat  flat%   sum%        cum   cum%
  1838.33s 89.34% 89.34%   1838.92s 89.37%  runtime.cgocall
    22.74s  1.11% 90.44%     22.74s  1.11%  runtime.stdcall2
    13.68s  0.66% 91.11%     14.74s  0.72%  encoding/json.appendString[go.shape.string]
     8.74s  0.42% 91.53%     35.08s  1.70%  internal/filepathlite.Clean
     7.52s  0.37% 91.90%     13.55s  0.66%  internal/filepathlite.(*lazybuf).append (inline)
     6.14s   0.3% 92.20%     15.96s  0.78%  runtime.concatstrings
     5.56s  0.27% 92.47%     22.20s  1.08%  syscall.UTF16ToString
     4.69s  0.23% 92.70%     35.77s  1.74%  runtime.mallocgcSmallNoscan
     4.08s   0.2% 92.89%     24.12s  1.17%  runtime.scanobject
     3.94s  0.19% 93.09%     56.95s  2.77%  runtime.mallocgc
(pprof)
--- lan 2 ----
Fetching profile over HTTP from http://localhost:6060/debug/pprof/profile?seconds=30
Saved profile in C:\Users\MSI CYBORG\pprof\pprof.wx.exe.samples.cpu.003.pb.gz
File: wx.exe
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 19:37:54.301224 +0700 +07
Type: cpu
Time: 2025-10-04 19:48:02 +07
Duration: 30.12s, Total samples = 2277.61s (7562.37%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 2138.92s, 93.91% of 2277.61s total
Dropped 797 nodes (cum <= 11.39s)
Showing top 10 nodes out of 72
      flat  flat%   sum%        cum   cum%
  2060.47s 90.47% 90.47%   2060.91s 90.49%  runtime.cgocall
    23.39s  1.03% 91.49%     23.39s  1.03%  runtime.stdcall2
    13.43s  0.59% 92.08%     14.41s  0.63%  encoding/json.appendString[go.shape.string]
     9.01s   0.4% 92.48%     37.17s  1.63%  internal/filepathlite.Clean
     8.40s  0.37% 92.85%     14.82s  0.65%  internal/filepathlite.(*lazybuf).append (inline)
     6.04s  0.27% 93.11%     14.83s  0.65%  runtime.concatstrings
     5.44s  0.24% 93.35%     21.14s  0.93%  syscall.UTF16ToString
     4.81s  0.21% 93.56%     34.25s  1.50%  runtime.mallocgcSmallNoscan
     4.11s  0.18% 93.74%     54.95s  2.41%  runtime.mallocgc
     3.82s  0.1
*/
/*
-- version 1---
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 20:01:48.8000865 +0700 +07
Type: cpu
Time: 2025-10-04 20:02:30 +07
Duration: 30.13s, Total samples = 742.99s (2466.26%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 695.61s, 93.62% of 742.99s total
Dropped 804 nodes (cum <= 3.71s)
Showing top 10 nodes out of 55
      flat  flat%   sum%        cum   cum%
   495.03s 66.63% 66.63%    496.08s 66.77%  runtime.cgocall
   153.88s 20.71% 87.34%    164.19s 22.10%  encoding/json.appendString[go.shape.string]
    10.69s  1.44% 88.78%     10.69s  1.44%  runtime.memmove
     8.97s  1.21% 89.98%      8.97s  1.21%  runtime.stdcall2
     6.43s  0.87% 90.85%    186.90s 25.16%  encoding/json.stringEncoder
     5.60s  0.75% 91.60%      7.74s  1.04%  bytes.(*Buffer).Write
     4.91s  0.66% 92.26%    205.44s 27.65%  encoding/json.arrayEncoder.encode
     4.89s  0.66% 92.92%      6.34s  0.85%  bytes.(*Buffer).WriteByte
     4.54s  0.61% 93.53%      7.28s  0.98%  reflect.Value.Index
     0.67s  0.09% 93.62%    458.90s 61.76%  reflect.Value.call
(pprof)
--- version 2 kg dung json-iterator---
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 20:01:48.8000865 +0700 +07
Type: cpu
Time: 2025-10-04 20:03:54 +07
Duration: 30.12s, Total samples = 289.78s (962.16%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 250.02s, 86.28% of 289.78s total
Dropped 687 nodes (cum <= 1.45s)
Showing top 10 nodes out of 93
      flat  flat%   sum%        cum   cum%
   117.38s 40.51% 40.51%    125.01s 43.14%  encoding/json.appendString[go.shape.string]
    49.21s 16.98% 57.49%     50.06s 17.28%  runtime.cgocall
    31.71s 10.94% 68.43%     31.72s 10.95%  runtime.stdcall0
    18.80s  6.49% 74.92%     18.80s  6.49%  runtime.stdcall2
    12.61s  4.35% 79.27%     12.61s  4.35%  runtime.memmove
     4.73s  1.63% 80.90%    142.88s 49.31%  encoding/json.stringEncoder
     4.48s  1.55% 82.45%     11.06s  3.82%  bytes.(*Buffer).Write
     3.91s  1.35% 83.80%    157.25s 54.27%  encoding/json.arrayEncoder.encode
     3.74s  1.29% 85.09%      4.84s  1.67%  bytes.(*Buffer).WriteByte
     3.45s  1.19% 86.28%      5.62s  1.94%  reflect.Value.Index
(pprof) 
--- dung json-iterator ---
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 20:41:03.1021 +0700 +07
Type: cpu
Time: 2025-10-04 20:41:20 +07
Duration: 30.17s, Total samples = 110.08s (364.82%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 88.23s, 80.15% of 110.08s total
Dropped 693 nodes (cum <= 0.55s)
Showing top 10 nodes out of 128
      flat  flat%   sum%        cum   cum%
    43.97s 39.94% 39.94%     44.37s 40.31%  runtime.cgocall
    29.13s 26.46% 66.41%     29.99s 27.24%  github.com/json-iterator/go.(*Stream).WriteStringWithHTMLEscaped
     4.66s  4.23% 70.64%      4.66s  4.23%  runtime.stdcall2
     3.23s  2.93% 73.57%      3.23s  2.93%  runtime.stdcall1
     2.64s  2.40% 75.97%      2.64s  2.40%  runtime.memmove
     1.36s  1.24% 77.21%      1.37s  1.24%  runtime.stdcall0
     0.95s  0.86% 78.07%     30.98s 28.14%  github.com/json-iterator/go.(*htmlEscapedStringEncoder).Encode 
     0.95s  0.86% 78.93%      0.95s  0.86%  runtime.stdcall6
     0.74s  0.67% 79.61%      2.48s  2.25%  runtime.scanobject
     0.60s  0.55% 80.15%     33.49s 30.42%  github.com/json-iterator/go.(*sliceEncoder).Encode
(pprof)
--- lan 2 dung json-iterator cai thien dung var jsonLib = jsoniter.ConfigFastest---
Build ID: D:\code\go\wx-example\wx-example\media\wx\wx.exe2025-10-04 20:45:58.0408396 +0700 +07
Type: cpu
Time: 2025-10-04 20:46:18 +07
Duration: 30.17s, Total samples = 322.65s (1069.51%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top10
Showing nodes accounting for 276.56s, 85.72% of 322.65s total
Dropped 783 nodes (cum <= 1.61s)
Showing top 10 nodes out of 127
      flat  flat%   sum%        cum   cum%
   119.71s 37.10% 37.10%    119.71s 37.10%  runtime.stdcall2
    65.51s 20.30% 57.41%     67.84s 21.03%  github.com/json-iterator/go.(*Stream).WriteString
    39.69s 12.30% 69.71%     39.76s 12.32%  runtime.stdcall0
    28.30s  8.77% 78.48%     29.50s  9.14%  runtime.cgocall
     6.29s  1.95% 80.43%      6.29s  1.95%  runtime.memmove
     5.77s  1.79% 82.22%      5.77s  1.79%  runtime.stdcall1
     3.84s  1.19% 83.41%      3.84s  1.19%  runtime.procyield
     2.86s  0.89% 84.29%      3.55s  1.10%  runtime.findObject
     2.36s  0.73% 85.02%      9.04s  2.80%  runtime.scanobject
     2.23s  0.69% 85.72%     70.36s 21.81%  github.com/json-iterator/go.(*stringCodec).Encode
(pprof)
*/