package hypr

// SetDRMSysfsDir lets the external test package redirect the EDID lookup.
func SetDRMSysfsDir(dir string) { drmSysfsDir = dir }
