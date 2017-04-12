## Instructions

1. Running router.ServeHTTP benchmarks
    - go test -bench=. -test.benchmem

2. Running http.ListenAndServe benchmarks
    - Uncomment listen API's - one at a time
    - go test -bench=listen -test.benchmem
    - go get github.com/wg/wrk
    - ./wrk -c100 -d10 -t60 http://127.0.0.1:8080/user/index

## References

1. github.com/gorilla/mux
2. github.com/go-zoo/bone
3. github.com/szxp/mux
4. github.com/julienschmidt/httprouter
5. github.com/julienschmidt/go-http-routing-benchmark
