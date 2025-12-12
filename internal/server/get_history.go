package server

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"
)

type GetHistoryClient interface {
	GetCombinedEvents(ctx sirius.Context, uid string) (sirius.APIEvents, error)
	CaseSummary(ctx sirius.Context, uid string) (sirius.CaseSummary, error)
}

type getHistory struct {
	CaseSummary sirius.CaseSummary
	EventData   any
}

func GetHistory(client GetHistoryClient, tmpl template.Template) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		uid := r.PathValue("uid")
		ctx := getContext(r)

		caseSummary, err := client.CaseSummary(ctx, uid)
		if err != nil {
			return err
		}

		eventDetails, err := client.GetCombinedEvents(ctx, uid)

		for i := range eventDetails {
			if eventDetails[i].IsLpaStore() {
				formattedUUID, _ := LPAEventIDFromUUID(eventDetails[i].ID)
				eventDetails[i].FormattedLpaStoreId = formattedUUID

				lsc := getLpaStoreCategoryFromChanges(eventDetails[i].Changes)
				eventDetails[i].Category = lsc.Readable()

				for j, c := range eventDetails[i].Changes {
					ct := getLpaStoreChangeTypeFromChange(c)

					eventDetails[i].Changes[j].Template = ct.GetTemplate()
					eventDetails[i].Changes[j].Readable = ct.Readable()
				}
			}
		}

		if err != nil {
			return err
		}

		data := getHistory{
			CaseSummary: caseSummary,
			EventData:   eventDetails,
		}

		return tmpl(w, data)
	}
}

func LPAEventIDFromUUID(id string) (string, error) {
	clean := strings.ReplaceAll(id, "-", "")

	idBytes, err := hex.DecodeString(clean)
	if err != nil {
		return "", err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	base32Str := encoder.EncodeToString(idBytes)

	if len(base32Str) < 8 {
		return "", fmt.Errorf("unexpected Base32 length for UID")
	}
	return base32Str[:8], nil
}
