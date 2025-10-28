package httpclient

import (
  "context"
	"bytes"
	"crypto/x509"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type Request struct {
  Method       string
  Endpoint     string
  Body         interface{}
  Headers      map[string]string
  Insecure     bool
  Username     string
  Password     string
  Token        string
  CACert       string
  Timeout      time.Duration
  MaxRedirects int
}

func DoRequest(ctx context.Context, request Request) ([]byte, error) {
  start := time.Now()
  if request.Endpoint == "" {
    return nil, fmt.Errorf("Doing request: URL is nil")
  }
  tflog.Debug(ctx, fmt.Sprintf("Doing %s to %s", request.Method, request.Endpoint))
  tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", request))
  client := &http.Client{
    Timeout: request.Timeout,
    CheckRedirect: func(req *http.Request, via []*http.Request) error {
      if len(via) >= request.MaxRedirects {
        return fmt.Errorf("Stopped after %d redirects", request.MaxRedirects)
      }
      if len(request.Token) > 0 {
        req.Header.Add("Authorization", "Bearer "+request.Token)
      } else if len(request.Username) > 0 && len(request.Password) > 0 {
        s := request.Username + ":" + request.Password
        req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
      }
      return nil
    },
  }
  
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{InsecureSkipVerify: request.Insecure},
    Proxy:           http.ProxyFromEnvironment,
  }
  
  if request.CACert != "" {
    // Get the SystemCertPool, continue with an empty pool on error
    rootCAs, _ := x509.SystemCertPool()
    if rootCAs == nil {
      rootCAs = x509.NewCertPool()
    }
    
    // Append our cert to the system pool
    if ok := rootCAs.AppendCertsFromPEM([]byte(request.CACert)); !ok {
      log.Println("No certs appended, using system certs only")
    }
    transport.TLSClientConfig.RootCAs = rootCAs
  }
  client.Transport = transport
}

func DoGet(req Request) ([]byte, error) {
	start := time.Now()

	if req.Endpoint == "" {
		return nil, fmt.Errorf("Doing get: URL is nil")
	}
  tflog("Getting from ", req.Endpoint)
  tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", req))

	client := &http.Client{
		Timeout: time.Duration(60 * time.Second),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxHTTPRedirect {
				return fmt.Errorf("Stopped after %d redirects", maxHTTPRedirect)
			}
			if len(req.Token) > 0 {
				req.Header.Add("Authorization", "Bearer "+req.Token)
			} else if len(req.Username) > 0 && len(req.Password) > 0 {
				s := req.Username + ":" + req.Password
				req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
			}
			return nil
		},
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: req.Insecure},
		Proxy:           http.ProxyFromEnvironment,
	}

	if req.CACert != "" {
		// Get the SystemCertPool, continue with an empty pool on error
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}

		// Append our cert to the system pool
		if ok := rootCAs.AppendCertsFromPEM([]byte(req.CACert)); !ok {
			log.Println("No certs appended, using system certs only")
		}
		transport.TLSClientConfig.RootCAs = rootCAs
	}
	client.Transport = transport

	r, err := http.NewRequest("GET", req.Endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("Doing get: %v", err)
	}
	if len(req.Token) > 0 {
		r.Header.Add("Authorization", "Bearer "+req.Token)
	} else if len(req.Username) > 0 && len(req.Password) > 0 {
		s := req.Username + ":" + req.Password
		r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s)))
	}
	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Doing get: %v", err)
	}
	defer resp.Body.Close()

	// Timings recorded as part of internal metrics
	log.Println("Time to get req: ", float64((time.Since(start))/time.Millisecond), " ms")

	return ioutil.ReadAll(resp.Body)
}
