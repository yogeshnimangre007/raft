package raft

import (
	crand "crypto/rand"
	"fmt"
	"math/rand"
	"time"
)

// UUID is used to generate a random UUID
func UUID() string {
	buf := make([]byte, 16)
	if _, err := crand.Read(buf); err != nil {
		panic(fmt.Errorf("Failed to read random bytes: %v", err))
	}

	return fmt.Sprintf("%08x-%04x-%04x-%04x-%12x",
		buf[0:4],
		buf[4:6],
		buf[6:8],
		buf[8:10],
		buf[10:16])
}

// returns channel with timeout between [d, 2*d)
func randomTimeout(d time.Duration) <-chan time.Time {
	delta := time.Duration(rand.Int63n(int64(d)))
	return time.After(delta + d)
}
