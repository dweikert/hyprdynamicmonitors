package hypr

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadEDIDHash(t *testing.T) {
	dir := t.TempDir()
	conn := filepath.Join(dir, "card1-DP-6")
	require.NoError(t, os.MkdirAll(conn, 0o750))
	edid := []byte{0, 255, 255, 255, 255, 255, 255, 0, 1, 2, 3}
	require.NoError(t, os.WriteFile(filepath.Join(conn, "edid"), edid, 0o600))

	sum := sha256.Sum256(edid)
	assert.Equal(t, hex.EncodeToString(sum[:]), readEDIDHash(dir, "DP-6"))
	assert.Empty(t, readEDIDHash(dir, "HDMI-A-1"))

	// connectors without a monitor expose an empty edid file
	empty := filepath.Join(dir, "card1-HDMI-A-1")
	require.NoError(t, os.MkdirAll(empty, 0o750))
	require.NoError(t, os.WriteFile(filepath.Join(empty, "edid"), nil, 0o600))
	assert.Empty(t, readEDIDHash(dir, "HDMI-A-1"))
}
