package routes

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	bruteforceprotection "github.com/fbuedding/iota-admin/internal/pkg/bruteForceProtection"
	"github.com/fbuedding/iota-admin/internal/pkg/cookies"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type Credentials struct {
	Username string `schema:"username,required"`
	Password string `schema:"password,required"`
}

func Auth(a auth.Authenticator, st sessionStore.SessionStore, bfp bruteforceprotection.BrutForceProtection) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		decoder := schema.NewDecoder()
		err := r.ParseForm()
		if err != nil {
			log.Error().Err(err).Msg("Could not parse form")
			templates.HandleError(r.Context(), w, err, 400)
			return
		}

		var cred Credentials
		err = decoder.Decode(&cred, r.PostForm)
		if err != nil {
			log.Error().Err(err).Msg("Could not decode form")
			templates.HandleError(r.Context(), w, err, 400)
			return
		}
		if bfp.IsBlocked(auth.Username(cred.Username)) {
			templates.HandleError(r.Context(), w, fmt.Errorf("To many login attempts"), http.StatusUnauthorized)
			return
		}
		usr, err := a.Authenticate(auth.Username(cred.Username), auth.Password(cred.Password))
		if err != nil {
			log.Debug().Err(err).Str("user", cred.Username).Msg("Could not log in user")
			bfp.Hit(auth.Username(cred.Username))
			templates.HandleError(r.Context(), w, err, http.StatusUnauthorized)
			return
		}
		bfp.Delete(auth.Username(cred.Username))
		session := sessionStore.Session{
			Username: usr.Username,
			Expiry:   time.Now().Add(120 * time.Second),
		}
		sessionToken, err := st.Add(&session)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}
		err = cookies.WriteSigned(w, cookies.New("session_token", string(sessionToken), session.Expiry),
			getCookieSecret())
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}

		log.Debug().Any("User", usr).Msg("Login succesfull!")
		w.Header().Add("HX-Redirect", "/index")
	})

	r.Delete("/login", func(w http.ResponseWriter, r *http.Request) {
		sessionToken, err := cookies.ReadSigned(r, "session_token", getCookieSecret())
		log.Debug().Str("Session Token:", sessionToken).AnErr("Error:", err).Msg("delete login")
		cookies.Delete(w, "session_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				handleUnauthorized(w, r)
			case errors.Is(err, cookies.ErrInvalidValue):
				handleUnauthorized(w, r)
			default:
				handleUnauthorized(w, r)
			}
			return
		}
		st.Remove(sessionStore.SessionToken(sessionToken))
		w.Header().Add("HX-Redirect", "/login")
		return
	})
	return r
}

func AuthMiddleware(st sessionStore.SessionStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if globals.Conf.BypassAuth {
				ctx := context.WithValue(r.Context(), "user", "fbuedding")

				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			sessionToken, err := cookies.ReadSigned(r, "session_token", getCookieSecret())
			if err != nil {
				log.Error().Err(err).Msg("Error while reading cookie!")
				cookies.Delete(w, "session_token")
				switch {
				case errors.Is(err, http.ErrNoCookie):
					handleUnauthorized(w, r)
				case errors.Is(err, cookies.ErrInvalidValue):
					handleUnauthorized(w, r)
				default:
					handleUnauthorized(w, r)
				}
				return
			}

			session, err := st.Get(sessionStore.SessionToken(sessionToken))
			if err != nil {
				cookies.Delete(w, "session_token")
				handleUnauthorized(w, r)
				return
			}
			if session.IsExpired() {
				st.Remove(sessionStore.SessionToken(sessionToken))
				cookies.Delete(w, "session_token")
				handleUnauthorized(w, r)
				return
			}
			session.Refresh(time.Now().Add(2 * time.Minute))
			err = cookies.WriteSigned(w, &http.Cookie{
				Name:     "session_token",
				Value:    string(sessionToken),
				HttpOnly: true,
				Expires:  session.Expiry,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
				Path:     "/",
			},
				getCookieSecret())
			if err != nil {
				templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), "user", string(session.Username))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getCookieSecret() []byte {
	secretKey, err := hex.DecodeString(globals.Conf.CookieSecret)
	if err != nil {
		panic(err)
	}
	return secretKey
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("HX-Request", r.Header.Get("HX-Request")).Msg("Handling unauthorized")
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Add("HX-Redirect", "/login")
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
