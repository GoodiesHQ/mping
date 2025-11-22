//go:build windows
// +build windows

// check privileges on Windows systems
package privcheck

import "golang.org/x/sys/windows"

// HasAdmin reports whether the current process is running elevated
// (i.e. as Administrator) on Windows.
func HasAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	admin, err := windows.Token(0).IsMember(sid)
	return err == nil && admin
}
