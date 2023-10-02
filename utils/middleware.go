package utils

import (
	"csye6225-mainproject/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

func UserExtractor(context *gin.Context, provider services.ServiceProvider) {
	authorization := context.Request.Header.Get("Authorization")

	authorization, found := strings.CutPrefix(authorization, "Basic ")

	if !found {
		context.JSON(http.StatusUnauthorized, struct {
			Error string
		}{"The authorization header is not in correct format. Please check again"})
	}

	emailAndPassword := strings.Split(authorization, ":")

	user, err := provider.MyAccountStore.GetOneByEmail(emailAndPassword[0])

	if err != nil {
		context.JSON(http.StatusUnauthorized, struct {
			Error string
		}{"User not found in the database. Please check your email and try again"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(emailAndPassword[1]))

	if err != nil {
		context.JSON(http.StatusUnauthorized, struct {
			Error string
		}{"Password is wrong. Please check and try again``"})
	}

	context.Set("currenUserAccount", user)

	context.Next()
}
