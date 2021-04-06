package middlewares

import (
	"context"
	"errors"

	// "go-graphql-jwt/graph/model"
	// "go-graphql-jwt/graph/utils"
	"net"
	"net/http"
	"strings"

	"github.com/amenabe22/chachata_backend/graph/model"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// UserAuth - user auth middleware structure
type UserAuth struct {
	UserID    string
	Roles     []string
	IPAddress string
	Token     string
}

var userCtxKey = &contextKey{name: "user"}

type contextKey struct {
	name string
}

// UserClaims for jwt
type UserClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

// JwtMiddleware middleware for http server
func JwtMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := TokenFromHTTPRequestgo(r)

			userID := UserIDFromHTTPRequestgo(token)
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)
			userAuth := UserAuth{
				UserID:    userID,
				IPAddress: ip,
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, &userAuth)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// TokenFromHTTPRequestgo - get jwt token from request
func TokenFromHTTPRequestgo(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}

// UserIDFromHTTPRequestgo - get user id from request
func UserIDFromHTTPRequestgo(tokenString string) string {
	token, err := DecodeJwt(tokenString)
	if err != nil {
		return "err"
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		if claims == nil {
			return "err"
		}
		return claims.UserID
	}
	return "Err"
}

// GetAuthFromContext - Gets context
func GetAuthFromContext(ctx context.Context) *UserAuth {
	raw := ctx.Value(userCtxKey)
	if raw == nil {
		return nil
	}

	return raw.(*UserAuth)
}

var issuer = []byte("github/thealmarques")

// DecodeJwt decode jwt
func DecodeJwt(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return issuer, nil
	})
}

// GenerateJwt create jwt
func GenerateJwt(userID string, expiredAt int64) string {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredAt,
			Issuer:    string(issuer),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(issuer)

	return signedToken
}

func GetAuthStat(coredb *gorm.DB, ctx context.Context, errorMessage string) (bool, model.User, error) {
	userAuth := GetAuthFromContext(ctx)
	if userAuth.UserID == "err" {
		return false, model.User{}, errors.New(errorMessage)
	}
	user := model.User{}
	coredb.First(&user, "id = ?", userAuth.UserID)

	return true, user, nil
}
