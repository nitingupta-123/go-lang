package main
import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

// transport is a custom transport that modifies requests and responses.
type transport struct {
	http.RoundTripper
}

// RoundTrip modifies the request and response as needed.
func (t *transport) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	// Set the Host header for the request
	req.Host = "demo.miniorange.com"

	// Perform the actual RoundTrip using the embedded RoundTripper
	resp, err = t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Modify the Location header if present
	if location, ok := resp.Header["Location"]; ok {
		for i, value := range location {
			// Replace "demo.miniorange" with "example" in the Location header
			resp.Header["Location"][i] = strings.Replace(value, "demo.miniorange.com", "example.com", -1)
		}
	}

	// Read the response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Close the original response body
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Replace "demo.miniorange" with "example" in the response body
	b = bytes.Replace(b, []byte("demo.miniorange"), []byte("example"), -1)

	// Print the modified body to the console
	fmt.Println("Modified Response Body:", string(b))

	// Create a new io.ReadCloser from the modified body
	body := io.NopCloser(bytes.NewReader(b))
	resp.Body = body
	resp.ContentLength = int64(len(b))
	resp.Header.Set("Content-Length", strconv.Itoa(len(b)))

	return resp, nil
}

var _ http.RoundTripper = &transport{}

func main() {
	// Parse the target URL
	target, err := url.Parse("https://demo.miniorange.com")
	if err != nil {
		panic(err)
	}

	// Create a reverse proxy with a custom transport
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &transport{http.DefaultTransport}
	proxy.Director = func(req *http.Request) {
		// Set the scheme, host, and Host header for the request
		req.URL.Scheme = "https"
		req.URL.Host = "demo.miniorange.com"
		req.Header.Set("Host", "demo.miniorange.com")
	}

	// Handle requests for the root path using the reverse proxy
	http.Handle("/", proxy)

	// Listen on port 443 with TLS
	err = http.ListenAndServeTLS(":443", "./example.com.crt", "./example.com.key", nil)
	if err != nil {
		panic(err)
	}
}
