package source

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expect = []byte("Text file\n")

func TestReadAny(t *testing.T) {
	t.Run("non existing", func(t *testing.T) {
		r := reader{}
		if err := ReadAny("test_data/non-existing.txt", r.read); assert.Error(t, err) {
			assert.IsType(t, &os.PathError{}, err)
		}
	})
	t.Run("bad gzip", func(t *testing.T) {
		r := reader{}
		if err := ReadAny("test_data/bad-gzip.gz", r.read); assert.Error(t, err) {
			assert.Equal(t, gzip.ErrHeader, err)
		}
	})
	t.Run("plain text", func(t *testing.T) {
		r := reader{}
		if err := ReadAny("test_data/file.txt", r.read); assert.NoError(t, err) {
			assert.Equal(t, expect, r.content)
		}
	})
	t.Run("gzipped", func(t *testing.T) {
		r := reader{}
		if err := ReadAny("test_data/file.txt.gz", r.read); assert.NoError(t, err) {
			assert.Equal(t, expect, r.content)
		}
	})
	t.Run("bzipped", func(t *testing.T) {
		r := reader{}
		if err := ReadAny("test_data/file.txt.bz2", r.read); assert.NoError(t, err) {
			assert.Equal(t, expect, r.content)
		}
	})
	s := startMockServer()
	t.Run("bad code", func(t *testing.T) {
		r := reader{}
		if err := ReadAny(s.Address()+"bad", r.read); assert.Error(t, err) {
			assert.Equal(t, "bad response code 400", err.Error())
		}
	})
	t.Run("http", func(t *testing.T) {
		r := reader{}
		if err := ReadAny(s.Address(), r.read); assert.NoError(t, err) {
			assert.Equal(t, expect, r.content)
		}
	})
	s.Stop()
	t.Run("bad server", func(t *testing.T) {
		r := reader{}
		if err := ReadAny(s.Address(), r.read); assert.Error(t, err) {
			assert.IsType(t, &url.Error{}, err)
		}
	})
	t.Run("stdin", func(t *testing.T) {
		tmp, err := ioutil.TempFile("", "example")
		if err != nil {
			panic(err)
		}
		defer func() { _ = os.Remove(tmp.Name()) }()
		if _, err = tmp.Write(expect); err != nil {
			panic(err)
		}
		if _, err = tmp.Seek(0, 0); err != nil {
			panic(err)
		}
		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()
		os.Stdin = tmp
		r := reader{}
		if err = ReadAny("--", r.read); assert.NoError(t, err) {
			assert.Equal(t, expect, r.content)
		}
	})
}

type reader struct {
	content []byte
}

func (r *reader) read(ior io.Reader) (err error) {
	r.content, err = ioutil.ReadAll(ior)
	return
}

type mockServer struct {
	server  *http.Server
	address string
}

func startMockServer() *mockServer {
	l, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	server := mockServer{
		address: "http://" + l.Addr().String() + "/",
	}
	server.server = &http.Server{Handler: &server}
	go func() { _ = server.server.Serve(l) }()
	return &server
}
func (m mockServer) Address() string { return m.address }
func (m *mockServer) Stop()          { _ = m.server.Close() }

func (m *mockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/bad" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, _ = w.Write(expect)
}
