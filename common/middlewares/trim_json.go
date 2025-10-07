package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func TrimJSONMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ContentType() == "application/json" {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil && len(bodyBytes) > 0 {
				var data map[string]interface{}
				if err := json.Unmarshal(bodyBytes, &data); err == nil {
					for k, v := range data {
						if str, ok := v.(string); ok {
							data[k] = strings.TrimSpace(str)
						}
					}
					newBody, _ := json.Marshal(data)
					c.Request.Body = io.NopCloser(bytes.NewReader(newBody))
				}
			}
		}
		c.Next()
	}
}
