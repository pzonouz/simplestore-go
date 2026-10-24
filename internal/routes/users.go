package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simplestore_go/internal/services"
	"simplestore_go/utils"
	"time"

	sqlc "simplestore_go/internal/db"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wneessen/go-mail"
)

func RegisterUsersRoutes(mux *chi.Mux, service services.Service) {
	mux.Route("/auth", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			var userFromRequest sqlc.CreateUserParams
			err := json.NewDecoder(r.Body).Decode(&userFromRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			user, err := service.CreateUsers(r.Context(), userFromRequest)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			m := mail.NewMsg()
			if err := m.From("master@nanochemshop.ir"); err != nil {
				log.Printf("failed to set From: %v", err)
			}
			if err := m.To(user.Email.String); err != nil {
				log.Printf("failed to set To: %v", err)
			}
			m.Subject("ایمیل فعالسازی Nanochemshop.ir")

			// 2. Set the HTML body
			htmlBody := `<html><body><h1>Hello!</h1><p>برای فعالسازی روی </a>لینک<a href="http://localhost/auth/activation/` + user.Token.String + `"> کلیک کنید</p></body></html>`
			m.SetBodyString(mail.TypeTextHTML, htmlBody)

			// 3. Configure the client (example with a generic SMTP server)
			c, err := mail.NewClient("smtp.gmail.com",
				mail.WithPort(587),
				mail.WithSMTPAuth(mail.SMTPAuthPlain),
				mail.WithUsername(os.Getenv("GMAIL_USERNAME")),
				mail.WithPassword(os.Getenv("GMAIL_PASSWORD")),
				mail.WithTLSPolicy(mail.TLSMandatory),
			)
			if err != nil {
				log.Printf("failed to create client: %v", err)
			}

			// 4. Send the email
			if err := c.DialAndSend(m); err != nil {
				log.Printf("failed to send email: %v", err)
			}
			log.Println("HTML email sent successfully!")
			w.WriteHeader(http.StatusCreated)

		})
		r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
			var user sqlc.User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			userFromDB, err := service.GetUserByEmail(r.Context(), user.Email)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			fmt.Print(userFromDB.Password)
			ok, err := argon2id.ComparePasswordAndHash(user.Password, userFromDB.Password)
			if err != nil {
				fmt.Println(err.Error())
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			if !ok {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				panic("JWT_SECRET environment variable is required")

			}

			t := jwt.New(jwt.SigningMethodHS256)

			jwtClaims := utils.JWTClaims{
				Is_admin:  userFromDB.IsAdmin.Bool,
				Email:     userFromDB.Email.String,
				User_id:   userFromDB.ID.String(),
				Issuer:    "nanochemshop",
				Audience:  jwt.ClaimStrings{"nanochemshop-web"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			}
			t.Claims = jwtClaims
			token, err := t.SignedString([]byte(secret))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Add("Content-Type", "application/json")
			http.SetCookie(w, &http.Cookie{
				Name:     "access",
				Value:    token,
				Path:     "/",
				Expires:  jwtClaims.ExpiresAt.Time,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			})
			w.WriteHeader(http.StatusCreated)
		})
		r.Post("/activation", func(w http.ResponseWriter, r *http.Request) {
			type TokenBody struct {
				Token string `json:"token"`
			}
			var tokenBody TokenBody
			err := json.NewDecoder(r.Body).Decode(&tokenBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			err = service.ActivateUserByToken(r.Context(), pgtype.Text{String: tokenBody.Token, Valid: true})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
		})
		r.With(utils.AuthenticatedMiddleware).Get("/me", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			user, ok := r.Context().Value("user").(*utils.UserInfo)
			if !ok {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
			err := json.NewEncoder(w).Encode(user)
			if err != nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}
		})
	})
}
