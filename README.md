# go-lang

# Reverse Proxy with Request and Response Modification

This Go program demonstrates the implementation of a reverse proxy server using the `net/http` package. The reverse proxy modifies both incoming requests and outgoing responses, specifically replacing occurrences of "demo.miniorange" with "example" in the request headers, response headers, and response body.

## Code Overview

### Custom Transport

The `transport` type is a custom implementation of the `http.RoundTripper` interface. It is responsible for modifying both requests and responses.

- **RoundTrip:** This method intercepts the HTTP request, modifies its host and performs the actual round trip. It then processes the response, replacing relevant information.

### Main Function

The `main` function initializes the reverse proxy and sets up a simple HTTP server.

- **Target URL Parsing:** The target URL (https://demo.miniorange.com) is parsed using `url.Parse`.

- **Reverse Proxy Configuration:**
  - A reverse proxy is created using `httputil.NewSingleHostReverseProxy` with a custom transport.
  - The `Director` function is set to modify the request before it is sent.

- **HTTP Server Configuration:**
  - Requests for the root path ("/") are handled by the reverse proxy.
  - The server listens on port 443 with TLS using `http.ListenAndServeTLS`.

### TLS Configuration

The server is configured to listen on port 443 with TLS, and TLS certificates are provided (example.com.crt and example.com.key).

## Running the Code

To run the code, ensure that the necessary dependencies are installed, and then execute the following:

```bash
go run main.go
