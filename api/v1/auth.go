package v1

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/snehmatic/mindloop/internal/config"
)

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()

		if !uc.AuthConfig.Enabled {
			next.ServeHTTP(w, r)
			return
		}

		// Allow static files and login endpoints to pass through
		if r.URL.Path == "/login" || strings.HasPrefix(r.URL.Path, "/static/") || r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		// Check for cookie
		cookie, err := r.Cookie("mindloop_session")
		if err != nil || cookie.Value != uc.AuthConfig.PasswordHash {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (mlh *MindloopHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Just render a simple login template
		// We might need to ensure layout handles missing partials or we use a separate layout.
		mlh.renderTemplate(w, "login.html", map[string]interface{}{
			"Title": "Login",
		})
		return
	}

	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		uc := config.UserConfig{}
		_ = uc.ReadFromYAML()

		if hashPassword(password) == uc.AuthConfig.PasswordHash {
			http.SetCookie(w, &http.Cookie{
				Name:     "mindloop_session",
				Value:    uc.AuthConfig.PasswordHash,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
				Path:     "/",
			})
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		mlh.renderTemplate(w, "login.html", map[string]interface{}{
			"Title":        "Login",
			"ErrorMessage": "Invalid password",
		})
	}
}

func (mlh *MindloopHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "mindloop_session",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
