package client

import (
	"regexp"
	"strings"
)

// RedactSensitiveInfo removes sensitive information from error messages
// such as S3 bucket names, repository paths, IP addresses, etc.
func RedactSensitiveInfo(message string) string {
	if message == "" {
		return ""
	}

	// Redact S3 bucket names and paths (e.g., https://s3.amazonaws.com/bucket-name/...)
	s3URLRe := regexp.MustCompile(`https?://s3\.amazonaws\.com/[^\s/]+/[^\s"]*`)
	message = s3URLRe.ReplaceAllString(message, "[REDACTED_S3_URL]")

	// Redact S3 bucket subdomains (e.g., https://bucket-name.s3.amazonaws.com/...)
	s3SubdomainRe := regexp.MustCompile(`https?://[^\s/]+\.s3\.amazonaws\.com/[^\s"]*`)
	message = s3SubdomainRe.ReplaceAllString(message, "[REDACTED_S3_SUBDOMAIN_URL]")

	// Redact repository paths (e.g., /docker/registry/v2/repositories/org/repo/...)
	repoPathRe := regexp.MustCompile(`/docker/registry/v2/repositories/[^/\s"]+/[^/\s"]+(/[^\s"]*)?`)
	message = repoPathRe.ReplaceAllString(message, "/docker/registry/v2/repositories/[REDACTED_ORG]/[REDACTED_REPO]$1")

	// Redact upload IDs (UUIDs)
	uuidRe := regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	message = uuidRe.ReplaceAllString(message, "[REDACTED_UUID]")

	// Redact IP addresses with ports (e.g., 54.231.9.56:443)
	ipPortRe := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}:\d+`)
	message = ipPortRe.ReplaceAllString(message, "[REDACTED_IP]")

	// Redact IP addresses without ports (e.g., 54.231.9.56)
	ipRe := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	message = ipRe.ReplaceAllString(message, "[REDACTED_IP]")

	return message
}
