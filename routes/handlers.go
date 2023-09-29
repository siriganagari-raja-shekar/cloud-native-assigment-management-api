package routes

import (
	"csye6225-mainproject/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func invalidHandler(context *gin.Context) {
	context.Header("Cache-Control", "no-cache")
	context.String(http.StatusNotFound, "")
}

func createHealthzHandler(dbHelper db.DatabaseHelper) func(ctx *gin.Context) {

	return func(context *gin.Context) {
		context.Header("Cache-Control", "no-cache")
		switch context.Request.Method {
		case http.MethodGet:

			if !isPayloadEmpty(context) || len(context.Request.URL.RawQuery) > 0 {
				context.String(http.StatusBadRequest, "")
			}

			dbHelper.OpenDBConnection(db.CreateDialectorFromEnv(), db.CreateDBConfig())

			if dbHelper.GetDBConnection() != nil {
				err := dbHelper.CloseDBConnection()
				if err != nil {
					context.String(http.StatusInternalServerError, "")
				} else {
					context.String(http.StatusOK, "")
				}
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
