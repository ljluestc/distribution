// Package bloblink provides a middleware that links manifests to blobs
package bloblink

import (
	"context"
	"io"
	"net/http"
	"path"
	"strings"

	dcontext "github.com/distribution/distribution/v3/internal/dcontext"
	storagedriver "github.com/distribution/distribution/v3/registry/storage/driver"
)

// blobLinkPathSuffix is used to distinguish blob link files from actual blob data files.
const blobLinkPathSuffix = ".bloblink"

// manifestBlobLinkPath formats a manifest path as a blob link path
func manifestBlobLinkPath(blobPath string) string {
	return blobPath + blobLinkPathSuffix
}

// isBlobLinkPath checks if a path is a blob link path
func isBlobLinkPath(path string) bool {
	return strings.HasSuffix(path, blobLinkPathSuffix)
}

// getTargetFromLink gets the target path from a blob link
func getTargetFromLink(ctx context.Context, driver storagedriver.StorageDriver, linkPath string) (string, error) {
	targetBytes, err := driver.GetContent(ctx, linkPath)
	if err != nil {
		return "", err
	}

	return string(targetBytes), nil
}

// Middleware implements a storagedriver middleware that redirects blob requests to manifests
// when appropriate, allowing manifests to be accessed via the blob API.
type Middleware struct {
	storagedriver.StorageDriver
}

// BlobLinkDriverFactory creates a middleware that links manifests to blobs.
func BlobLinkDriverFactory(driver storagedriver.StorageDriver) storagedriver.StorageDriver {
	return &Middleware{
		StorageDriver: driver,
	}
}

// GetContent retrieves content from a path. If the path is not found and it's a blob path,
// it will check if there's a corresponding manifest.
func (m *Middleware) GetContent(ctx context.Context, path string) ([]byte, error) {
	content, err := m.StorageDriver.GetContent(ctx, path)
	if err != nil {
		if _, ok := err.(storagedriver.PathNotFoundError); ok {
			// Check if there's a blob link for this path
			linkPath := manifestBlobLinkPath(path)
			targetPath, linkErr := getTargetFromLink(ctx, m.StorageDriver, linkPath)
			if linkErr == nil {
				dcontext.GetLogger(ctx).Debugf("bloblink: redirecting blob request from %s to manifest at %s", path, targetPath)
				return m.StorageDriver.GetContent(ctx, targetPath)
			}
		}
		return nil, err
	}
	return content, nil
}

// Reader retrieves a reader for the content at the given path. If the path is not found
// and it's a blob path, it will check if there's a corresponding manifest.
func (m *Middleware) Reader(ctx context.Context, path string, offset int64) (io.ReadCloser, error) {
	reader, err := m.StorageDriver.Reader(ctx, path, offset)
	if err != nil {
		if _, ok := err.(storagedriver.PathNotFoundError); ok {
			// Check if there's a blob link for this path
			linkPath := manifestBlobLinkPath(path)
			targetPath, linkErr := getTargetFromLink(ctx, m.StorageDriver, linkPath)
			if linkErr == nil {
				dcontext.GetLogger(ctx).Debugf("bloblink: redirecting blob reader request from %s to manifest at %s", path, targetPath)
				return m.StorageDriver.Reader(ctx, targetPath, offset)
			}
		}
		return nil, err
	}
	return reader, nil
}

// Stat returns information about the provided path. If the path is not found
// and it's a blob path, it will check if there's a corresponding manifest.
func (m *Middleware) Stat(ctx context.Context, path string) (storagedriver.FileInfo, error) {
	fi, err := m.StorageDriver.Stat(ctx, path)
	if err != nil {
		if _, ok := err.(storagedriver.PathNotFoundError); ok {
			// Check if there's a blob link for this path
			linkPath := manifestBlobLinkPath(path)
			targetPath, linkErr := getTargetFromLink(ctx, m.StorageDriver, linkPath)
			if linkErr == nil {
				dcontext.GetLogger(ctx).Debugf("bloblink: redirecting blob stat request from %s to manifest at %s", path, targetPath)
				return m.StorageDriver.Stat(ctx, targetPath)
			}
		}
		return nil, err
	}
	return fi, nil
}
