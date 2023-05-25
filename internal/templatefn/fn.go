package templatefn

import (
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
		"radios":  radios,
		"item":    item,
		"fieldID": fieldID,
		"select":  select_,
		"options": options,
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
		"ToUpper": strings.ToUpper,
		"capitalise": func(text string) string {
			return cases.Title(language.English).String(text)
		},
		"contains": func(xs []string, needle string) bool {
			for _, x := range xs {
				if x == needle {
					return true
				}
			}
			return false
		},
		"minus1": func(i int) int {
			return i - 1
		},
		"statusColour": func(s string) string {
			switch s {
			case "registered":
				return "green"
			case "perfect":
				return "turquoise"
			case "pending":
				return "blue"
			case "payment pending", "reduced fees pending":
				return "purple"
			case "cancelled", "rejected", "revoked", "withdrawn", "return - unpaid", "deleted":
				return "red"
			default:
				return "grey"
			}
		},
		"replace": func(s, find, replace string) string {
			return strings.ReplaceAll(s, find, replace)
		},
		"dateYear": func(s sirius.DateString) (string, error) {
			if s != "" {
				return s.GetYear()
			}
			return "", nil
		},
		"filterContent": func(content string) string {
			//Fixes extra newline appearing in text editor due to newline present between the doctype and html tags
			return strings.Replace(content, "<!DOCTYPE html>\n<html lang=\"en\">", "<!DOCTYPE html><html lang=\"en\">", -1)
		},
		"abs": func(num int) int {
			if num < 0 {
				return -num
			}
			return num
		},
		"attr": func(s string) template.HTMLAttr {
			return template.HTMLAttr(s)
		},
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
	return fieldData{
		Name:  name,
		Label: label,
		Value: value,
		Error: error,
		Attrs: collectAttrs(attrs),
	}
}

type radiosData struct {
	Name   string
	Label  string
	Value  interface{}
	Errors map[string]string
	Items  []itemData
}

func radios(name, label string, value interface{}, errors map[string]string, items ...itemData) radiosData {
	return radiosData{
		Name:   name,
		Label:  label,
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
	return itemData{
		Value: value,
		Label: label,
		Attrs: collectAttrs(attrs),
	}
}

func fieldID(name string, i int) string {
	if i == 0 {
		return name
	}

	return fmt.Sprintf("%s-%d", name, i+1)
}

type selectData struct {
	Name    string
	Label   string
	Value   interface{} // string|int
	Errors  map[string]string
	Options []optionData
	Attrs   map[string]interface{}
}

func select_(name, label string, value interface{}, errors map[string]string, options []optionData, attrs ...interface{}) selectData {
	return selectData{
		Name:    name,
		Label:   label,
		Value:   value,
		Errors:  errors,
		Options: options,
		Attrs:   collectAttrs(attrs),
	}
}

type optionData struct {
	Value interface{} // string|int
	Label string
}

func options(list interface{}, attrs ...interface{}) []optionData {
	attributes := collectAttrs(attrs)
	var datas []optionData

	switch v := list.(type) {
	case []string:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item, Label: item}
		}

	case []sirius.MiConfigEnum:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item.Name, Label: item.Description}
		}

	case []sirius.RefDataItem:
		if attributes["filterSelectable"] == true {
			for _, item := range v {
				if item.UserSelectable {
					datas = append(datas, optionData{Value: item.Handle, Label: item.Label})
				}
			}
		} else {
			datas = make([]optionData, len(v))
			for i, item := range v {
				datas[i] = optionData{Value: item.Handle, Label: item.Label}
			}
		}

	case []sirius.Team:
		datas = make([]optionData, len(v))
		for i, item := range v {
			datas[i] = optionData{Value: item.ID, Label: item.DisplayName}
		}
	}

	return datas
}

func collectAttrs(attrs []interface{}) map[string]interface{} {
	attributes := map[string]interface{}{}
	if len(attrs)%2 != 0 {
		panic("must have even number of attrs")
	}

	for i := 0; i < len(attrs); i += 2 {
		attributes[attrs[i].(string)] = attrs[i+1]
	}

	return attributes
}
