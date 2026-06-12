package profilemaker_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
	"github.com/fiffeek/hyprdynamicmonitors/internal/profilemaker"
	"github.com/fiffeek/hyprdynamicmonitors/internal/testutils"
	"github.com/fiffeek/hyprdynamicmonitors/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// A lua profile frozen from the current setup must carry the hl_roles block as
// unexpanded template directives: the daemon re-evaluates roleByPosition on
// every render, so the roles follow the live monitor layout instead of being
// baked in at freeze time.
func TestService_FreezeGivenAs_LuaRolesBlock(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.toml")
	require.NoError(t, os.WriteFile(configPath, []byte(""), 0o600))

	cfg := testutils.NewTestConfig(t).
		WithConfigPath(configPath).
		WithDestination(filepath.Join(tempDir, "monitors.lua")).
		Get()

	monitors := []*hypr.MonitorSpec{
		{
			ID: utils.IntPtr(0), Name: "DP-1", Description: "Main Display",
			Width: 3840, Height: 2160, RefreshRate: 60.0, X: 1920, Scale: 1.0,
		},
		{
			ID: utils.IntPtr(1), Name: "HDMI-A-1", Description: "",
			Width: 1920, Height: 1080, RefreshRate: 60.0, X: 0, Scale: 1.0,
		},
	}
	for _, monitor := range monitors {
		require.NoError(t, monitor.Validate(), "monitor spec should be correct")
	}

	profileFile := filepath.Join(tempDir, "frozen.lua.tmpl")
	service := profilemaker.NewService(cfg, nil)
	require.NoError(t, service.FreezeGivenAs("frozen", profileFile, monitors))

	content, err := os.ReadFile(profileFile)
	require.NoError(t, err)
	assert.Contains(t, string(content), "roleByPosition .Monitors",
		"roles block should survive freeze as a directive, not be expanded")

	// The frozen file must parse and render under a roleByPosition-style
	// funcMap (the real one lives in the generators package), emitting desc:
	// matches with a connector-name fallback for monitors without one.
	tmpl, err := template.New("frozen").Funcs(template.FuncMap{
		"roleByPosition": func(mons []*hypr.MonitorSpec) map[string]*hypr.MonitorSpec {
			return map[string]*hypr.MonitorSpec{"center": mons[0], "left": mons[1]}
		},
	}).Parse(string(content))
	require.NoError(t, err)

	var rendered bytes.Buffer
	require.NoError(t, tmpl.Execute(&rendered, map[string]any{"Monitors": monitors}))
	assert.Contains(t, rendered.String(), "hl_roles = {")
	assert.Contains(t, rendered.String(), "center = [[desc:Main Display]],")
	assert.Contains(t, rendered.String(), "left = [[HDMI-A-1]],")
}
