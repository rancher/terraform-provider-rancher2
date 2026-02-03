package client

// Client is the interface for a client that can make requests to the Rancher API.
type Client interface {
	Do(req *Request, resp *Response) error
	Set(client Client) (Client, error)
	GetApiUrl() string
}

// Request is the request object for the client.
type Request struct {
	Method   string
	Endpoint string
	Body     any // this will be marshalled to json
	Headers  map[string]string
}

func (r *Request) Set(req Request) *Request {
	r.Method = req.Method
	r.Endpoint = req.Endpoint
	r.Body = req.Body
	r.Headers = req.Headers
	return r
}

// Response is the response object from the client.
type Response struct {
	Body       []byte
	Headers    map[string][]string
	StatusCode int
}

func (r *Response) Set(resp Response) *Response {
	r.Body = resp.Body
	r.Headers = resp.Headers
	r.StatusCode = resp.StatusCode
	return r
}

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return e.Message
}
