package middlewares

import (
	"time"

	"github.com/amenabe22/chachata_backend/graph/model"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var identityKey = "id"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Middleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				ID: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			println("We here bruh")
			// userID := loginVals.Username
			// password := loginVals.Password

			// if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			// 	return &model.User{
			// 		UserName:  userID,
			// 		LastName:  "Bo-Yi",
			// 		FirstName: "Wu",
			// 	}, nil
			// }

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			// if v, ok := data.(*model.User); ok && v.ID == "admin" {
			// 	return true
			// }
			println("This is the authorizor!!!")
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	if err != nil {
		return nil, err
	}
	return authMiddleware, nil
}

// package middlewares

// import (
// 	"context"
// 	"net/http"

// 	"github.com/amenabe22/chachata_backend/graph/jwt"
// 	"github.com/amenabe22/chachata_backend/graph/model"
// 	"github.com/amenabe22/chachata_backend/graph/setup"
// )

// var userCtxKey = &contextKey{"user"}

// type contextKey struct {
// 	name string
// }

// var db = setup.SetupModels()

// func Middleware() func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			header := r.Header.Get("Authorization")

// 			// Allow unauthenticated users in
// 			if header == "" {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			//validate jwt token
// 			tokenStr := header
// 			usrID, err := jwt.ParseToken(tokenStr)
// 			if err != nil {
// 				http.Error(w, "Invalid token", http.StatusForbidden)
// 				return
// 			}

// 			// create user and check if user exists in db
// 			user := &model.User{ID: usrID}
// 			db.Find(&user)
// 			// println(user.ID, "HERE")
// 			// println(uint(finID), "TST")
// 			// id, err := model.GetUserIdByUsername(username)
// 			if err != nil {
// 				next.ServeHTTP(w, r)
// 				return
// 			}
// 			// put it in context
// 			user.ID = usrID
// 			ctx := context.WithValue(r.Context(), userCtxKey, &user)
// 			// println(user.Name)
// 			r = r.WithContext(ctx)
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// // ForContext finds the user from the context. REQUIRES Middleware to have run.

// func ForContext(ctx context.Context) model.User {
// 	// println(ctx.Err().Error())
// 	// TODO: FIX server is crashing when request is sent with no header
// 	raw, _ := ctx.Value(userCtxKey).(*model.User)
// 	// println(&raw, "here")
// 	return *raw
// }
