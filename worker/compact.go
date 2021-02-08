package worker

import (
	"encoding/hex"
)

// compact is 16-byte representation of 36-characters UUID, 56% smaller, still comparable
type compact [16]byte

// fromUUID returns Compact representation of `uuid`
func fromUUID(uuid string) compact {
	c := compact{}
	bid := []byte(uuid)
	// We expect UUID to be always valid here
	_, _ = hex.Decode(c[0:8], bid[0:8])
	_, _ = hex.Decode(c[4:6], bid[9:13])
	_, _ = hex.Decode(c[6:10], bid[14:18])
	_, _ = hex.Decode(c[8:12], bid[19:23])
	_, _ = hex.Decode(c[10:16], bid[24:])
	return c
}

// String returns standard 36-char uuid representation
func (c compact) String() string {
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], c[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], c[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], c[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], c[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], c[10:])
	return string(buf)
}
