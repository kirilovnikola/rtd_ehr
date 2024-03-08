package main

import (
	"bytes"
    "fmt"
    "log"
    "net/http"
    "net/url"
    "net/http/httputil"
)

type EHRReverseProxy struct {
    Proxy *httputil.ReverseProxy
}

func NewProxy(requestUrl string) (*EHRReverseProxy, error) {
    url, err := url.Parse(requestUrl)
    if err != nil {
        return nil, err
    }
    p := &EHRReverseProxy{httputil.NewSingleHostReverseProxy(url)}

    // Modify response
    p.Proxy.ModifyResponse = func(r *http.Response) error {
		method := r.Request.Method
		if method == "POST" || method == "PUT" || method == "DELETE"{
		status := r.Status
		if r.StatusCode == http.StatusOK || r.StatusCode == http.StatusCreated {
		location := r.Header.Get("Content-Location")
		fmt.Println("Sent real-time data for Request "+method+" with Response status: "+status+" location: "+location)

		var jsonData = []byte(`{
			"method": "`+r.Request.Method+`",
			"location": "`+r.Header.Get("Content-Location")+`"
		}`)

		z, err := http.NewRequest("POST", "http://localhost:5001/", bytes.NewBuffer(jsonData))
		z.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(z)
		if err != nil {
		panic(err)
		}
		fmt.Println("Response from target: "+resp.Status)
		}
		}
        return nil
    }

    return p, nil
}

func (p *EHRReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    p.Proxy.ServeHTTP(rw, req)
}

func main() {

    rev_proxy, err := NewProxy("http://localhost:8080")
    if err != nil {
        panic(err)
    }
    http.Handle("/", rev_proxy)
    log.Fatal(http.ListenAndServe(":5000", nil))
}