//go:build !windows
// +build !windows

// check privileges on non-Windows systems
package privcheck

import "os"

func HasAdmin() bool {
	return os.Geteuid() == 0
}
