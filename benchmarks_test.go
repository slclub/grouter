package grouter

import (
	//"fmt"
	"net/http"
	"testing"
)

// go test -v -run="none" -bench=.
//
func BenchmarkFirst(B *testing.B) {
	app := NewEngine()
	url1 := "/common/static/rest/get/ping"
	app.router.GET(url1, func(ctx Contexter) {
		//fmt.Println("ping!")
	})
	url2 := "/common/static/ok/r1/r2/r3/r4/r5"
	app.router.GET(url2, func(ctx Contexter) {
	})

	url3 := "/common/static/book?a=1&b=2"

	//run_request(B, app, "GET", url1)
	run_request(B, app, "GET", url2)
	run_request(B, app, "GET", url3)
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
