package main

import (
	"golang.org/x/sys/windows"
)

var sid *windows.SID

// IsAdminElevated is a function that checks if the current process token belongs to the administrator's group and if it's elevated.
// It first tries to allocate and initialize a security identifier (SID) for the administrators group. If this fails, it returns an error.
// If it succeeds, it defers a call to free the SID, ensuring that the SID is released when the function exits.
// It then gets the current process token and checks if it's a member of the administrators group. If this fails, it returns an error.
// Finally, it checks if the token is elevated and returns this status along with the membership status and any error that occurred.
// Note: Even if a token belongs to the administrators group, it is not necessarily elevated.
// For more information about process elevation, see: https://github.com/mozey/run-as-admin
func IsAdminElevated() (bool, bool, error) {
	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		//log.Fatalf("SID Error: %s", err)
		return false, false, err
	}
	defer windows.FreeSid(sid)

	token := windows.GetCurrentProcessToken()
	member, err := token.IsMember(sid)
	if err != nil {
		return false, token.IsElevated(), err
	}

	// Also note that an admin is _not_ necessarily considered
	// elevated.
	// For elevation see https://github.com/mozey/run-as-admin
	return member, token.IsElevated(), nil
}
