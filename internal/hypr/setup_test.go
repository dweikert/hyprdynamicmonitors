package hypr_test

import (
	"flag"
	"os"
	"testing"

	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
)

var debug = flag.Bool("debug", false, "enable debug logs")

func TestMain(m *testing.M) {
	// never read the real EDIDs of the host running the tests
	hypr.SetDRMSysfsDir(os.TempDir() + "/hdm-no-such-drm-dir")
	os.Exit(m.Run())
}
