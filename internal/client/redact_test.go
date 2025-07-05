package client
package client

import (
	"testing"
)

func TestRedactSensitiveInfo(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "original issue example",
			input:    `DEBU[0264] Unexpected response from server: "{\"errors\":[{\"code\":\"UNKNOWN\",\"message\":\"unknown error\",\"detail\":\"Put https://s3.amazonaws.com/docker-images-prod/registry-v2/docker/registry/v2/repositories/lyst/sortinghat/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/startedat: dial tcp 54.231.9.56:443: connection timed out\"}]}\n"`,
			expected: `DEBU[0264] Unexpected response from server: "{\"errors\":[{\"code\":\"UNKNOWN\",\"message\":\"unknown error\",\"detail\":\"Put [REDACTED_S3_URL]: dial tcp [REDACTED_IP]: connection timed out\"}]}\n"`,
		},
		{
			name:     "s3 bucket URL",
			input:    "https://s3.amazonaws.com/bucket-name/path/to/object",
			expected: "[REDACTED_S3_URL]",
		},
		{
			name:     "s3 bucket subdomain",
			input:    "https://bucket-name.s3.amazonaws.com/path/to/object",
			expected: "[REDACTED_S3_SUBDOMAIN_URL]",
		},
		{
			name:     "repository path",
			input:    "/docker/registry/v2/repositories/org/repo/_uploads/123456",
			expected: "/docker/registry/v2/repositories/[REDACTED_ORG]/[REDACTED_REPO]/_uploads/123456",
		},
		{
			name:     "uuid",
			input:    "fafc5d11-0c70-4887-9c1b-24cb01d04f02",
			expected: "[REDACTED_UUID]",
		},
		{
			name:     "ip address with port",
			input:    "dial tcp 54.231.9.56:443: connection timed out",
			expected: "dial tcp [REDACTED_IP]: connection timed out",
		},
		{
			name:     "ip address without port",
			input:    "connecting to 54.231.9.56 failed",
			expected: "connecting to [REDACTED_IP] failed",
		},
		{
			name:     "multiple sensitive information",
			input:    "Error accessing https://s3.amazonaws.com/bucket-name/registry-v2/docker/registry/v2/repositories/org/repo/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/data: dial tcp 54.231.9.56:443: connection timed out",
			expected: "Error accessing [REDACTED_S3_URL]: dial tcp [REDACTED_IP]: connection timed out",
		},
		{
			name:     "regular text",
			input:    "This is a regular message with no sensitive info",
			expected: "This is a regular message with no sensitive info",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := RedactSensitiveInfo(tc.input)
			if result != tc.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tc.expected, result)
			}
		})
	}
}
import (
	"testing"
)

func TestRedactSensitiveInfo(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "s3 bucket URL",
			input:    `{"errors":[{"code":"UNKNOWN","message":"unknown error","detail":"Put https://s3.amazonaws.com/docker-images-prod/registry-v2/docker/registry/v2/repositories/lyst/sortinghat/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/startedat: dial tcp 54.231.9.56:443: connection timed out"}]}`,
			expected: `{"errors":[{"code":"UNKNOWN","message":"unknown error","detail":"Put https://s3.amazonaws.com/[REDACTED]/registry-v2/docker/registry/v2/repositories/[REDACTED]/[REDACTED]/_uploads/[REDACTED-UUID]/startedat: dial tcp [REDACTED-IP]:443: connection timed out"}]}`,
		},
		{
			name:     "s3 bucket subdomain",
			input:    "https://my-bucket.s3.amazonaws.com/path/to/object.jpg",
			expected: "https://[REDACTED].s3.amazonaws.com/path/to/object.jpg",
		},
		{
			name:     "repository path",
			input:    "/docker/registry/v2/repositories/myorg/myrepo/_uploads/123456",
			expected: "/docker/registry/v2/repositories/[REDACTED]/[REDACTED]/_uploads/123456",
		},
		{
			name:     "uuid",
			input:    "fafc5d11-0c70-4887-9c1b-24cb01d04f02",
			expected: "[REDACTED-UUID]",
		},
		{
			name:     "ip address",
			input:    "dial tcp 54.231.9.56:443: connection timed out",
			expected: "dial tcp [REDACTED-IP]:443: connection timed out",
		},
		{
			name:     "regular text",
			input:    "This is a regular message with no sensitive info",
			expected: "This is a regular message with no sensitive info",
		},
		{
			name:     "original issue example",
			input:    `DEBU[0264] Unexpected response from server: "{\"errors\":[{\"code\":\"UNKNOWN\",\"message\":\"unknown error\",\"detail\":\"Put https://s3.amazonaws.com/docker-images-prod/registry-v2/docker/registry/v2/repositories/lyst/sortinghat/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/startedat: dial tcp 54.231.9.56:443: connection timed out\"}]}\n"`,
			expected: `DEBU[0264] Unexpected response from server: "{\"errors\":[{\"code\":\"UNKNOWN\",\"message\":\"unknown error\",\"detail\":\"Put https://s3.amazonaws.com/[REDACTED]/registry-v2/docker/registry/v2/repositories/[REDACTED]/[REDACTED]/_uploads/[REDACTED-UUID]/startedat: dial tcp [REDACTED-IP]:443: connection timed out\"}]}\n"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := redactSensitiveInfo(tc.input)
			if result != tc.expected {
				t.Errorf("Expected: %s\nGot: %s", tc.expected, result)
			}
		})
	}
}

// TestSanitizeErrorMessage tests the exported SanitizeErrorMessage function
func TestSanitizeErrorMessage(t *testing.T) {
	input := `{"errors":[{"code":"UNKNOWN","message":"unknown error","detail":"Put https://s3.amazonaws.com/docker-images-prod/registry-v2/docker/registry/v2/repositories/lyst/sortinghat/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/startedat: dial tcp 54.231.9.56:443: connection timed out"}]}`
	expected := `{"errors":[{"code":"UNKNOWN","message":"unknown error","detail":"Put https://s3.amazonaws.com/[REDACTED]/registry-v2/docker/registry/v2/repositories/[REDACTED]/[REDACTED]/_uploads/[REDACTED-UUID]/startedat: dial tcp [REDACTED-IP]:443: connection timed out"}]}`
	
	result := SanitizeErrorMessage(input)
	if result != expected {
		t.Errorf("SanitizeErrorMessage failed\nExpected: %s\nGot: %s", expected, result)
	}
}
