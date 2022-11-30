package server

import (
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
	"golang.org/x/exp/slices"
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
	DocumentInsertTypes     []InsertDisplayData
	InsertCategories        []string
	TemplateSelected        sirius.DocumentTemplateData
}

type InsertDisplayData struct {
	Handle string
	Label  string
	Key    string
}

func CreateDocument(client CreateDocumentClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if err := r.ParseForm(); err != nil {
			return err
		}
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

			templateId := r.FormValue("templateId")
			if templateId != "" {
				for _, dt := range data.DocumentTemplates {
					if dt.TemplateId == templateId {
						data.TemplateSelected = dt
						break
					}
				}
			}
			if data.TemplateSelected.TemplateId != "" {
				data.DocumentInsertTypes, data.InsertCategories = translateInsertData(data.TemplateSelected.Inserts, data.DocumentTemplateRefData)
			}

			//inserts := r.Form["insert"]
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

func translateInsertData(selectedTemplateInserts []sirius.Insert, documentTemplateRefData []sirius.RefDataItem) ([]InsertDisplayData, []string) {
	var documentTemplateInserts []InsertDisplayData
	var insertCategories []string
	for _, in := range selectedTemplateInserts {
		if !slices.Contains(insertCategories, in.Key) {
			insertCategories = append(insertCategories, in.Key)
		}
		for _, refData := range documentTemplateRefData {
			if refData.Handle == in.OnScreenSummary {
				translatedRefDataItem := InsertDisplayData{
					Handle: in.InsertId,
					Label:  refData.Label,
					Key:    in.Key,
				}
				documentTemplateInserts = append(documentTemplateInserts, translatedRefDataItem)
			}
		}
	}
	return documentTemplateInserts, insertCategories
}
