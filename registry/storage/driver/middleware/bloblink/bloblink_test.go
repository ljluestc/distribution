package bloblink

import (
	"context"
	"io"
	"testing"

	"github.com/distribution/distribution/v3/registry/storage/driver/inmemory"
	"github.com/stretchr/testify/require"
)

func TestBlobLinkMiddleware(t *testing.T) {
	ctx := context.Background()
	baseDriver := inmemory.New()
	driver := BlobLinkDriverFactory(baseDriver)

	// Create a manifest
	manifestPath := "/docker/registry/v2/repositories/test/manifests/revisions/sha256:abcdef"
	manifestContent := []byte("manifest content")
	err := baseDriver.PutContent(ctx, manifestPath, manifestContent)
	require.NoError(t, err)

	// Create a blob link pointing to the manifest
	blobPath := "/docker/registry/v2/blobs/sha256/ab/sha256:abcdef"
	err = baseDriver.PutContent(ctx, manifestBlobLinkPath(blobPath), []byte(manifestPath))
	require.NoError(t, err)

	// Test GetContent
	content, err := driver.GetContent(ctx, blobPath)
	require.NoError(t, err)
	require.Equal(t, manifestContent, content)

	// Test Reader
	reader, err := driver.Reader(ctx, blobPath, 0)
	require.NoError(t, err)
	defer reader.Close()
	
	readContent, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.Equal(t, manifestContent, readContent)

	// Test Stat
	fileInfo, err := driver.Stat(ctx, blobPath)
	require.NoError(t, err)
	require.Equal(t, int64(len(manifestContent)), fileInfo.Size())

	// Test with non-existent path
	_, err = driver.GetContent(ctx, "/non-existent")
	require.Error(t, err)
}
