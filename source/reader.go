package source

import (
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ReadAny attempts to call `fn` with `io.Reader` made from `src`.
// `src` may be a path to a local file (can be compressed gzip/bzip2),
// or URL, or `--` for OS stdin stream
func ReadAny(src string, fn func(io.Reader) error) error {
	// StdIn
	if src == "--" {
		return fn(os.Stdin)
	}
	// Remote URL
	if strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "http://") {
		resp, err := http.Get(src)
		if err != nil {
			return err
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode == http.StatusOK {
			return fn(resp.Body)
		}
		return errors.New("bad response code " + strconv.Itoa(resp.StatusCode))
	}
	// Local file
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	// Compressed?
	if strings.HasSuffix(src, ".gz") {
		z, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer func() { _ = z.Close() }()
		return fn(z)
	} else if strings.HasSuffix(src, ".bz2") {
		return fn(bzip2.NewReader(f))
	}
	return fn(f)
}
