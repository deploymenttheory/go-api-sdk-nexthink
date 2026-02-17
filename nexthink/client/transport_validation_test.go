package client

import (
	"strings"
	"testing"
)

func TestValidateTransportConfig(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		instance     string
		region       string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "valid configuration",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       RegionUS,
			wantErr:      false,
		},
		{
			name:         "valid EU region",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       RegionEU,
			wantErr:      false,
		},
		{
			name:         "valid PAC region",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       RegionPAC,
			wantErr:      false,
		},
		{
			name:         "valid META region",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       RegionMETA,
			wantErr:      false,
		},
		{
			name:         "empty client ID",
			clientID:     "",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "client ID cannot be empty",
		},
		{
			name:         "empty client secret",
			clientID:     "test-client-id",
			clientSecret: "",
			instance:     "myinstance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "client secret cannot be empty",
		},
		{
			name:         "empty instance",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "instance cannot be empty",
		},
		{
			name:         "empty region",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       "",
			wantErr:      true,
			errContains:  "region cannot be empty",
		},
		{
			name:         "invalid region",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "myinstance",
			region:       "invalid",
			wantErr:      true,
			errContains:  "invalid region",
		},
		{
			name:         "instance with spaces",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "my instance",
			region:       RegionUS,
			wantErr:      true,
			errContains:  "instance name cannot contain spaces",
		},
		{
			name:         "instance too long",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     strings.Repeat("a", 101),
			region:       RegionUS,
			wantErr:      true,
			errContains:  "instance name too long",
		},
		{
			name:         "instance at max length (100 chars)",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     strings.Repeat("a", 100),
			region:       RegionUS,
			wantErr:      false,
		},
		{
			name:         "instance with hyphen (valid)",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "my-instance",
			region:       RegionUS,
			wantErr:      false,
		},
		{
			name:         "instance with underscore (valid)",
			clientID:     "test-client-id",
			clientSecret: "test-secret",
			instance:     "my_instance",
			region:       RegionUS,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransportConfig(tt.clientID, tt.clientSecret, tt.instance, tt.region)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTransportConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateTransportConfig() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestValidateBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid HTTPS URL",
			baseURL: "https://api.nexthink.com",
			wantErr: false,
		},
		{
			name:    "valid HTTP URL",
			baseURL: "http://localhost:8080",
			wantErr: false,
		},
		{
			name:    "valid HTTPS with subdomain",
			baseURL: "https://myinstance.api.us.nexthink.cloud",
			wantErr: false,
		},
		{
			name:        "empty URL",
			baseURL:     "",
			wantErr:     true,
			errContains: "base URL cannot be empty",
		},
		{
			name:        "URL without protocol",
			baseURL:     "api.nexthink.com",
			wantErr:     true,
			errContains: "must start with http:// or https://",
		},
		{
			name:        "URL with trailing slash",
			baseURL:     "https://api.nexthink.com/",
			wantErr:     true,
			errContains: "should not end with a trailing slash",
		},
		{
			name:        "URL with path",
			baseURL:     "https://api.nexthink.com/v1",
			wantErr:     false,
		},
		{
			name:        "URL with port",
			baseURL:     "https://api.nexthink.com:8443",
			wantErr:     false,
		},
		{
			name:        "FTP URL (invalid)",
			baseURL:     "ftp://api.nexthink.com",
			wantErr:     true,
			errContains: "must start with http:// or https://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBaseURL(tt.baseURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateBaseURL() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestValidateTimeout(t *testing.T) {
	tests := []struct {
		name        string
		timeout     int
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid timeout (30 seconds)",
			timeout: 30,
			wantErr: false,
		},
		{
			name:    "valid timeout (60 seconds)",
			timeout: 60,
			wantErr: false,
		},
		{
			name:    "valid timeout (1 second)",
			timeout: 1,
			wantErr: false,
		},
		{
			name:    "valid timeout (max 3600 seconds)",
			timeout: 3600,
			wantErr: false,
		},
		{
			name:        "zero timeout",
			timeout:     0,
			wantErr:     true,
			errContains: "timeout must be greater than 0",
		},
		{
			name:        "negative timeout",
			timeout:     -1,
			wantErr:     true,
			errContains: "timeout must be greater than 0",
		},
		{
			name:        "timeout too large",
			timeout:     3601,
			wantErr:     true,
			errContains: "timeout too large",
		},
		{
			name:        "timeout way too large",
			timeout:     10000,
			wantErr:     true,
			errContains: "timeout too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimeout(tt.timeout)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateTimeout() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestValidateRetryCount(t *testing.T) {
	tests := []struct {
		name        string
		retryCount  int
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid retry count (0)",
			retryCount: 0,
			wantErr:    false,
		},
		{
			name:       "valid retry count (3)",
			retryCount: 3,
			wantErr:    false,
		},
		{
			name:       "valid retry count (5)",
			retryCount: 5,
			wantErr:    false,
		},
		{
			name:       "valid retry count (max 10)",
			retryCount: 10,
			wantErr:    false,
		},
		{
			name:        "negative retry count",
			retryCount:  -1,
			wantErr:     true,
			errContains: "retry count cannot be negative",
		},
		{
			name:        "retry count too large",
			retryCount:  11,
			wantErr:     true,
			errContains: "retry count too large",
		},
		{
			name:        "retry count way too large",
			retryCount:  100,
			wantErr:     true,
			errContains: "retry count too large",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRetryCount(tt.retryCount)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRetryCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateRetryCount() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestValidateProxyURL(t *testing.T) {
	tests := []struct {
		name        string
		proxyURL    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "empty proxy URL (valid - no proxy)",
			proxyURL: "",
			wantErr:  false,
		},
		{
			name:     "valid HTTP proxy",
			proxyURL: "http://proxy.company.com:8080",
			wantErr:  false,
		},
		{
			name:     "valid HTTPS proxy",
			proxyURL: "https://proxy.company.com:8443",
			wantErr:  false,
		},
		{
			name:     "valid SOCKS5 proxy",
			proxyURL: "socks5://127.0.0.1:1080",
			wantErr:  false,
		},
		{
			name:     "valid HTTP proxy without port",
			proxyURL: "http://proxy.company.com",
			wantErr:  false,
		},
		{
			name:     "valid localhost proxy",
			proxyURL: "http://localhost:8888",
			wantErr:  false,
		},
		{
			name:     "valid IP proxy",
			proxyURL: "http://192.168.1.100:3128",
			wantErr:  false,
		},
		{
			name:        "invalid proxy (no protocol)",
			proxyURL:    "proxy.company.com:8080",
			wantErr:     true,
			errContains: "must start with http://, https://, or socks5://",
		},
		{
			name:        "invalid proxy (FTP protocol)",
			proxyURL:    "ftp://proxy.company.com:21",
			wantErr:     true,
			errContains: "must start with http://, https://, or socks5://",
		},
		{
			name:        "invalid proxy (socks4 protocol)",
			proxyURL:    "socks4://127.0.0.1:1080",
			wantErr:     true,
			errContains: "must start with http://, https://, or socks5://",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProxyURL(tt.proxyURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProxyURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("ValidateProxyURL() error = %v, want to contain %q", err, tt.errContains)
			}
		})
	}
}

func TestValidateTransportConfig_AllRegions(t *testing.T) {
	// Test all valid regions systematically
	validRegions := []string{RegionUS, RegionEU, RegionPAC, RegionMETA}

	for _, region := range validRegions {
		t.Run("region_"+region, func(t *testing.T) {
			err := ValidateTransportConfig("test-id", "test-secret", "test-instance", region)
			if err != nil {
				t.Errorf("ValidateTransportConfig() failed for valid region %q: %v", region, err)
			}
		})
	}
}

func TestValidateTransportConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		instance     string
		region       string
		wantErr      bool
		description  string
	}{
		{
			name:         "whitespace in client ID",
			clientID:     "  test-id  ",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       RegionUS,
			wantErr:      false,
			description:  "whitespace is allowed in client ID",
		},
		{
			name:         "special characters in client ID",
			clientID:     "test-id!@#$%",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       RegionUS,
			wantErr:      false,
			description:  "special characters allowed in client ID",
		},
		{
			name:         "numbers in instance",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "instance123",
			region:       RegionUS,
			wantErr:      false,
			description:  "numbers allowed in instance",
		},
		{
			name:         "uppercase region",
			clientID:     "test-id",
			clientSecret: "test-secret",
			instance:     "test-instance",
			region:       "US",
			wantErr:      true,
			description:  "region is case-sensitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTransportConfig(tt.clientID, tt.clientSecret, tt.instance, tt.region)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: ValidateTransportConfig() error = %v, wantErr %v", tt.description, err, tt.wantErr)
			}
		})
	}
}

func TestValidateBaseURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		baseURL     string
		wantErr     bool
		description string
	}{
		{
			name:        "URL with multiple slashes at end",
			baseURL:     "https://api.nexthink.com//",
			wantErr:     true,
			description: "multiple trailing slashes should fail",
		},
		{
			name:        "URL with query parameters",
			baseURL:     "https://api.nexthink.com?version=1",
			wantErr:     false,
			description: "query parameters are allowed",
		},
		{
			name:        "URL with fragment",
			baseURL:     "https://api.nexthink.com#section",
			wantErr:     false,
			description: "fragments are allowed",
		},
		{
			name:        "URL with authentication",
			baseURL:     "https://user:pass@api.nexthink.com",
			wantErr:     false,
			description: "basic auth in URL is allowed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBaseURL(tt.baseURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("%s: ValidateBaseURL() error = %v, wantErr %v", tt.description, err, tt.wantErr)
			}
		})
	}
}

func TestValidateTimeout_BoundaryValues(t *testing.T) {
	// Test boundary values
	tests := []struct {
		timeout int
		wantErr bool
	}{
		{timeout: 0, wantErr: true},      // Just below valid range
		{timeout: 1, wantErr: false},     // Lower boundary
		{timeout: 1800, wantErr: false},  // Middle value
		{timeout: 3600, wantErr: false},  // Upper boundary
		{timeout: 3601, wantErr: true},   // Just above valid range
	}

	for _, tt := range tests {
		t.Run("timeout_"+strings.ReplaceAll(string(rune(tt.timeout)), "-", "neg"), func(t *testing.T) {
			err := ValidateTimeout(tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTimeout(%d) error = %v, wantErr %v", tt.timeout, err, tt.wantErr)
			}
		})
	}
}

func TestValidateRetryCount_BoundaryValues(t *testing.T) {
	// Test boundary values
	tests := []struct {
		retryCount int
		wantErr    bool
	}{
		{retryCount: -1, wantErr: true},    // Just below valid range
		{retryCount: 0, wantErr: false},    // Lower boundary
		{retryCount: 5, wantErr: false},    // Middle value
		{retryCount: 10, wantErr: false},   // Upper boundary
		{retryCount: 11, wantErr: true},    // Just above valid range
	}

	for _, tt := range tests {
		t.Run("retry_count_"+strings.ReplaceAll(string(rune(tt.retryCount)), "-", "neg"), func(t *testing.T) {
			err := ValidateRetryCount(tt.retryCount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRetryCount(%d) error = %v, wantErr %v", tt.retryCount, err, tt.wantErr)
			}
		})
	}
}
