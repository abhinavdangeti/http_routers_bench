package routerCheck

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
	"testing"

	gmux "github.com/gorilla/mux"
	httprouter "github.com/julienschmidt/httprouter"
	bone "github.com/go-zoo/bone"
	smux "github.com/szxp/mux"
)

// common
type route struct {
	method string
	path string
}

var routes = []route{
	{"GET", "/user/:id"},
	{"GET", "/user/:id/name/:lastname"},
	{"GET", "/user/:id/age"},
	{"GET", "/user/:id/sex"},
	{"POST", "/user/:id/address/:zip"},
	{"GET", "/user"},
	{"DELETE", "/user/:id/address/:country"},
	{"GET", "/user/:id/height"},
	{"POST", "/user/:id/phone"},
	{"GET", "/user/:id/weight"},
}

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

func benchServeSingleRoute(b *testing.B, router http.Handler, r *http.Request) {
	w := new(mockResponseWriter)
	u := r.URL
	r.RequestURI = u.RequestURI()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
	}
}

func benchServeMultipleRoutes(b *testing.B, router http.Handler) {
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/", nil)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r.Method = routes[i%len(routes)].method
		r.RequestURI = routes[i%len(routes)].path
		router.ServeHTTP(w, r)
	}
}

// ===================== gmux =======================
func gmuxWrite(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, r.RequestURI)
}

func loadGmuxSingleRoute(method, path string) http.Handler {
	m := gmux.NewRouter()
	m.HandleFunc(path, gmuxWrite).Methods(method)
	return m
}

func loadGmuxMultipleRoutes() http.Handler {
	m := gmux.NewRouter()
	re := regexp.MustCompile(":([^/]*)")
	for _, route := range routes {
		m.HandleFunc(
			re.ReplaceAllString(route.path, "{$1}"),
			gmuxWrite,
		).Methods(route.method)
	}
	return m
}

func Benchmark_gmux_single_route_serve_GET(b *testing.B) {
	router := loadGmuxSingleRoute("GET", "/user/{name}")
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServeSingleRoute(b, router, r)
}

func Benchmark_gmux_single_route_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadGmuxSingleRoute("POST", "/user/{name}")
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServeSingleRoute(b, router, r)
}

func Benchmark_gmux_multiple_routes_serve(b *testing.B) {
	router := loadGmuxMultipleRoutes()
	benchServeMultipleRoutes(b, router)
}

/*
func Benchmark_gmux_listen(b *testing.B) {
	router := loadGmuxSingleRoute("GET", "/user/{name}")
	http.ListenAndServe(":8080", router)
}
*/

// ===================== bone =======================
func boneWrite(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, bone.GetValue(r, "name"))
}

func addToBone(method, path string, b *bone.Mux) *bone.Mux {
	switch method {
	case "GET":
		b.Get(path, http.HandlerFunc(boneWrite))
	case "POST":
		b.Post(path, http.HandlerFunc(boneWrite))
	case "PUT":
		b.Put(path, http.HandlerFunc(boneWrite))
	case "PATCH":
		b.Patch(path, http.HandlerFunc(boneWrite))
	case "DELETE":
		b.Delete(path, http.HandlerFunc(boneWrite))
	default:
		panic("What is this method!? ..." + method)
	}
	return b
}

func loadBoneSingleRoute(method, path string) http.Handler {
	b := bone.New()
	return addToBone(method, path, b)
}

func loadBoneMultipleRoutes() http.Handler {
	b := bone.New()
	for _, route := range routes {
		b = addToBone(route.method, route.path, b)
	}
	return b
}

func Benchmark_bone_single_route_serve_GET(b *testing.B) {
	router := loadBoneSingleRoute("GET", "/user/:name")
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServeSingleRoute(b, router, r)
}

func Benchmark_bone_single_route_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadBoneSingleRoute("POST", "/user/:name")
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServeSingleRoute(b, router, r)
}

func Benchmark_bone_multiple_routes_serve(b *testing.B) {
	router := loadBoneMultipleRoutes()
	benchServeMultipleRoutes(b, router)
}

/*
func Benchmark_bone_listen(b *testing.B) {
	router := loadBoneSingleRoute("GET", "/user/:name")
	http.ListenAndServe(":8080", router)
}
*/

// ===================== smux =======================
func smuxWrite(rw http.ResponseWriter, r *http.Request) {
	io.WriteString(rw, r.Context().Value(smux.CtxKey("name")).(string))
}

func loadSmuxSingleRoute(method, path string) http.Handler {
	m := smux.NewMuxer()
	m.HandleFunc(path, smuxWrite, method)
	return m
}

func loadSmuxMultipleRoutes() http.Handler {
	m := smux.NewMuxer()
	for _, route := range routes {
		m.HandleFunc(route.path, smuxWrite, route.method)
	}
	return m
}

func Benchmark_smux_single_route_serve_GET(b *testing.B) {
	router := loadSmuxSingleRoute("GET", "/user/:name")
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServeSingleRoute(b, router, r)
}

func Benchmark_smux_single_route_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadSmuxSingleRoute("POST", "/user/:name")
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServeSingleRoute(b, router, r)
}

func Benchmark_smux_multiple_routes_serve(b *testing.B) {
	router := loadSmuxMultipleRoutes()
	benchServeMultipleRoutes(b, router)
}

/*
func Benchmark_smux_listen(b *testing.B) {
	router := loadSmuxSingleRoute("GET", "/user/:name")
	http.ListenAndServe(":8080", router)
}
*/

// ================== httprouter ====================
func httprouterWrite(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	io.WriteString(rw, p.ByName("name"))
}

func loadHttprouterSingleRoute(method, path string) http.Handler {
	hr := httprouter.New()
	hr.Handle(method, path, httprouterWrite)
	return hr
}

func loadHttprouterMultipleRoutes() http.Handler {
	hr := httprouter.New()
	for _, route := range routes {
		hr.Handle(route.method, route.path, httprouterWrite)
	}
	return hr
}

func Benchmark_httprouter_single_route_serve_GET(b *testing.B) {
	router := loadHttprouterSingleRoute("GET", "/user/:name")
	r, _ := http.NewRequest("GET", "/user/index", nil)
	benchServeSingleRoute(b, router, r)
}

func Benchmark_httprouter_single_route_serve_POST(b *testing.B) {
	index := "primary_index"
	router := loadHttprouterSingleRoute("POST", "/user/:name")
	r, _ := http.NewRequest("POST", "/user/index", bytes.NewReader([]byte(index)))
	benchServeSingleRoute(b, router, r)
}

func Benchmark_httprouter_multiple_routes_serve(b *testing.B) {
	router := loadHttprouterMultipleRoutes()
	benchServeMultipleRoutes(b, router)
}

/*
func Benchmark_httprouter_listen(b *testing.B) {
	router := loadHttprouterSingleRoute("GET", "/user/:name")
	http.ListenAndServe(":8080", router)
}
*/
