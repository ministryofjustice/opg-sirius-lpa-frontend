package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

const flashCookieName string = "flash-lpa-frontend"

type FlashNotification struct {
	Title       string `json:"name"`
	Description string `json:"description"`
}

func SetFlash(w http.ResponseWriter, notification FlashNotification) {
	str, err := json.Marshal(&notification)
	if err != nil {
		return
	}

	c := &http.Cookie{
		Name:     flashCookieName,
		Value:    base64.URLEncoding.EncodeToString(str),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func GetFlash(w http.ResponseWriter, r *http.Request) (FlashNotification, error) {
	c, err := r.Cookie(flashCookieName)

	if err != nil {
		switch err {
		case http.ErrNoCookie:
			return FlashNotification{}, nil
		default:
			return FlashNotification{}, err
		}
	}

	str, err := base64.URLEncoding.DecodeString(c.Value)
	if err != nil {
		return FlashNotification{}, err
	}

	var v FlashNotification
	err = json.Unmarshal([]byte(str), &v)

	if err != nil {
		return FlashNotification{}, err
	}

	dc := &http.Cookie{Name: flashCookieName, MaxAge: -1, Expires: time.Unix(1, 0)}
	http.SetCookie(w, dc)
	return v, nil
}
