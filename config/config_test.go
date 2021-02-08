package config

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	testCases := []struct {
		config Config
		err    error
	}{
		{
			err: errors.New("no input specified"),
		},
		{
			config: Config{
				CSVFile: "/some/file",
			},
			err: errors.New("no API URL specified"),
		},
		{
			config: Config{
				APIURL:  "ftp://invalid.url",
				CSVFile: "/some/file",
			},
			err: errors.New("invalid API URL"),
		},
		{
			config: Config{
				APIURL:  "http://valid.url",
				CSVFile: "/some/file",
			},
			err: errors.New("workers count should be greater than zero"),
		},
		{
			config: Config{
				APIURL:  "http://valid.url",
				CSVFile: "/some/file",
				Workers: -1,
			},
			err: errors.New("workers count should be greater than zero"),
		},
		{
			config: Config{
				APIURL:  "http://valid.url",
				CSVFile: "/some/file",
				Workers: 1,
			},
			err: errors.New("interval should be at least a second"),
		},
		{
			config: Config{
				APIURL:   "http://valid.url",
				CSVFile:  "/some/file",
				Workers:  1,
				Interval: -1,
			},
			err: errors.New("interval should be at least a second"),
		},
		{
			config: Config{
				APIURL:   "http://valid.url",
				CSVFile:  "/some/file",
				Workers:  1,
				Interval: time.Second,
			},
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, tc.err, tc.config.validate())
		})
	}
}

func TestMustLoad(t *testing.T) {
	keepArgs := os.Args
	defer func() { os.Args = keepArgs }()
	args := []string{
		os.Args[0],
		"-api", "http://example.com",
		"--",
	}
	os.Args = args
	cfg := MustLoad()
	assert.Equal(t, Config{
		APIURL:   "http://example.com",
		CSVFile:  "--",
		Interval: 60 * time.Second,
		Workers:  1,
	}, cfg)
}

// This one is not picked up by test coverage analyser
func TestMustExit(t *testing.T) {
	// Run subprocess that is expected to exit with code 1
	if os.Getenv("MUST_EXIT") == "1" {
		keepArgs := os.Args
		defer func() { os.Args = keepArgs }()
		args := []string{
			os.Args[0],
			"-api", "http://example.com",
		}
		os.Args = args
		MustLoad()
		panic("I am unreachable")
	}
	// Run subprocess and collect its output
	cmd := exec.Command(os.Args[0], "-test.run=TestMustExit")
	cmd.Env = append(os.Environ(), "MUST_EXIT=1")
	out, err := cmd.CombinedOutput()
	if e, ok := err.(*exec.ExitError); assert.True(t, ok) {
		assert.Equal(t, 1, e.ExitCode())
		assert.True(t, bytes.Contains(out, []byte(`Error: no input specified`)))
	}
}
