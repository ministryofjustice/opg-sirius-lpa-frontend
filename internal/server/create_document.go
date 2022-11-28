package server

import (
	"fmt"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
)

type CreateDocumentClient interface {
	Case(ctx sirius.Context, id int) (sirius.Case, error)
	RefDataByCategory(ctx sirius.Context, category string) ([]sirius.RefDataItem, error)
	DocumentTemplates(ctx sirius.Context, caseType sirius.CaseType) ([]sirius.DocumentTemplateData, error)
}

type createDocumentData struct {
	XSRFToken               string
	Success                 bool
	Error                   sirius.ValidationError
	Case                    sirius.Case
	Document                sirius.Document
	DocumentTemplates       []sirius.DocumentTemplateData
	DocumentTemplateTypes   []sirius.RefDataItem
	DocumentTemplateRefData []sirius.RefDataItem
	TemplateId              string
	TemplateSelected        sirius.DocumentTemplateData
}

func CreateDocument(client CreateDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := getContext(r)
		data := createDocumentData{
			XSRFToken: ctx.XSRFToken,
		}

		switch r.Method {
		case http.MethodGet:
			caseID, err := strconv.Atoi(r.FormValue("id"))
			if err != nil {
				return err
			}

			caseType, err := sirius.ParseCaseType(r.FormValue("case"))
			if err != nil {
				return err
			}

			group, groupCtx := errgroup.WithContext(ctx.Context)

			group.Go(func() error {
				caseItem, err := client.Case(ctx.With(groupCtx), caseID)
				if err != nil {
					return err
				}

				data.Case = caseItem
				return nil
			})

			group.Go(func() error {
				documentTemplates, err := client.DocumentTemplates(ctx.With(groupCtx), caseType)
				if err != nil {
					return err
				}

				data.DocumentTemplates = documentTemplates
				return nil
			})

			group.Go(func() error {
				documentTemplateRefData, err := client.RefDataByCategory(ctx, sirius.DocumentTemplateIdCategory)
				if err != nil {
					return err
				}
				data.DocumentTemplateRefData = documentTemplateRefData
				return nil
			})

			if err := group.Wait(); err != nil {
				return err
			}

			data.DocumentTemplateTypes = translateDocumentData(data.DocumentTemplates, data.DocumentTemplateRefData)

			data.TemplateId = r.FormValue("templateId")

			fmt.Println(data.TemplateSelected)

			if data.TemplateId != "" {
				for _, dt := range data.DocumentTemplates {
					if dt.TemplateId == data.TemplateId {
						data.TemplateSelected = dt
						break
					}
				}
			}

			fmt.Println(data.TemplateSelected)

		case http.MethodPost:
		}

		return tmpl(w, data)
	}
}

func translateDocumentData(documentTemplateData []sirius.DocumentTemplateData, documentTemplateRefData []sirius.RefDataItem) []sirius.RefDataItem {
	var documentTemplateTypes []sirius.RefDataItem
	for _, dt := range documentTemplateData {
		for _, refData := range documentTemplateRefData {
			if refData.Handle == dt.OnScreenSummary {
				translatedRefDataItem := sirius.RefDataItem{
					Handle:         dt.TemplateId,
					Label:          refData.Label,
					UserSelectable: false,
				}
				documentTemplateTypes = append(documentTemplateTypes, translatedRefDataItem)
			}
		}
	}

	return documentTemplateTypes
}
