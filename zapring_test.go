package zapring_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"moul.io/zapring"
)

func TestSync(t *testing.T) {
	ring := zapring.New(10 * 1024 * 1024)
	defer ring.Close()
	require.NoError(t, ring.Sync())
}

func TestSync_wrapped(t *testing.T) {
	ring := zapring.New(10 * 1024 * 1024)
	defer ring.Close()
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	ring.SetNextCore(logger.Core())
	require.NotPanics(t, func() { _ = ring.Sync() })
}
