package auth

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jrhoward/cognito-authflow-handler/config"
)

const cookiePaths = "/"

var domain = config.Get("domain")

var LogError = log.New(os.Stdout, "ERROR ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)
var LogInfo = log.New(os.Stdout, "INFO ", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

func AuthWrapper(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code != "" {
		TokenHandler(w, r)
	} else {
		ProtectedHandler(w, r)
	}
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	code := values.Get("code")
	token, err := setCognitoToken(code, false)
	if err != nil {
		LogError.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	setAuthCookies(w, token)
	http.Redirect(w, r, config.Get("authHandlerRedirect"), http.StatusSeeOther)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	id_token, err := r.Cookie(config.Get("idCookieName"))
	if err != nil {
		LogError.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenString := id_token.Value
	err = validateCognitoToken(tokenString)
	if err == nil {
		return
	} else if strings.Contains(err.Error(), "TOKEN EXPIRED") {
		refresh_token, err := r.Cookie(config.Get("refreshCookieName"))
		if err != nil {
			LogError.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			expireCookies(w)
			return
		}
		token, err := setCognitoToken(refresh_token.Value, true)
		if err != nil {
			LogError.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			expireCookies(w)
			return
		}
		setAuthCookies(w, token)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		LogError.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := r.Cookie(config.Get("refreshCookieName"))
	if err != nil {
		LogError.Println(err)
	} else {
		err = revokeToken(refresh_token.Value)
		if err != nil {
			LogError.Println(err)
		}
	}
	expireCookies(w)
	http.Redirect(w, r, config.Get("logoutHandlerRedirect"), http.StatusSeeOther)
}

func setAuthCookies(w http.ResponseWriter, token Tokens) {
	cookieMaxAge := config.GetCookieMaxAge()
	http.SetCookie(w, &http.Cookie{Name: config.Get("idCookieName"), Value: token.Id_token,
		SameSite: http.SameSiteLaxMode, HttpOnly: true, Secure: true, Domain: domain, Path: cookiePaths, MaxAge: cookieMaxAge})
	http.SetCookie(w, &http.Cookie{Name: config.Get("refreshCookieName"), Value: token.Refresh_token,
		SameSite: http.SameSiteLaxMode, HttpOnly: true, Secure: true, Domain: domain, Path: cookiePaths, MaxAge: cookieMaxAge})
}

func expireCookies(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: config.Get("idCookieName"), MaxAge: -1})
	http.SetCookie(w, &http.Cookie{Name: config.Get("refreshCookieName"), MaxAge: -1})
}
