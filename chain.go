package httputil

import (
	"net/http"
	"net/http/httputil"
)

// ChainItem request=>response
type ChainItem struct {
	Request    []byte
	Response   []byte
	StatusCode int
}

// GetChain if redirects
func GetChain(r *http.Response) (chain []ChainItem, err error) {
	lastresp := r
	for lastresp != nil {
		lastreq := lastresp.Request
		lastreqDump, err := httputil.DumpRequestOut(r.Request, false)
		if err != nil {
			return nil, err
		}
		lastrespDump, err := httputil.DumpResponse(lastresp, false)
		if err != nil {
			return nil, err
		}
		chain = append(chain, ChainItem{Request: lastreqDump, Response: lastrespDump, StatusCode: lastresp.StatusCode})
		// process next
		lastresp = lastreq.Response
	}
	// reverse the slice in order to have the chain in progressive order
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}

	return
}
