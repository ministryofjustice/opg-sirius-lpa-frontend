package server

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ministryofjustice/opg-sirius-lpa-frontend/internal/sirius"

	"github.com/ministryofjustice/opg-go-common/template"
)

type GetHistoryClient interface {
	GetEvents(ctx sirius.Context, donorId int, caseId int) (any, error)
	GetCombinedEvents(ctx sirius.Context, uid string) (sirius.APIEvent, error)
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

		var eventDetails sirius.APIEvent
		var traditionalEventDetails any

		if caseSummary.DigitalLpa.LpaStoreData.Status.ReadableString() != "" {
			// Digital LPA - use combined events
			eventDetails, err = client.GetCombinedEvents(ctx, uid)
			for i := range eventDetails {
				formattedUUID, err := LPAEventIDFromUUID(eventDetails[i].UUID)
				if err != nil {
					log.Println("Error generating formattedUUID:", err)
					continue
				}
				eventDetails[i].FormattedUUID = formattedUUID
			}
		} else {
			// Traditional LPA - use Sirius events only
			donorId := caseSummary.DigitalLpa.SiriusData.Donor.ID
			caseId := caseSummary.DigitalLpa.SiriusData.ID
			traditionalEventDetails, err = client.GetEvents(ctx, donorId, caseId)
		}

		if err != nil {
			return err
		}

		data := getHistory{
			CaseSummary: caseSummary,
		}
		if eventDetails != nil {
			data.EventData = eventDetails
		} else if traditionalEventDetails != nil {
			data.EventData = traditionalEventDetails
		}

		return tmpl(w, data)
	}
}

func LPAEventIDFromUUID(uuidStr string) (string, error) {
	clean := strings.ReplaceAll(uuidStr, "-", "")

	idBytes, err := hex.DecodeString(clean)
	if err != nil {
		return "", err
	}

	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	base32Str := encoder.EncodeToString(idBytes)

	if len(base32Str) < 8 {
		return "", fmt.Errorf("unexpected Base32 length for UUID")
	}
	return base32Str[:8], nil
}
