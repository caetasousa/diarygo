//go:build integration

// Package factories provides builders that INSERT rows directly into the test
// database, bypassing services. Each factory returns a domain struct reflecting
// the inserted data. Tests compose factories to assemble arbitrary scenarios.
package factories

import (
	"fmt"
	"sync/atomic"
)

// cpfCounter ensures every call to GenerateCPF returns a unique valid CPF
// within a test process. Starts at a base that's unlikely to collide with
// hand-picked CPFs in existing unit tests.
var cpfCounter uint64 = 100_000_000

// GenerateCPF returns a valid CPF (11 digits, correct check digits) unique to
// this process. Useful for factories that need a non-colliding CPF per call.
func GenerateCPF() string {
	n := atomic.AddUint64(&cpfCounter, 1)
	base := fmt.Sprintf("%09d", n%1_000_000_000)
	return base + checkDigits(base)
}

// checkDigits computes the two CPF verifier digits for a 9-digit base.
// Replicates the algorithm in domain.ValidarCPF.
func checkDigits(base string) string {
	d1 := verifier(base, 10)
	d2 := verifier(base+fmt.Sprint(d1), 11)
	return fmt.Sprintf("%d%d", d1, d2)
}

func verifier(s string, startWeight int) int {
	sum := 0
	for i, c := range s {
		sum += int(c-'0') * (startWeight - i)
	}
	rem := sum % 11
	if rem < 2 {
		return 0
	}
	return 11 - rem
}
