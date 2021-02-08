package api

import (
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	s := startMockServer()
	c := New(s.addr)
	// Test GET requests
	for uuid, r := range getRequests {
		if item, err := c.GetItem(uuid); r.err == nil && assert.NoError(t, err) {
			assert.Equal(t, r.item, item)
		} else {
			assert.IsType(t, r.err, err)
		}
	}
	// Test POST requests
	for uuid, r := range postRequests {
		if err := c.PostAlert(uuid); r.err == nil {
			assert.NoError(t, err)
		} else {
			assert.Equal(t, err, r.err)
		}
	}
	// Stop the server and test unreachable API
	s.stop()
	if _, err := c.GetItem("00000000-0000-0000-0000-000000000200"); assert.Error(t, err) {
		assert.IsType(t, &url.Error{}, err)
	}
	if err := c.PostAlert("00000000-0000-0000-0000-000000000200"); assert.Error(t, err) {
		assert.IsType(t, &url.Error{}, err)
	}
}

type mockResponse struct {
	code        int    // expected status code
	contentType string // expected content type
	body        []byte // expected body
	err         error  // expected error
	item        *Item  // expected item
}

var (
	getRequests = map[string]mockResponse{
		// OK
		"00000000-0000-0000-0000-000000000200": {
			code:        200,
			contentType: "application/json",
			body:        []byte(`{"uuid":"00000000-0000-0000-0000-000000000200", "name": "item name", "quantity": 10}`),
			item: &Item{
				UUID:     "00000000-0000-0000-0000-000000000200",
				Name:     "item name",
				Quantity: 10,
			},
		},
		// Invalid JSON
		"10000000-0000-0000-0000-000000000200": {
			code:        200,
			contentType: "application/json",
			body:        []byte(`{ inVaLid JsOn }`),
			err:         &json.SyntaxError{},
		},
		// Wrong Content-Type header
		"20000000-0000-0000-0000-000000000200": {
			code:        200,
			contentType: "text/html",
			body:        []byte(`{}`),
			err:         ErrInvalidContentType{contentType: "text/html"},
		},
		// Unexpected response field
		"30000000-0000-0000-0000-000000000200": {
			code:        200,
			contentType: "application/json",
			body:        []byte(`{"id":"unknown", "uuid":"00000000-0000-0000-0000-000000000200", "name": "item name", "quantity": 10}`),
			err:         errors.New(`json: unknown field "id"`),
		},
		// Error 400
		"00000000-0000-0000-0000-000000000400": {
			code: 400,
			err:  ErrBadRequest,
			body: nil,
		},
		// Error 500
		"00000000-0000-0000-0000-000000000500": {
			code: 500,
			err:  ErrServerError,
			body: nil,
		},
		// Unexpected error
		"00000000-0000-0000-0000-000000000999": {
			code: 999,
			err:  ErrUnexpectedStatusCode{code: 999},
		},
	}
	postRequests = map[string]mockResponse{
		// Response 201
		"00000000-0000-0000-0000-000000000201": {
			code: 201,
		},
		// Response 400
		"00000000-0000-0000-0000-000000000400": {
			code: 400,
			err:  ErrBadRequest,
		},
		// Response 500
		"00000000-0000-0000-0000-000000000500": {
			code: 500,
			err:  ErrServerError,
		},
		// Unexpected response
		"00000000-0000-0000-0000-000000000999": {
			code: 999,
			err:  ErrUnexpectedStatusCode{code: 999},
		},
	}
)

type mockAPIServer struct {
	server *http.Server
	addr   string
}

func startMockServer() *mockAPIServer {
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := mockAPIServer{
		addr: "http://" + l.Addr().String() + "/",
	}
	server.server = &http.Server{Handler: &server}
	go func() { _ = server.server.Serve(l) }()
	return &server
}

func (m *mockAPIServer) stop() { _ = m.server.Close() }

func (m *mockAPIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if resp, ok := getRequests[strings.TrimPrefix(r.URL.Path, getItemPath)]; ok {
			w.Header().Add("Content-Type", resp.contentType)
			w.WriteHeader(resp.code)
			if _, err := w.Write(resp.body); err != nil {
				panic(err)
			}
			return
		}
		panic("unexpected GET request " + r.URL.Path)
	case "POST":
		if resp, ok := postRequests[strings.TrimPrefix(r.URL.Path, postAlertPath)]; ok {
			w.WriteHeader(resp.code)
			return
		}
		panic("unexpected POST request " + r.URL.Path)
	default:
		panic("unexpected method " + r.Method)
	}
}
