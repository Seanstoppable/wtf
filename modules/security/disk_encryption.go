package security

import (
	"os/exec"
	"runtime"

	"github.com/wtfutil/wtf/wtf"
)

// DiskEncryptionState checks for in built disk encryption of primary drive
func DiskEncryptionState() string {
	switch runtime.GOOS {
	case "darwin":
		return diskEncryptionStateMacOS()
	default:
		return "Unknown"
	}
}

func diskEncryptionStateMacOS() string {

	cmd := exec.Command("fdesetup", "status")
	str := wtf.ExecuteCommand(cmd)

	if str == "FileVault is On.\n" {
		return "Enabled"
	} else if str == "FileVault is Off.\n" {
		return "Disabled"
	}
	return "Unknown/Module is broken"
}
