package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	{
		e := ErrUnexpectedStatusCode{code: 999}
		if assert.Error(t, e) {
			assert.Equal(t, "unexpected status code 999", e.Error())
		}
	}
	{
		e := ErrInvalidContentType{contentType: "text/html"}
		if assert.Error(t, e) {
			assert.Equal(t, "invalid content type 'text/html'", e.Error())
		}
	}

}
