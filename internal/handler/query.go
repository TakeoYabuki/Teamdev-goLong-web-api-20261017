package handler

import (
	"fmt"
	"maps"
	"net/url"
	"slices"
)

func checkQuery(query url.Values, allowed ...string) error {
	for _, name := range slices.Sorted(maps.Keys(query)) {
		values := query[name]

		if !slices.Contains(allowed, name) {
			return fmt.Errorf("%sは使えないパラメータです", name)
		}
		if len(values) > 1 {
			return fmt.Errorf("%sは1つだけ指定してください", name)
		}
		if values[0] == "" {
			return fmt.Errorf("%sに値を指定してください", name)
		}
	}
	return nil
}
