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
		"field":   field,
		"items":   items,
		"item":    item,
		"fieldID": fieldID,
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

type fieldData struct {
	Name  string
	Label string
	Value interface{}
	Error map[string]string
	Attrs map[string]interface{}
}

func field(name, label string, value interface{}, error map[string]string, attrs ...interface{}) fieldData {
	field := fieldData{
		Name:  name,
		Label: label,
		Value: value,
		Error: error,
		Attrs: map[string]interface{}{},
	}

	if len(attrs)%2 != 0 {
		panic("must have even number of attrs")
	}

	for i := 0; i < len(attrs); i += 2 {
		field.Attrs[attrs[i].(string)] = attrs[i+1]
	}

	return field
}

type itemsData struct {
	Name   string
	Value  interface{}
	Errors map[string]string
	Items  []itemData
}

func items(name string, value interface{}, errors map[string]string, items ...itemData) itemsData {
	return itemsData{
		Name:   name,
		Value:  value,
		Errors: errors,
		Items:  items,
	}
}

type itemData struct {
	Value string
	Label string
	Attrs map[string]interface{}
}

func item(value, label string, attrs ...interface{}) itemData {
	item := itemData{
		Value: value,
		Label: label,
		Attrs: map[string]interface{}{},
	}

	if len(attrs)%2 != 0 {
		panic("must have even number of attrs")
	}

	for i := 0; i < len(attrs); i += 2 {
		item.Attrs[attrs[i].(string)] = attrs[i+1]
	}

	return item
}

func fieldID(name string, i int) string {
	if i == 0 {
		return name
	}

	return fmt.Sprintf("%s-%d", name, i+1)
}
