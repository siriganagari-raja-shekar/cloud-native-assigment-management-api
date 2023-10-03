package utils

import (
	"csye6225-mainproject/services"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func UserExtractor(provider *services.ServiceProvider) func(context *gin.Context) {

	return func(context *gin.Context) {
		authorization := context.Request.Header.Get("Authorization")

		authorization, found := strings.CutPrefix(authorization, "Basic ")

		if !found {
			context.AbortWithStatusJSON(http.StatusUnauthorized, struct {
				Error string
			}{"The authorization header is not in correct format. Please check again"})
		}

		decodedString, _ := base64.StdEncoding.DecodeString(authorization)

		emailAndPasswordJoined := string(decodedString)

		emailAndPassword := strings.Split(emailAndPasswordJoined, ":")

		user, err := provider.MyAccountStore.GetOneByEmail(emailAndPassword[0])

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, struct {
				Error string
			}{"User not found in the database. Please check your email and try again"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(emailAndPassword[1]))

		if err != nil {
			context.AbortWithStatusJSON(http.StatusUnauthorized, struct {
				Error string
			}{"Password is wrong. Please check and try again``"})
			return
		}

		context.Set("currentUserAccount", user)

		context.Next()
	}

}

func InvalidHandler(context *gin.Context) {
	context.Header("Cache-Control", "no-cache")
	context.String(http.StatusNotFound, "")
}
