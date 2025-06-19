package xeger_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/travisnewhouse/go-xeger"
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

func Test_WithLimit(t *testing.T) {
	expr := `[a-z]{2,}`
	re := regexp.MustCompile(expr)

	t.Run("without limit", func(t *testing.T) {
		x, err := xeger.NewXeger(expr)
		require.NoError(t, err)
		s1 := x.Generate()
		require.Regexp(t, re, s1)
		// Default limit is 10.
		require.GreaterOrEqual(t, len(s1), 2, "Expect string of at least 2 characters")
		require.LessOrEqual(t, len(s1), 10, "Expect string not longer than 10 characters")
	})

	t.Run("with limit equal to lower range", func(t *testing.T) {
		x, err := xeger.NewXeger(expr, xeger.WithLimit(2))
		require.NoError(t, err)
		s1 := x.Generate()
		require.Regexp(t, re, s1)
		require.Len(t, s1, 2, "Expect string length to be 2 characters due to limit")
	})

	t.Run("with limit greater than lower range, less than default", func(t *testing.T) {
		x, err := xeger.NewXeger(expr, xeger.WithLimit(5))
		require.NoError(t, err)
		s1 := x.Generate()
		require.Regexp(t, re, s1)
		require.GreaterOrEqual(t, len(s1), 2, "Expect string of at least 2 characters")
		require.LessOrEqual(t, len(s1), 5, "Expect string not longer than 5 characters")
	})
}

func Test_WithSource(t *testing.T) {
	rng := constantInt{n: 42}
	expr := `[a-zA-Z0-9]{32}`
	re := regexp.MustCompile(expr)

	x, err := xeger.NewXeger(expr, xeger.WithSource(rng))
	require.NoError(t, err)
	s1 := x.Generate()
	require.Regexp(t, re, s1)

	xDefault, err := xeger.NewXeger(expr)
	require.NoError(t, err)
	s2 := xDefault.GenerateWithSource(rng)
	require.Regexp(t, re, s2)

	require.Equal(t, s1, s2, "Expect strings to be equal when source of random numbers is the repeatable.")
}

// constantInt implements RandSource but always returns the same value.
type constantInt struct {
	n int64
}

func (c constantInt) Int63() int64 { return c.n }
