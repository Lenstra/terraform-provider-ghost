package client

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	for _, env := range os.Environ() {
		name := strings.SplitN(env, "=", 2)[0]
		if strings.HasPrefix(name, "GHOST_") {
			os.Unsetenv(name)
		}
	}

	code := m.Run()
	os.Exit(code)
}

func TestConfig(t *testing.T) {
	conf := &Config{}
	_, err := NewClient(conf)
	require.ErrorContains(t, err, "address has not been set")
}
