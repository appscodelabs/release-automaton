package lib

import (
	"fmt"
	"github.com/Masterminds/semver"
	"k8s.io/apimachinery/pkg/util/sets"
	"sort"
)

func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func Values(m map[string]string) []string {
	values := sets.NewString()
	for _, v := range m {
		values.Insert(v)
	}
	return values.UnsortedList()
}

func SortVersions(versions []string) ([]string, error) {
	vs := make([]*semver.Version, len(versions))
	for i, v := range versions {
		v, err := semver.NewVersion(v)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %s", err)
		}
		vs[i] = v
	}
	sort.Sort(SemverCollection(vs))

	result := make([]string, len(vs))
	for i, v := range vs {
		result[i] = v.Original()
	}
	return result, nil
}
