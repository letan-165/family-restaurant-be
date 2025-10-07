package utils

import (
	"myapp/common/errors_code"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JSONError(c *gin.Context, err error) {
	appErr := errors_code.Wrap(err)
	c.JSON(appErr.Status, gin.H{
		"code":    appErr.Code,
		"message": appErr.Message,
	})
}

func JSONData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
