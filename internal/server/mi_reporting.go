package server

import (
	"net/http"
	"net/url"
	"sort"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type MiReportingClient interface {
	MiConfig(sirius.Context) (map[string]sirius.MiConfigProperty, error)
	MiReport(sirius.Context, url.Values) (*sirius.MiReportResponse, error)
}

type miReportingData struct {
	ReportTypes []sirius.MiConfigEnum
	ReportType  string
	ReportName  string
	Controls    []namedControl
	ResultCount int
	Download    string
}

type namedControl struct {
	Name       string
	Label      string
	Properties sirius.MiConfigProperty
}

var miLabels = map[string]string{
	"applicationType":   "Application",
	"endDate":           "To",
	"optionalEndDate":   "To",
	"optionalStartDate": "From",
	"paymentSource":     "Payment source",
	"reportSubject":     "Report subject",
	"reportType":        "Report type",
	"source":            "Source",
	"startDate":         "From",
	"status":            "Status",
	"state":             "Status",
}

func MiReporting(client MiReportingClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		data := miReportingData{}

		switch r.Method {
		case http.MethodGet:
			config, err := client.MiConfig(ctx)
			if err != nil {
				return err
			}

			data.ReportTypes = config["reportType"].Enum
			data.ReportType = r.FormValue("reportType")

			if data.ReportType != "" {
				for _, report := range data.ReportTypes {
					if report.Name == data.ReportType {
						data.ReportName = report.Description
						break
					}
				}

				for name, properties := range config {
					for _, reportType := range properties.DependsOn.ReportType {
						if reportType.Name == data.ReportType {
							data.Controls = append(data.Controls, namedControl{Name: name, Label: miLabels[name], Properties: properties})
						}
					}
				}

				orderMap := map[string]int{}
				for i, name := range []string{"reportType", "applicationType", "source", "status", "paymentSource", "startDate", "endDate", "state", "optionalStartDate", "optionalEndDate"} {
					orderMap[name] = i
				}

				sort.Slice(data.Controls, func(i, j int) bool {
					return orderMap[data.Controls[i].Name] < orderMap[data.Controls[j].Name]
				})
			}

		case http.MethodPost:
			form := r.PostForm

			for _, key := range []string{"startDate", "endDate", "optionalStartDate", "optionalEndDate"} {
				if form.Has(key) {
					value := form.Get(key)
					if value == "" {
						form.Del(key)
						continue
					}

					date, err := sirius.DateString(value).ToSirius()
					if err != nil {
						return err
					}
					form.Set(key, date)
				}
			}

			for _, key := range []string{"applicationType", "paymentSource", "source", "status", "state"} {
				if form.Has(key) {
					for _, value := range form[key] {
						form.Add(key+"[]", value)
					}
					form.Del(key)
				}
			}

			result, err := client.MiReport(ctx, form)
			if err != nil {
				return err
			}

			data.ResultCount = result.ResultCount
			data.ReportType = result.ReportType
			data.ReportName = result.ReportDescription

			form.Add("OPG-Bypass-Membrane", "1")
			data.Download = "/api/reporting/export?" + form.Encode()
		}

		return tmpl(w, data)
	}
}
