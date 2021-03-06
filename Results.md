# Results

```
Spec MacBook Pro
     Processor 2.3 GHz Intel Core i7
     Memory 16 GB 1600 MHz DDR3
```

## router.ServeHTTP

```
go test -bench=. -test.benchmem

Benchmark_gmux_single_route_serve_GET-8              1000000          1540 ns/op        1056 B/op         11 allocs/op
Benchmark_gmux_single_route_serve_POST-8             1000000          1568 ns/op        1056 B/op         11 allocs/op
Benchmark_gmux_multiple_routes_serve-8               1000000          1953 ns/op         880 B/op          8 allocs/op

Benchmark_bone_single_route_serve_GET-8              2000000           808 ns/op         688 B/op          5 allocs/op
Benchmark_bone_single_route_serve_POST-8             2000000           802 ns/op         688 B/op          5 allocs/op
Benchmark_bone_multiple_routes_serve-8               1000000          2046 ns/op        1149 B/op         16 allocs/op

Benchmark_smux_single_route_serve_GET-8              1000000          1105 ns/op         760 B/op         10 allocs/op
Benchmark_smux_single_route_serve_POST-8             1000000          1112 ns/op         760 B/op         10 allocs/op
Benchmark_smux_multiple_routes_serve-8               1000000          1261 ns/op         976 B/op         10 allocs/op

Benchmark_httprouter_single_route_serve_GET-8       10000000           127 ns/op          32 B/op          1 allocs/op
Benchmark_httprouter_single_route_serve_POST-8      10000000           111 ns/op          32 B/op          1 allocs/op
Benchmark_httprouter_multiple_routes_serve-8         2000000           945 ns/op         848 B/op          7 allocs/op
```

## http.ListenAndServe

```
wrk: ./wrk -c100 -d10 -t60 http://127.0.0.1:8080/user/index
```

### gmux
```
Running 10s test @ http://127.0.0.1:8080/user/123
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.23ms  220.08us   4.55ms   82.63%
    Req/Sec     8.11k   710.05    16.46k    93.23%
  811062 requests in 10.10s, 92.05MB read
Requests/sec:  80310.63
Transfer/sec:      9.11MB
```

### bone
```
Running 10s test @ http://127.0.0.1:8080/user/123
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.20ms  201.09us   4.41ms   80.91%
    Req/Sec     8.34k   492.79     9.75k    58.71%
  837970 requests in 10.10s, 95.10MB read
Requests/sec:  82955.95
Transfer/sec:      9.41MB
```

### smux
```
Running 10s test @ http://127.0.0.1:8080/user/123
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.21ms  195.50us   4.16ms   81.29%
    Req/Sec     8.26k   423.08     9.39k    62.97%
  830461 requests in 10.10s, 94.25MB read
Requests/sec:  82223.67
Transfer/sec:      9.33MB
```

### httprouter
```
Running 10s test @ http://127.0.0.1:8080/user/123
  10 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency     1.19ms  179.25us   4.22ms   76.00%
    Req/Sec     8.41k   519.04     9.93k    66.04%
  844642 requests in 10.10s, 95.86MB read
Requests/sec:  83625.84
Transfer/sec:      9.49MB
```
