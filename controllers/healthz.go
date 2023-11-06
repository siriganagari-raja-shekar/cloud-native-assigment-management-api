package controllers

import (
	"csye6225-mainproject/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetHealthzHandler(provider *services.ServiceProvider) func(ctx *gin.Context) {

	return func(context *gin.Context) {

		provider.MyStatsStore.Client.Incr("api.requests.assignments.healthz", 1)

		context.Header("Cache-Control", "no-cache")
		switch context.Request.Method {
		case http.MethodGet:

			if context.Request.ContentLength > 0 || len(context.Request.URL.RawQuery) > 0 {
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
