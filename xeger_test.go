package xeger_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/takahiromiyamoto/go-xeger"
)

func Test_Generate(t *testing.T) {
	expr := `[a-zA-Z0-9]{32}`
	re := regexp.MustCompile(expr)
	x, err := xeger.NewXeger(expr)
	require.NoError(t, err)

	s1 := x.Generate()
	require.Regexp(t, re, s1)

	s2 := x.Generate()
	require.Regexp(t, re, s2)

	require.NotEqual(t, s1, s2, "Expect strings to differ when using default RandSource")
}

func Test_GenerateWithSource(t *testing.T) {
	expr := `[a-zA-Z0-9]{32}`
	re := regexp.MustCompile(expr)
	x, err := xeger.NewXeger(expr)
	require.NoError(t, err)

	rng := constantInt{n: 42}

	s1 := x.GenerateWithSource(rng)
	require.Regexp(t, re, s1)

	s2 := x.GenerateWithSource(rng)
	require.Regexp(t, re, s2)

	require.Equal(t, s1, s2, "Expect strings to be equal when source of random numbers is the repeatable.")
}

// constantInt implements RandSource but always returns the same value.
type constantInt struct {
	n int64
}

func (c constantInt) Int63() int64 { return c.n }
