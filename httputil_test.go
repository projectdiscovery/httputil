package httputil

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDumpRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)

	reqdump, err := DumpRequest(req)
	require.Nil(t, err)
	exp := "GET /foo HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n"
	require.Equal(t, exp, reqdump)
}

func TestDumpResponseHeadersAndRaw(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Del("Date")
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	require.Nil(t, err)

	headersdump, respdump, err := DumpResponseHeadersAndRaw(res)
	headersdump = strings.Split(headersdump, "Date")[0]
	tokens := strings.Split(respdump, "\r\n")
	respdump = ""
	for _, token := range tokens {
		if !strings.HasPrefix(token, "Date") {
			respdump += token + "\r\n"
		}
	}
	require.Nil(t, err)
	headers := "HTTP/1.1 200 OK\r\nContent-Length: 14\r\nContent-Type: text/plain; charset=utf-8\r\n"
	resp := "HTTP/1.1 200 OK\r\nContent-Length: 14\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nHello, client\n\r\n"
	require.Equal(t, headers, headersdump)
	require.Equal(t, resp, respdump)
}
