package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	// "net/http/httputil"
	"net/url"

	"golang.org/x/net/http2"
)

func main() {
}

func initclient() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func fosdemClient() {
	// demoUrl, err := url.Parse("http://172.17.0.2")
	demoUrl, err := url.Parse("https://172.17.0.2")
	if err != nil {
		log.Fatal(err)
	}
	// proxy := httputil.NewSingleHostReverseProxy(demoUrl)
	proxy := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.Host = demoUrl.Host
		req.URL.Host = demoUrl.Host
		req.URL.Scheme = demoUrl.Scheme
		req.RequestURI = "" // Needs to be cleaned or may be rejected
		// Use the caller ip
		s, _, _ := net.SplitHostPort(req.RemoteAddr)
		req.Header.Set("X-Forwarded-For", s)

		// Tell client you can do http2
		http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, err)
			return
		}
		// Copy content type as well as other headers.
		for key, values := range resp.Header {
			for _, value := range values {
				rw.Header().Set(key, value) // Maybe Header().Add?
			}
		}

		// Flush , stream
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-time.Tick(10 * time.Millisecond):
					rw.(http.Flusher).Flush()
				case <-done:
					return
				}
			}
		}()

		// Trailer
		trailerKeys := []string{}
		for key := range resp.Trailer {
			trailerKeys = append(trailerKeys, key)
		}
		rw.Header().Set("Trailer", strings.Join(trailerKeys, ","))

		// Build response
		rw.WriteHeader(resp.StatusCode)
		// Trailer part2 set Trailer header value
		for key, values := range resp.Trailer {
			for _, value := range values {
				rw.Header().Set(key, value)
			}
		}

		io.Copy(rw, resp.Body)
		close(done)
	})

	// Http2 -> TLS ALPN
	http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", proxy)

	// http.ListenAndServe(":8080", proxy)
}

func simpleProxy() {
	listener, err := net.Listen("tcp", ":1111")
	if err != nil {
		log.Fatalln("Unable to bind on port")
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Unable to accept connection")
		}
		go handleConnection(conn)
	}
}

func handleConnection(src net.Conn) {
	dst, err := net.Dial("tcp", "google.com:443") // "http://www.google.com:80")
	if err != nil {
		log.Fatalln("Unable to connect to target server.")
	}
	defer dst.Close()

	// Copy the output to destination
	go func() {
		if _, err := io.Copy(dst, src); err != nil {
			log.Fatalln(err)
		}
	}()

	// Copy the response from target server conn to src conn
	if _, err := io.Copy(src, dst); err != nil {
		log.Fatalln(err)
	}
}
