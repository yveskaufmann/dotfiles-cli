package osutil

import (
	"bufio"
	"os"
	"runtime"
	"strings"
)

type OSType string

const (
	OSLinux     OSType = "linux"
	OSMac       OSType = "darwin"
	OSUbuntu    OSType = "ubuntu"
	OSUbuntuWSL OSType = "ubuntu-wsl"
	OSOther     OSType = "other"
)

func Is(osType OSType) bool {
	switch osType {
	case OSLinux:
		return IsLinux()
	case OSMac:
		return IsMac()
	case OSUbuntu:
		return IsUbuntu() && !isWSL()
	case OSUbuntuWSL:
		return IsUbuntu() && isWSL()
	default:
		return false
	}
}

// Check if the current os is linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// Check if the current os is macOS
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// Check if the current linux distro is ubuntu
func IsUbuntu() bool {
	if !IsLinux() {
		return false
	}

	osInfo, err := readOSReleaseFile()
	if err != nil {
		panic(err)
	}

	if id, exists := osInfo["ID"]; exists {
		return id == "ubuntu"
	}

	if name, exists := osInfo["Name"]; exists {
		return name == "ubuntu"
	}

	return false
}

// Check if the current linux distro is running under WSL
func isWSL() bool {
	if !IsLinux() {
		return false
	}

	// check /proc/version and kernel osrelease for "microsoft" / "wsl"
	if b, err := os.ReadFile("/proc/version"); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "microsoft") || strings.Contains(s, "wsl") {
			return true
		}
	}
	if b, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil {
		s := strings.ToLower(string(b))
		if strings.Contains(s, "microsoft") || strings.Contains(s, "wsl") {
			return true
		}
	}

	return false
}

func readOSReleaseFile() (map[string]string, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := make(map[string]string)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := parts[0]
			value := strings.Trim(parts[1], `"`)
			info[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return info, nil
}
