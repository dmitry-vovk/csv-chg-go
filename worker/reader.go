package worker

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

// ReadUUIDs scans `r` for UUIDs, one per line
func (w *Worker) ReadUUIDs(r io.Reader) error {
	rUUID := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	scanner := bufio.NewScanner(r)
	skipped := 0
	start := time.Now()
	for line := 1; scanner.Scan(); line++ {
		uuid := strings.TrimSpace(scanner.Text())
		if !rUUID.MatchString(uuid) {
			log.Printf("Invalid UUID in line %d: %q", line, uuid)
			skipped++
			continue
		}
		compactUUID := fromUUID(uuid)
		if _, ok := w.uuids[compactUUID]; ok {
			log.Printf("Duplicate UUID in line %d: %q", line, uuid)
			skipped++
		} else {
			w.uuids[compactUUID] = struct{}{}
		}
	}
	log.Printf("%d records loaded, %d skipped in %s", len(w.uuids), skipped, time.Since(start))
	return scanner.Err()
}
