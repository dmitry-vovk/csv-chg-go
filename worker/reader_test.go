package worker

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerReader(t *testing.T) {
	w := New(nil)
	f, err := os.Open("test_data/file.csv")
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()
	// capture error log
	logBuffer := &bytes.Buffer{}
	log.SetOutput(logBuffer)
	log.SetFlags(0)
	if assert.NoError(t, w.ReadUUIDs(f)) {
		assert.Equal(t, 3, len(w.uuids))
	}
	// check for logged error messages
	logString := logBuffer.String()
	assert.Contains(t, logString, `Invalid UUID in line 3: "ee88ff32-f753-4a49-abf1-2885fdfcafbaee88ff32-f753-4a49-abf1-2885fdfcafba"`)
	assert.Contains(t, logString, `Invalid UUID in line 4: "ee88ff32-f753-4a49-abf1-2885fdfcafbz"`)
	assert.Contains(t, logString, `Invalid UUID in line 5: ""`)
	assert.Contains(t, logString, `Invalid UUID in line 7: "..."`)
	assert.Contains(t, logString, `Duplicate UUID in line 8: "9E2CB4dd-bd6e-48aa-9c0d-696a058226ed"`)
	assert.Contains(t, logString, `3 records loaded, 5 skipped in `)
	w.running = true
	if err := w.ReadUUIDs(f); assert.Error(t, err) {
		assert.Equal(t, "worker already running", err.Error())
	}
}
