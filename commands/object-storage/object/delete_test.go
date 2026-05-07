package object

import (
	"testing"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string { return &s }

func TestObjectsToIdentifiers(t *testing.T) {
	tests := []struct {
		name    string
		objects []objectstorage.Object
		want    []objectstorage.ObjectIdentifier
	}{
		{
			name:    "empty slice",
			objects: []objectstorage.Object{},
			want:    []objectstorage.ObjectIdentifier{},
		},
		{
			name: "multiple objects with keys",
			objects: []objectstorage.Object{
				{Key: strPtr("photos/a.jpg")},
				{Key: strPtr("docs/b.pdf")},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "photos/a.jpg"},
				{Key: "docs/b.pdf"},
			},
		},
		{
			name: "object with nil Key returns empty string",
			objects: []objectstorage.Object{
				{Key: nil},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := objectsToIdentifiers(tt.objects)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersionsToIdentifiers(t *testing.T) {
	tests := []struct {
		name     string
		versions []objectstorage.ObjectVersion
		want     []objectstorage.ObjectIdentifier
	}{
		{
			name:     "empty slice",
			versions: []objectstorage.ObjectVersion{},
			want:     []objectstorage.ObjectIdentifier{},
		},
		{
			name: "version with VersionId set",
			versions: []objectstorage.ObjectVersion{
				{Key: strPtr("file.txt"), VersionId: strPtr("v1")},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "file.txt", VersionId: strPtr("v1")},
			},
		},
		{
			name: "version with nil VersionId results in nil VersionId",
			versions: []objectstorage.ObjectVersion{
				{Key: strPtr("file.txt"), VersionId: nil},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "file.txt", VersionId: nil},
			},
		},
		{
			name: "mix of versioned and unversioned",
			versions: []objectstorage.ObjectVersion{
				{Key: strPtr("a.txt"), VersionId: strPtr("v1")},
				{Key: strPtr("b.txt"), VersionId: nil},
				{Key: strPtr("c.txt"), VersionId: strPtr("v3")},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "a.txt", VersionId: strPtr("v1")},
				{Key: "b.txt", VersionId: nil},
				{Key: "c.txt", VersionId: strPtr("v3")},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := versionsToIdentifiers(tt.versions)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDeleteMarkersToIdentifiers(t *testing.T) {
	tests := []struct {
		name    string
		markers []objectstorage.DeleteMarkerEntry
		want    []objectstorage.ObjectIdentifier
	}{
		{
			name:    "empty slice",
			markers: []objectstorage.DeleteMarkerEntry{},
			want:    []objectstorage.ObjectIdentifier{},
		},
		{
			name: "marker with VersionId set",
			markers: []objectstorage.DeleteMarkerEntry{
				{Key: strPtr("removed.txt"), VersionId: strPtr("dm1")},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "removed.txt", VersionId: strPtr("dm1")},
			},
		},
		{
			name: "marker with nil VersionId results in nil VersionId",
			markers: []objectstorage.DeleteMarkerEntry{
				{Key: strPtr("removed.txt"), VersionId: nil},
			},
			want: []objectstorage.ObjectIdentifier{
				{Key: "removed.txt", VersionId: nil},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deleteMarkersToIdentifiers(tt.markers)
			assert.Equal(t, tt.want, got)
		})
	}
}
