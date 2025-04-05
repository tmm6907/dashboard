package routes

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmm6907/dashboard/auth"
	"github.com/tmm6907/dashboard/models"
	"golang.org/x/oauth2"
)

type JWTClaims struct {
	OAuthID    string
	FirstName  string
	LastName   string
	AuthCode   string
	Expiration int64
	IAT        int64
}

func NewJWTClaims(oauthID string, fname string, lname string, authCode string) *JWTClaims {
	return &JWTClaims{
		OAuthID:    oauthID,
		FirstName:  fname,
		LastName:   lname,
		AuthCode:   authCode,
		Expiration: time.Now().Add(time.Hour * 24).Unix(),
		IAT:        time.Now().Unix(),
	}
}

type Handler struct {
	db          *sqlx.DB
	oauthConfig *oauth2.Config
}

func NewHandler(db *sqlx.DB, oauthConfig *oauth2.Config) *Handler {
	return &Handler{
		db:          db,
		oauthConfig: oauthConfig,
	}
}

func (h *Handler) SaveNewToken(c *fiber.Ctx, oauthID string, firstName string, lastName string, code string) error {
	token, err := h.GenerateJWTToken(oauthID, firstName, lastName, code)
	if err != nil {
		return err
	}
	log.Debug("Token created: ", token)

	cookie := &fiber.Cookie{
		Name:     "token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	}

	if c.Protocol() == "https" {
		cookie.Secure = true
	}
	c.Cookie(cookie)
	return nil
}

func (h *Handler) FetchGoogleInfo(code string) (map[string]string, error) {
	// Fetch user info from Google
	userInfo, err := auth.GetUserInfo(h.GetOauthConfig(), code)
	if err != nil {
		return nil, err
	}

	// Extract relevant details
	oauthID, ok := userInfo["id"].(string)
	if !ok {
		return nil, errors.New("Unable to retrieve oauthID")
	}
	firstName, ok := userInfo["given_name"].(string) // First name
	if !ok {
		return nil, errors.New("Unable to retrieve oauth first name")
	}
	lastName, ok := userInfo["family_name"].(string) // Last name
	if !ok {
		return nil, errors.New("Unable to retrieve oauth last name")
	}
	return map[string]string{
		"oauthID":   oauthID,
		"firstName": firstName,
		"lastName":  lastName,
	}, nil
}

func (h *Handler) ParseTokenString(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
}

func (h *Handler) ParseExpiredTokenString(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	return parser.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
}

func (h *Handler) IsValidJWT(tokenString string) bool {
	token, err := h.ParseTokenString(tokenString)
	if err != nil || !token.Valid {
		return false
	}
	return true
}

func (h *Handler) CheckAuthHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies("token")
		if tokenString == "" {
			log.Error("token is empty")
			redirectURL := auth.GetLoginURL(h.GetOauthConfig(), "random-state")
			return c.Status(http.StatusFound).SendString(redirectURL)
		}
		if !h.IsValidJWT(tokenString) {
			log.Error("invalid token, resetting")
			tokenData, err := h.GetJWTInfo(tokenString, true)
			if err != nil {
				log.Error("Unable to get user info from token: ", err)
				return c.Status(http.StatusInternalServerError).SendString("Unable to get user info from token")
			}
			if err := h.SaveNewToken(c, tokenData["user_id"], tokenData["fname"], tokenData["lname"], tokenData["auth_code"]); err != nil {
				log.Error(err)
				return c.Status(http.StatusInternalServerError).SendString("Unable to generate auth token")
			}
			var user models.User
			db := h.GetDB()
			if err := db.Get(&user, "SELECT * FROM users WHERE oauth_id = ?;", tokenData["oauth_id"]); err != nil {
				newUUID := uuid.New()
				mashboardEmail := newUUID.String() + "@mashboard.app"
				_, err := db.Exec("INSERT INTO users (id, oauth_provider, oauth_id, first_name, last_name, mashboard_email) VALUES (?, ?, ?, ?, ?, ?)",
					newUUID[:], "google", tokenData["user_id"], tokenData["firstName"], tokenData["lastName"], mashboardEmail)
				if err != nil {
					log.Error("Failed to insert user:", err)
					return c.Status(http.StatusInternalServerError).SendString("User creation failed")
				}
			}
		}
		return c.Next()
	}
}

func (h *Handler) GetDB() *sqlx.DB {
	return h.db
}

func (h *Handler) GetOauthConfig() *oauth2.Config {
	return h.oauthConfig
}

func (h *Handler) GenerateJWTToken(oauthID string, fname string, lname string, authCode string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    oauthID,
		"first_name": fname,
		"last_name":  lname,
		"auth_code":  authCode,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":        time.Now().Unix(),                     // Issued at
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (h *Handler) GetUserFromToken(tokenString string) (*models.User, error) {
	token, err := h.ParseTokenString(tokenString)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Debug(token)
	// log.Debug(token)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		oauthID, exists := claims["user_id"].(string)
		if exists {
			db := h.GetDB()
			log.Debug(oauthID)
			var user models.User
			if err = db.Get(&user, "SELECT * from users WHERE oauth_id = ? LIMIT 1 ;", oauthID); err != nil {
				return nil, err
			}
			return &user, nil
		}
	}
	return nil, errors.New("invalid token")

}

func (h *Handler) GetJWTInfo(tokenString string, expired bool) (map[string]string, error) {
	token := &jwt.Token{}
	if expired {
		t, err := h.ParseExpiredTokenString(tokenString)
		if err != nil {
			return nil, err
		}
		token = t

	} else {
		t, err := h.ParseTokenString(tokenString)
		if err != nil {
			return nil, err
		}
		token = t
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		oauthID := ""
		lname := ""
		fname := ""
		authCode := ""
		oauthID, exists := claims["user_id"].(string)
		if !exists {
			return nil, errors.New("unexpected error: oauth not in token")
		}
		fname, exists = claims["first_name"].(string)
		if !exists {
			return nil, errors.New("unexpected error: fname not in token")
		}
		lname, exists = claims["last_name"].(string)
		if !exists {
			return nil, errors.New("unexpected error: lname not in token")
		}

		authCode, exists = claims["auth_code"].(string)
		if !exists {
			return nil, errors.New("unexpected error: auth code not in token")
		}
		return map[string]string{
			"user_id":   oauthID,
			"fname":     fname,
			"lname":     lname,
			"auth_code": authCode,
		}, nil
	}
	return nil, errors.New("unexpected error")

}

// Parses uuid string into uuid and converts it to bytes for storage in SQLite.
func (h *Handler) ParseUUIDString(uuidStr string) ([]byte, error) {
	feedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return nil, err
	}
	return feedUUID[:], nil
}

func (h *Handler) ValidateURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close() // Prevent resource leak

	// Ensure we have a successful response
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return false
	}
	contentType := res.Header.Get("Content-Type")
	validTypes := []string{"application/xml", "text/xml", "application/rss+xml"}

	for _, validType := range validTypes {
		if strings.Contains(contentType, validType) {
			return true
		}
	}
	return false
}

func (h *Handler) QueryRows(items *[]map[string]any, query string, args ...any) error {
	db := h.GetDB()
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	for rows.Next() {
		values := make([]any, len(cols))
		valuePtrs := make([]any, len(cols))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Error(err)
			return err
		}
		item := make(map[string]any)
		for i, colName := range cols {
			item[colName] = values[i]
		}
		*items = append(*items, item)
	}
	return nil
}
