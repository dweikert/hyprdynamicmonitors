package generators

import "testing"

func mon(name string, x int, disabled bool) *MonitorSpec {
	return &MonitorSpec{Name: name, X: x, Disabled: disabled}
}

func TestRoleByPosition(t *testing.T) {
	tests := []struct {
		name     string
		monitors []*MonitorSpec
		want     map[string]string // role -> monitor name
	}{
		{
			name:     "none",
			monitors: nil,
			want:     map[string]string{},
		},
		{
			name:     "single is center",
			monitors: []*MonitorSpec{mon("DP-1", 0, false)},
			want:     map[string]string{"center": "DP-1"},
		},
		{
			name:     "two: leftmost left, other center",
			monitors: []*MonitorSpec{mon("HDMI-A-1", 2560, false), mon("HDMI-A-2", 0, false)},
			want:     map[string]string{"left": "HDMI-A-2", "center": "HDMI-A-1"},
		},
		{
			name: "three by x: left center right",
			monitors: []*MonitorSpec{
				mon("DP-3", 5120, false), mon("DP-1", 0, false), mon("DP-2", 2560, false),
			},
			want: map[string]string{"left": "DP-1", "center": "DP-2", "right": "DP-3"},
		},
		{
			name: "disabled monitors ignored",
			monitors: []*MonitorSpec{
				mon("eDP-1", -1, true), mon("HDMI-A-2", 0, false), mon("HDMI-A-1", 2560, false),
			},
			want: map[string]string{"left": "HDMI-A-2", "center": "HDMI-A-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roleByPosition(tt.monitors)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d roles, want %d: %v", len(got), len(tt.want), got)
			}
			for role, name := range tt.want {
				m, ok := got[role]
				if !ok {
					t.Errorf("missing role %q", role)
					continue
				}
				if m.Name != name {
					t.Errorf("role %q = %q, want %q", role, m.Name, name)
				}
			}
		})
	}
}
