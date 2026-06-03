package server

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

// TODO: implement logic
func SiriusHeaderPeopleInfo(tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return tmpl(w, struct{}{})
	}
}
