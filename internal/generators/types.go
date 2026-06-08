package generators

import (
	"sort"

	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
)

type MonitorSpec struct {
	Name            string   `json:"name"`
	ID              *int     `json:"id"`
	Description     string   `json:"description"`
	Disabled        bool     `json:"disabled"`
	X               int      `json:"x"`
	AvailableModes  []string `json:"availableModes"`
	Mirror          string   `json:"mirrorOf"`
	CurrentFormat   string   `json:"currentFormat"`
	DpmsStatus      bool     `json:"dpmsStatus"`
	ActivelyTearing bool     `json:"activelyTearing"`
	DirectScanoutTo string   `json:"directScanoutTo"`
	Solitary        string   `json:"solitary"`
}

func NewMonitorSpec(monitor *hypr.MonitorSpec) *MonitorSpec {
	return &MonitorSpec{
		Name:            monitor.Name,
		Description:     monitor.Description,
		ID:              monitor.ID,
		Disabled:        monitor.Disabled,
		X:               monitor.X,
		AvailableModes:  monitor.AvailableModes,
		Mirror:          monitor.Mirror,
		CurrentFormat:   monitor.CurrentFormat,
		DpmsStatus:      monitor.DpmsStatus,
		ActivelyTearing: monitor.ActivelyTearing,
		DirectScanoutTo: monitor.DirectScanoutTo,
		Solitary:        monitor.Solitary,
	}
}

// roleByPosition assigns monitor roles (left/center/right) from the physical
// horizontal layout, so profiles don't need to hand-tag monitors. Connected
// monitors are sorted left-to-right by their X coordinate and mapped by count:
//
//	1 monitor   -> center
//	2 monitors  -> left, center        (leftmost = left, the other = center/main)
//	3+ monitors -> left, center, right (leftmost = left, rightmost = right, middle = center)
//
// Exposed to profile templates as the `roleByPosition` function, e.g.:
//
//	hl_roles = {
//	{{- range $tag, $mon := roleByPosition .Monitors }}
//	  {{ $tag }} = [[{{ $mon.Name }}]],
//	{{- end }}
//	}
//
// NOTE (future adjustment): this assumes the main monitor is the right-hand one
// on a two-monitor setup. If a machine's main display is the LEFT monitor, this
// rule can't tell from position alone — swap left/center here, or add a
// per-profile override hint (e.g. a static_template_value) and honour it below.
// Disabled monitors are ignored.
func roleByPosition(monitors []*MonitorSpec) map[string]*MonitorSpec {
	enabled := make([]*MonitorSpec, 0, len(monitors))
	for _, m := range monitors {
		if m != nil && !m.Disabled {
			enabled = append(enabled, m)
		}
	}
	sort.SliceStable(enabled, func(i, j int) bool { return enabled[i].X < enabled[j].X })

	roles := make(map[string]*MonitorSpec)
	switch len(enabled) {
	case 0:
	case 1:
		roles["center"] = enabled[0]
	case 2:
		roles["left"] = enabled[0]
		roles["center"] = enabled[1]
	default:
		roles["left"] = enabled[0]
		roles["center"] = enabled[len(enabled)/2]
		roles["right"] = enabled[len(enabled)-1]
	}
	return roles
}
