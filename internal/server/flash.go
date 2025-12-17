package server

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

const flashCookieName string = "flash-lpa-frontend"

var secureCookies bool = os.Getenv("INSECURE_COOKIES") != "1"

type FlashNotification struct {
	Title string `json:"name"`
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
		Path:     "/",
		Secure:   secureCookies,
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

	dc := &http.Cookie{Name: flashCookieName, MaxAge: -1, Expires: time.Unix(1, 0), Path: "/", Secure: true}

	http.SetCookie(w, dc)
	return v, nil
}
