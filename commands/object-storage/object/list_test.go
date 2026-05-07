package object

import (
	"testing"
	"time"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/stretchr/testify/assert"
)

func ptr[T any](v T) *T { return &v }

func TestConvertObjects(t *testing.T) {
	fixedTime := time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC)
	storageClass := objectstorage.OBJECTSTORAGECLASS_STANDARD

	tests := []struct {
		name     string
		input    []objectstorage.Object
		expected []listObjectInfo
	}{
		{
			name:     "empty slice",
			input:    []objectstorage.Object{},
			expected: []listObjectInfo{},
		},
		{
			name: "all fields populated",
			input: []objectstorage.Object{
				{
					Key:          ptr("photos/cat.jpg"),
					Size:         ptr(int32(1024)),
					LastModified: &objectstorage.IonosTime{Time: fixedTime},
					StorageClass: &storageClass,
					ETag:         ptr("\"abc123\""),
				},
			},
			expected: []listObjectInfo{
				{
					Key:          "photos/cat.jpg",
					Size:         "1.0 KiB",
					LastModified: fixedTime,
					StorageClass: "STANDARD",
					ETag:         "\"abc123\"",
				},
			},
		},
		{
			name: "nil LastModified yields zero time",
			input: []objectstorage.Object{
				{
					Key:          ptr("file.txt"),
					Size:         ptr(int32(512)),
					LastModified: nil,
					StorageClass: &storageClass,
					ETag:         ptr("\"def456\""),
				},
			},
			expected: []listObjectInfo{
				{
					Key:          "file.txt",
					Size:         "512 B",
					LastModified: time.Time{},
					StorageClass: "STANDARD",
					ETag:         "\"def456\"",
				},
			},
		},
		{
			name: "nil StorageClass yields empty string",
			input: []objectstorage.Object{
				{
					Key:          ptr("data.bin"),
					Size:         ptr(int32(2048)),
					LastModified: &objectstorage.IonosTime{Time: fixedTime},
					StorageClass: nil,
					ETag:         ptr("\"ghi789\""),
				},
			},
			expected: []listObjectInfo{
				{
					Key:          "data.bin",
					Size:         "2.0 KiB",
					LastModified: fixedTime,
					StorageClass: "",
					ETag:         "\"ghi789\"",
				},
			},
		},
		{
			name: "nil Size returns 0 B",
			input: []objectstorage.Object{
				{
					Key:          ptr("empty.txt"),
					Size:         nil,
					LastModified: &objectstorage.IonosTime{Time: fixedTime},
					StorageClass: &storageClass,
					ETag:         ptr("\"jkl012\""),
				},
			},
			expected: []listObjectInfo{
				{
					Key:          "empty.txt",
					Size:         "0 B",
					LastModified: fixedTime,
					StorageClass: "STANDARD",
					ETag:         "\"jkl012\"",
				},
			},
		},
		{
			name: "multiple objects converted correctly",
			input: []objectstorage.Object{
				{
					Key:          ptr("small.txt"),
					Size:         ptr(int32(100)),
					LastModified: &objectstorage.IonosTime{Time: fixedTime},
					StorageClass: &storageClass,
					ETag:         ptr("\"aaa\""),
				},
				{
					Key:          ptr("medium.bin"),
					Size:         ptr(int32(1048576)), // 1 MiB
					LastModified: &objectstorage.IonosTime{Time: fixedTime.Add(time.Hour)},
					StorageClass: &storageClass,
					ETag:         ptr("\"bbb\""),
				},
				{
					Key:          ptr("large.iso"),
					Size:         ptr(int32(1073741824)), // 1 GiB (overflows int32, but tests the path)
					LastModified: nil,
					StorageClass: nil,
					ETag:         ptr("\"ccc\""),
				},
			},
			expected: []listObjectInfo{
				{
					Key:          "small.txt",
					Size:         "100 B",
					LastModified: fixedTime,
					StorageClass: "STANDARD",
					ETag:         "\"aaa\"",
				},
				{
					Key:          "medium.bin",
					Size:         "1.0 MiB",
					LastModified: fixedTime.Add(time.Hour),
					StorageClass: "STANDARD",
					ETag:         "\"bbb\"",
				},
				{
					Key:          "large.iso",
					Size:         "1.0 GiB",
					LastModified: time.Time{},
					StorageClass: "",
					ETag:         "\"ccc\"",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertObjects(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
