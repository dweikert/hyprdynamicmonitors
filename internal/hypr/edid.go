package hypr

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// drmSysfsDir is a var so tests can point it at a fixture.
var drmSysfsDir = "/sys/class/drm"

// ReadEDIDHash returns the hex sha256 of the EDID exposed by the kernel for the given
// connector name (e.g. DP-6), or an empty string if it cant be read.
func ReadEDIDHash(connector string) string {
	return readEDIDHash(drmSysfsDir, connector)
}

func readEDIDHash(drmDir, connector string) string {
	// connectors are exposed as card<N>-<connector>
	paths, err := filepath.Glob(filepath.Join(drmDir, "card*-"+connector, "edid"))
	if err != nil || len(paths) == 0 {
		return ""
	}
	for _, path := range paths {
		data, err := os.ReadFile(path) // #nosec G304 -- sysfs path built from a glob
		if err != nil || len(data) == 0 {
			logrus.WithError(err).WithField("path", path).Debug("Cant read EDID")
			continue
		}
		sum := sha256.Sum256(data)
		return hex.EncodeToString(sum[:])
	}
	return ""
}
