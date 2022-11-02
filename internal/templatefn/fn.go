package templatefn

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

func All(siriusPublicURL, prefix, staticHash string) map[string]interface{} {
	return map[string]interface{}{
		"sirius": func(s string) string {
			return siriusPublicURL + s
		},
		"prefix": func(s string) string {
			return prefix + s
		},
		"prefixAsset": func(s string) string {
			if len(staticHash) >= 11 {
				return prefix + s + "?" + url.QueryEscape(staticHash[3:11])
			} else {
				return prefix + s
			}
		},
		"today": func() string {
			return time.Now().Format("2006-01-02")
		},
		"field": func(name, label string, value interface{}, error map[string]string, attrs ...interface{}) map[string]interface{} {
			field := map[string]interface{}{
				"name":  name,
				"label": label,
				"value": value,
				"error": error,
			}

			if len(attrs)%2 != 0 {
				panic("must have even number of attrs")
			}

			for i := 0; i < len(attrs); i += 2 {
				field[attrs[i].(string)] = attrs[i+1]
			}

			return field
		},
		"fee": func(amount int) string {
			float := float64(amount)
			return fmt.Sprintf("%.2f", float/100)
		},
		"formatDate": func(s sirius.DateString) (string, error) {
			if s != "" {
				return s.ToSirius()
			}
			return "", nil
		},
		"translateRefData": func(types []sirius.RefDataItem, tmplHandle string) string {
			for _, refDataType := range types {
				if refDataType.Handle == tmplHandle {
					return refDataType.Label
				}
			}
			return tmplHandle
		},
		"ToLower": strings.ToLower,
	}
}
