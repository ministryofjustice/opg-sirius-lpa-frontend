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

		eventDetails, err := client.GetCombinedEvents(ctx, uid)

		for i := range eventDetails {
			formattedUUID, err := LPAEventIDFromUUID(eventDetails[i].ID)
			if err != nil {
				log.Println("Error generating formattedUID:", err)
				continue
			}
			eventDetails[i].FormattedLpaStoreId = formattedUUID
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
