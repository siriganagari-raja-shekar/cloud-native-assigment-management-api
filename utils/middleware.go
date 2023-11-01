package utils

import (
	"csye6225-mainproject/services"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func UserExtractor(provider *services.ServiceProvider) func(context *gin.Context) {

	return func(context *gin.Context) {

		connected, _ := provider.MyHealthzStore.Ping()

		if !connected {
			slog.Info("Health check before request: Cannot connect to database, aborting request processing")
			context.AbortWithStatus(http.StatusServiceUnavailable)
			return
		} else {
			slog.Info("Health check before request: Database is online, proceeding with request")
			provider.PopulateDBInServices()
		}

		authorization := context.Request.Header.Get("Authorization")

		authorization, found := strings.CutPrefix(authorization, "Basic ")

		if !found {
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "The authorization header is not in correct format. Please check again",
			})
			return
		}

		decodedString, _ := base64.StdEncoding.DecodeString(authorization)

		emailAndPasswordJoined := string(decodedString)

		emailAndPassword := strings.Split(emailAndPasswordJoined, ":")

		user, err := provider.MyAccountStore.GetOneByEmail(emailAndPassword[0])

		if err != nil {
			slog.Warn(fmt.Sprintf("Unauthorized user login with email %s", emailAndPassword[0]))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found in the database. Please check your email and try again",
			})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(emailAndPassword[1]))

		if err != nil {
			slog.Info(fmt.Sprintf("User with email %s logged in with wrong password", user.Email))
			context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Password is wrong. Please check and try again",
			})
			return
		}

		if len(context.Request.URL.RawQuery) > 0 {
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Please don't specify any query parameters in the request",
			})
			return
		}

		context.Set("currentUserAccount", user)

		context.Next()
	}

}

func StatsLogger(provider *services.ServiceProvider) func(context *gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests", 1)

		start := time.Now()

		context.Next()

		provider.MyStatsStore.Client.PrecisionTiming("api.request.processing.time", time.Since(start))

	}

}

func InvalidHandler(provider *services.ServiceProvider) func(context *gin.Context) {

	return func(context *gin.Context) {
		provider.MyStatsStore.Client.Incr("api.requests.invalid.count", 1)
		context.Header("Cache-Control", "no-cache")
		context.String(http.StatusNotFound, "")
	}
}

func MethodNotAllowedHandler(context *gin.Context) {
	context.Header("Cache-Control", "no-cache")
	context.String(http.StatusMethodNotAllowed, "")
}
