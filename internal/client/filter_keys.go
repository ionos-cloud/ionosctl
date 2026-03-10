package client

import (
	"reflect"
	"strings"
	"sync"

	"github.com/ionos-cloud/ionosctl/v6/internal/filterprops"
)

var (
	filterKeyMap  map[string]string // lowercase → correct camelCase
	filterKeyOnce sync.Once
)

// normalizeFilterKey maps any-cased filter key to the correct camelCase form
// expected by the IONOS Cloud API. If the key is unknown, it is returned as-is.
// E.g. "IMaGeTYpE" → "imageType", "imagetype" → "imageType", "imageType" → "imageType"
func normalizeFilterKey(key string) string {
	filterKeyOnce.Do(func() {
		filterKeyMap = buildFilterKeyMap()
	})

	if canonical, ok := filterKeyMap[strings.ToLower(key)]; ok {
		return canonical
	}
	return key
}

func buildFilterKeyMap() map[string]string {
	m := make(map[string]string)

	for _, t := range filterprops.AllCloudAPIv6Types() {
		rt := reflect.TypeOf(t)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		if rt.Kind() != reflect.Struct {
			continue
		}
		for i := 0; i < rt.NumField(); i++ {
			camelCase := filterprops.FirstLower(rt.Field(i).Name)
			m[strings.ToLower(camelCase)] = camelCase
		}
	}

	return m
}
