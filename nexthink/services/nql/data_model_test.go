package nql

import (
	"testing"
)

func TestTableConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Devices", TableDevices, "devices"},
		{"Users", TableUsers, "users"},
		{"Applications", TableApplications, "applications"},
		{"Binaries", TableBinaries, "binaries"},
		{"ExecutionCrashes", TableExecutionCrashes, "execution.crashes"},
		{"WebErrors", TableWebErrors, "web.errors"},
		{"DEXScores", TableDexScores, "dex.scores"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestFieldConstants_Device(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"DeviceName", FieldDeviceName, "device.name"},
		{"OSName", FieldOSName, "operating_system.name"},
		{"OSPlatform", FieldOSPlatform, "operating_system.platform"},
		{"HardwareType", FieldHardwareType, "hardware.type"},
		{"HardwareManufacturer", FieldHardwareManufacturer, "hardware.manufacturer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestFieldConstants_User(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"UserName", FieldUserName, "user.name"},
		{"UserType", FieldUserType, "user.type"},
		{"UserSID", FieldUserSID, "user.sid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestPlatformConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Windows", PlatformWindows, "Windows"},
		{"macOS", PlatformMacOS, "macOS"},
		{"Linux", PlatformLinux, "Linux"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestHardwareTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Laptop", HardwareTypeLaptop, "laptop"},
		{"Desktop", HardwareTypeDesktop, "desktop"},
		{"Server", HardwareTypeServer, "server"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestUserTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"LocalUser", UserTypeLocalUser, "LOCAL_USER"},
		{"LocalAdmin", UserTypeLocalAdmin, "LOCAL_ADMIN"},
		{"DomainUser", UserTypeDomainUser, "DOMAIN_USER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestExperienceLevelConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"Good", ExperienceLevelGood, "good"},
		{"Frustrating", ExperienceLevelFrustrating, "frustrating"},
		{"Tolerable", ExperienceLevelTolerable, "tolerable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Expected %s, got: %s", tt.expected, tt.constant)
			}
		})
	}
}

func TestNamespaceConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		contains string
	}{
		{"Device", NamespaceDevice, "device"},
		{"DevicePerformance", NamespaceDevicePerformance, "device_performance"},
		{"User", NamespaceUser, "user"},
		{"Application", NamespaceApplication, "application"},
		{"Binary", NamespaceBinary, "binary"},
		{"Web", NamespaceWeb, "web"},
		{"DEX", NamespaceDex, "dex"},
		{"Execution", NamespaceExecution, "execution"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.contains {
				t.Errorf("Expected %s, got: %s", tt.contains, tt.constant)
			}
		})
	}
}

func TestFieldConstants_Comprehensive(t *testing.T) {
	// Just verify that all exported field constants are non-empty
	fieldConstants := []struct {
		name  string
		value string
	}{
		{"FieldDeviceName", FieldDeviceName},
		{"FieldOSName", FieldOSName},
		{"FieldOSPlatform", FieldOSPlatform},
		{"FieldHardwareType", FieldHardwareType},
		{"FieldUserName", FieldUserName},
		{"FieldUserType", FieldUserType},
		{"FieldApplicationName", FieldApplicationName},
		{"FieldBinaryName", FieldBinaryName},
	}

	for _, field := range fieldConstants {
		t.Run(field.name, func(t *testing.T) {
			if field.value == "" {
				t.Errorf("Field constant %s is empty", field.name)
			}
		})
	}
}

func TestValueConstants_Comprehensive(t *testing.T) {
	// Verify value constants are non-empty
	valueConstants := []struct {
		name  string
		value string
	}{
		{"PlatformWindows", PlatformWindows},
		{"PlatformMacOS", PlatformMacOS},
		{"PlatformLinux", PlatformLinux},
		{"HardwareTypeLaptop", HardwareTypeLaptop},
		{"HardwareTypeDesktop", HardwareTypeDesktop},
		{"ExperienceLevelGood", ExperienceLevelGood},
		{"ExperienceLevelFrustrating", ExperienceLevelFrustrating},
	}

	for _, val := range valueConstants {
		t.Run(val.name, func(t *testing.T) {
			if val.value == "" {
				t.Errorf("Value constant %s is empty", val.name)
			}
		})
	}
}
