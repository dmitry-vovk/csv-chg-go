package config

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"
)

type Config struct {
	APIURL   string
	CSVFile  string
	Interval time.Duration
	Workers  int
}

func (c Config) validate() error {
	if c.CSVFile == "" {
		return errors.New("no input specified")
	}
	if c.APIURL == "" {
		return errors.New("no API URL specified")
	}
	if u, err := url.Parse(c.APIURL); err != nil {
		return fmt.Errorf("API URL parse error: %s", err)
	} else if !(u.Scheme == "http" || u.Scheme == "https") || u.Host == "" {
		return errors.New("invalid API URL")
	}
	if c.Workers < 1 {
		return errors.New("workers count should be greater than zero")
	}
	if c.Interval < time.Second {
		return errors.New("interval should be at least a second")
	}
	return nil
}

func MustLoad() Config {
	cfg := Config{}
	flag.StringVar(&cfg.APIURL, "api", "", "Base API URL")
	flag.StringVar(&cfg.CSVFile, "input", "", "CSV file source path")
	var interval string
	flag.StringVar(&interval, "interval", "60s", "Interval between checks in time.Duration format")
	flag.IntVar(&cfg.Workers, "workers", 1, "Number of parallel API requests")
	flag.Parse()
	if cfg.CSVFile == "" && os.Args[len(os.Args)-1] == "--" {
		cfg.CSVFile = "--"
	}
	// Ignoring the error here as it will be validated later
	cfg.Interval, _ = time.ParseDuration(interval)
	if err := cfg.validate(); err != nil {
		log.Printf("Error: %s", err)
		flag.Usage()
		os.Exit(1)
	}
	return cfg
}
