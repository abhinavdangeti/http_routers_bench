package routerCheck

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	gmux "github.com/gorilla/mux"
	httprouter "github.com/julienschmidt/httprouter"
	bone "github.com/go-zoo/bone"
	smux "github.com/szxp/mux"
)

// common
type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func benchServe(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	rq := u.RawQuery
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		u.RawQuery = rq
		router.ServeHTTP(w, r)
	}
}

// ===================== gmux =======================
func gmuxWrite(rw http.ResponseWriter, r *http.Request) {
	params := gmux.Vars(r)
	io.WriteString(rw, params["name"])
}

func loadgmux(method, path string, handler http.HandlerFunc) http.Handler {
	m := gmux.NewRouter()
	m.HandleFunc(path, handler).Methods(method)
	return m
}

func Benchmark_gmux_serve_GET(b *testing.B) {
	router := loadgmux("GET", "/user/{name}", gmuxWrite)
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServe(b, router, r)
}

func Benchmark_gmux_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadgmux("POST", "/user/{name}", gmuxWrite)
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServe(b, router, r)
}

/*
func Benchmark_gmux_listen(b *testing.B) {
	router := loadgmux("GET", "/user/{name}", gmuxWrite)
	http.ListenAndServe(":8080", router)
}
*/

// ===================== bone =======================
func boneWrite(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, bone.GetValue(r, "name"))
}

func loadbone(method, path string, handler http.Handler) http.Handler {
	b := bone.New()
	switch method {
	case "GET":
		b.Get(path, handler)
	case "POST":
		b.Post(path, handler)
	case "PUT":
		b.Put(path, handler)
	case "PATCH":
		b.Patch(path, handler)
	case "DELETE":
		b.Delete(path, handler)
	default:
		panic("What is this method!? ..." + method)
	}
	return b
}

func Benchmark_bone_serve_GET(b *testing.B) {
	router := loadbone("GET", "/user/:name", http.HandlerFunc(boneWrite))
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServe(b, router, r)
}

func Benchmark_bone_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadbone("POST", "/user/:name", http.HandlerFunc(boneWrite))
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServe(b, router, r)
}

/*
func Benchmark_bone_listen(b *testing.B) {
	router := loadbone("GET", "/user/:name", http.HandlerFunc(boneWrite))
	http.ListenAndServe(":8080", router)
}
*/

// ===================== smux =======================
func smuxWrite(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, r.Context().Value(smux.CtxKey("name")).(string))
}

func loadsmux(method, path string, handler http.HandlerFunc) http.Handler {
	m := smux.NewMuxer()
	m.HandleFunc(path, handler, method)
	return m
}

func Benchmark_smux_serve_GET(b *testing.B) {
	router := loadsmux("GET", "/user/:name", smuxWrite)
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServe(b, router, r)
}

func Benchmark_smux_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadsmux("POST", "/user/:name", smuxWrite)
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServe(b, router, r)
}

/*
func Benchmark_smux_listen(b *testing.B) {
	router := loadsmux("GET", "/user/:name", smuxWrite)
	http.ListenAndServe(":8080", router)
}
*/

// ================== httprouter ====================
func httprouterWrite(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(rw, p.ByName("name"))
}

func loadhttprouter(method, path string, handle httprouter.Handle) http.Handler {
	hr := httprouter.New()
	hr.Handle(method, path, handle)
	return hr
}

func Benchmark_httprouter_serve_GET(b *testing.B) {
	router := loadhttprouter("GET", "/user/:name", httprouterWrite)
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServe(b, router, r)
}

func Benchmark_httprouter_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadhttprouter("POST", "/user/:name", httprouterWrite)
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServe(b, router, r)
}

/*
func Benchmark_httprouter_listen(b *testing.B) {
	router := loadhttprouter("GET", "/user/:name", httprouterWrite)
	http.ListenAndServe(":8080", router)
}
*/
