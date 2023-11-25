package routes

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	"github.com/fbuedding/iota-admin/internal/pkg/cookies"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type Credentials struct {
	Username string `schema:"username,required"`
	Password string `schema:"password,required"`
}

func Auth(a auth.Authenticator, st sessionStore.SessionStore) chi.Router {
	r := chi.NewRouter()
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var decoder = schema.NewDecoder()
		err := r.ParseForm()
		if err != nil {
			log.Error().Err(err).Msg("Could not parse form")
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}

		var cred Credentials
		err = decoder.Decode(&cred, r.PostForm)
		if err != nil {
			log.Error().Err(err).Msg("Could not decode form")
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}
		usr, err := a.Login(auth.Username(cred.Username), auth.Password(cred.Password))
		if err != nil {
			log.Error().Err(err).Msg("Could not log in user")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		session := sessionStore.Session{
			Username: usr.Username,
			Expiry:   time.Now().Add(120 * time.Second),
		}
		sessionToken, err := st.Add(&session)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookies.WriteSigned(w, cookies.New("session_token", string(sessionToken), session.Expiry),
			getCookieSecret())

		w.Header().Add("HX-Redirect", "/index")

		w.Write([]byte(fmt.Sprintf("<div>Hallo %v</div>", cred.Username)))
	})

	r.Delete("/login", func(w http.ResponseWriter, r *http.Request) {

		sessionToken, err := cookies.ReadSigned(r, "session_token", getCookieSecret())
		cookies.Delete(w, "session_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			case errors.Is(err, cookies.ErrInvalidValue):
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			default:
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}
		st.Remove(sessionStore.SessionToken(sessionToken))
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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

			sessionToken, err := cookies.ReadSigned(r, "session_token", []byte(globals.Conf.CookieSecret))
			if err != nil {
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
			cookies.WriteSigned(w, &http.Cookie{
				Name:     "session_token",
				Value:    string(sessionToken),
				HttpOnly: true,
				Expires:  session.Expiry,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
				Path:     "/",
			},
				getCookieSecret())
			ctx := context.WithValue(r.Context(), "user", string(session.Username))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getCookieSecret() []byte {
	os.Getenv("COOKIE_SECRET")
	secretKey, err := hex.DecodeString(os.Getenv("COOKIE_SECRET"))
	if err != nil {
		panic(err)
	}
	return secretKey
}

func handleUnauthorized(w http.ResponseWriter, r *http.Request) {
	log.Debug().Str("HX-Request", r.Header.Get("HX-Request")).Msg("Handling Auth")
	if r.Header.Get("HX-Request") == "true" {

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
