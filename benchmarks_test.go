package grouter

import (
	//"fmt"
	"github.com/slclub/gnet"
	"github.com/slclub/link"
	"github.com/slclub/utils"
	"net/http"
	"testing"
)

// go test -v -run="none" -bench=.

func BenchmarkPath(B *testing.B) {
	url := "/amine/now"

	path_decoder := NewPath()

	B.ReportAllocs()
	B.ResetTimer()

	path := ""
	for i := 0; i < B.N; i++ {
		_, path, _ = path_decoder.Decode(url)
	}

	link.DEBUG_PRINT("path_decoder:", path, utils.EOL)
}

func BenchmarkMap(B *testing.B) {
	store := NewStore()
	store.Lookup("GET")
	store.Lookup("POST")
	store.Lookup("HEAD")
	store.Lookup("ANY")

	B.ReportAllocs()
	B.ResetTimer()

	for i := 0; i < B.N; i++ {
		store.Lookup("POST")
	}
}

//
func BenchmarkFirst(B *testing.B) {
	app := NewEngine()
	url1 := "/common/static/rest/get/ping"
	app.router.GET(url1, func(ctx gnet.Contexter) {
		//fmt.Println("ping!")
	})
	url2 := "/common/static/ok/r1/r2/r3/r4/r5"
	app.router.GET(url2, func(ctx gnet.Contexter) {
	})

	url3 := "/common/static/book?a=1&b=2"
	app.router.GET(url3, func(ctx gnet.Contexter) {
	})

	//run_request(B, app, "GET", url1)
	run_request(B, app, "GET", url2)
	//run_request(B, app, "GET", url3)
}

type header_writer struct {
	header http.Header
}

func new_mock_writer() *header_writer {
	return &header_writer{
		http.Header{},
	}
}

func (m *header_writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *header_writer) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *header_writer) Header() http.Header {
	return m.header
}

func (m *header_writer) WriteHeader(int) {}

func run_request(B *testing.B, en *Engine, method, path string) {
	// create fake request
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	w := new_mock_writer()
	B.ReportAllocs()
	B.ResetTimer()
	for i := 0; i < B.N; i++ {
		en.ServeHTTP(w, req)
	}
}
