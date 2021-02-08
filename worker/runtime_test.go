package worker

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dmitry-vovk/csv-chg-go/api"
	"github.com/stretchr/testify/assert"
)

type testCase struct {
	quantity   int    // quantity to return
	returnUUID string // UUID to return
	getError   error  // error to return on GET request
	postError  error  // error to return on POST request
}

var testCases = map[string]*testCase{
	"00000000-0000-0000-0000-000000000001": {
		quantity: 10,
	},
	"00000000-0000-0000-0000-000000000002": {
		quantity: 10,
		getError: api.ErrBadRequest,
	},
	"00000000-0000-0000-0000-000000000003": {
		quantity: 5,
	},
	"00000000-0000-0000-0000-000000000004": {
		quantity:   5,
		returnUUID: "00000000-dead-beef-0000-000000000004",
	},
	"00000000-0000-0000-0000-000000000005": {
		quantity: 10,
		getError: api.ErrServerError,
	},
	"00000000-0000-0000-0000-000000000006": {
		quantity:  4,
		postError: api.ErrBadRequest,
	},
	"00000000-0000-0000-0000-000000000007": {
		quantity:  4,
		postError: api.ErrServerError,
	},
}

func TestRuntime(t *testing.T) {
	var ids []string
	for uuid := range testCases {
		ids = append(ids, uuid)
	}
	c := &mockAPIClient{}
	w := New(c)
	w.WithInterval(time.Second)
	assert.NoError(t, w.ReadUUIDs(strings.NewReader(strings.Join(ids, "\n"))))
	// capture error log
	logBuffer := &bytes.Buffer{}
	log.SetOutput(logBuffer)
	log.SetFlags(0)
	defer log.SetOutput(os.Stderr)
	go w.Run()
	// Sleep for little longer than 'runtime' full cycles
	time.Sleep(2*time.Second + 70*time.Millisecond)
	w.Shutdown()
	assert.Equal(t, 12, c.gets)
	assert.Equal(t, 3, c.posts)
	// check some of the error messages
	logString := logBuffer.String()
	assert.Contains(t, logString, `API indicated UUID "00000000-0000-0000-0000-000000000002" not found, removing`)
	assert.Contains(t, logString, `APi returned wrong item, expected "00000000-0000-0000-0000-000000000004", got "00000000-dead-beef-0000-000000000004"`)
	assert.Contains(t, logString, `API error: internal server error`)
}

type mockAPIClient struct {
	gets  int
	posts int
	m     sync.Mutex
}

var _ APIClient = &mockAPIClient{}

func (m *mockAPIClient) GetItem(uuid string) (*api.Item, error) {
	m.m.Lock()
	defer m.m.Unlock()
	time.Sleep(10 * time.Millisecond)
	m.gets++
	if u, ok := testCases[uuid]; ok {
		if u.getError != nil {
			return nil, u.getError
		}
		if u.returnUUID != "" {
			return &api.Item{UUID: u.returnUUID, Quantity: 8}, nil
		}
		return &api.Item{UUID: uuid, Quantity: testCases[uuid].quantity}, nil
	}
	return nil, api.ErrBadRequest
}

func (m *mockAPIClient) PostAlert(uuid string) error {
	m.m.Lock()
	defer m.m.Unlock()
	time.Sleep(10 * time.Millisecond)
	m.posts++
	if u, ok := testCases[uuid]; ok {
		if u.postError != nil {
			return u.postError
		}
	}
	return nil
}
