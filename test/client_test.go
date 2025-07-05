package test

import (
	"testing"

	"github.com/distribution/distribution/v3/internal/client"
)

func TestRedactError(t *testing.T) {
	errorMsg := `{"errors":[{"code":"UNKNOWN","message":"unknown error","detail":"Put https://s3.amazonaws.com/docker-images-prod/registry-v2/docker/registry/v2/repositories/lyst/sortinghat/_uploads/fafc5d11-0c70-4887-9c1b-24cb01d04f02/startedat: dial tcp 54.231.9.56:443: connection timed out"}]}`
	redacted := client.RedactSensitiveInfo(errorMsg)
	
	// Verify S3 bucket is redacted
	if contains(redacted, "docker-images-prod") {
		t.Errorf("S3 bucket name 'docker-images-prod' was not redacted: %s", redacted)
	}
	
	// Verify repository name is redacted
	if contains(redacted, "lyst/sortinghat") {
		t.Errorf("Repository path 'lyst/sortinghat' was not redacted: %s", redacted)
	}
	
	// Verify IP address is redacted
	if contains(redacted, "54.231.9.56") {
		t.Errorf("IP address '54.231.9.56' was not redacted: %s", redacted)
	}
	
	// Verify UUID is redacted
	if contains(redacted, "fafc5d11-0c70-4887-9c1b-24cb01d04f02") {
		t.Errorf("UUID 'fafc5d11-0c70-4887-9c1b-24cb01d04f02' was not redacted: %s", redacted)
	}
}

func TestRedactNonJSONError(t *testing.T) {
	errorMsg := "Error: Put https://s3.amazonaws.com/docker-images-prod/test: dial tcp 54.231.9.56:443: connection timed out"
	redacted := client.RedactSensitiveInfo(errorMsg)
	
	// Verify S3 bucket is redacted
	if contains(redacted, "docker-images-prod") {
		t.Errorf("S3 bucket name 'docker-images-prod' was not redacted: %s", redacted)
	}
	
	// Verify IP address is redacted
	if contains(redacted, "54.231.9.56") {
		t.Errorf("IP address '54.231.9.56' was not redacted: %s", redacted)
	}
	
	expected := "Error: Put [REDACTED_S3_URL]: dial tcp [REDACTED_IP]: connection timed out"
	if redacted != expected {
		t.Errorf("Expected redacted message: %s, got: %s", expected, redacted)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
