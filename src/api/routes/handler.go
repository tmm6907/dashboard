package routes

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/tmm6907/dashboard/utils"
	"golang.org/x/oauth2"
)

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

func (h *Handler) GetDB() *sqlx.DB {
	return h.db
}

func (h *Handler) GetOauthConfig() *oauth2.Config {
	return h.oauthConfig
}

func (h *Handler) GenerateJWTToken(oauthID string, fname string, lname string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":    oauthID,
		"first_name": fname,
		"last_name":  lname,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
		"iat":        time.Now().Unix(),                     // Issued at
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func (h *Handler) GetUserIDFromToken(tokenString string) (utils.UUID, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		oauthID, exists := claims["user_id"].(string)
		if exists {
			db := h.GetDB()
			var userID utils.UUID
			if err = db.Get(&userID, "SELECT id from users WHERE oauth_id = ? ;", oauthID); err != nil {
				return nil, err
			}
			return userID, nil
		}

	}
	return nil, errors.New("invalid token")

}

func (h *Handler) GetOauthInfoFromToken(tokenString string) (string, string, string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if oauthID, exists := claims["user_id"].(string); exists {
			fname := claims["first_name"].(string)
			lname := claims["last_name"].(string)
			return oauthID, fname, lname, nil
		}
	}
	return "", "", "", errors.New("invalid token")

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
