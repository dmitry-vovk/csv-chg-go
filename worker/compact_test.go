package worker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompact(t *testing.T) {
	uuid := "767d967f-b55b-4457-bfee-685eaa6d0583"
	expected := compact{
		0x76, 0x7d, 0x96, 0x7f,
		0xb5, 0x5b,
		0x44, 0x57,
		0xbf, 0xee,
		0x68, 0x5e, 0xaa, 0x6d, 0x05, 0x83,
	}
	if id := fromUUID(uuid); assert.Equal(t, expected, id) {
		assert.Equal(t, uuid, id.String())
	}
}

// BenchmarkCompact-12
// 8261254
// 133 ns/op
// 96 B/op
// 2 allocs/op
func BenchmarkCompact(b *testing.B) {
	uuid := "767d967f-b55b-4457-bfee-685eaa6d0583"
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if uuid2 := fromUUID(uuid).String(); uuid2 != uuid {
			b.Fail()
		}
	}
}
