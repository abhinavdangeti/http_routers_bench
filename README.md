## Instructions

1. Running router.ServeHTTP benchmarks
    - go test -bench=. -test.benchmem

2. Running http.ListenAndServe benchmarks
    - Uncomment listen API's - one at a time
    - go test -bench=listen -test.benchmem
    - go get github.com/wg/wrk
    - ./wrk -c100 -d10 -t60 http://127.0.0.1:8080/user/index
