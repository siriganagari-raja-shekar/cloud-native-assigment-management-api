package controllers

import (
	"csye6225-mainproject/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealthzHandler(provider *services.ServiceProvider) func(ctx *gin.Context) {

	return func(context *gin.Context) {
		context.Header("Cache-Control", "no-cache")
		switch context.Request.Method {
		case http.MethodGet:

			if !isPayloadEmpty(context) || len(context.Request.URL.RawQuery) > 0 {
				context.String(http.StatusBadRequest, "")
			}

			connected, _ := provider.MyHealthzStore.Ping()

			if connected {
				context.String(http.StatusOK, "")
			} else {
				context.String(http.StatusServiceUnavailable, "")
			}

		default:
			context.String(http.StatusMethodNotAllowed, "")
		}
	}

}

func isPayloadEmpty(context *gin.Context) bool {
	body := make([]byte, 1)

	n, _ := context.Request.Body.Read(body)

	return n == 0
}
