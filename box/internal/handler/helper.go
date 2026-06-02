package handler

import "github.com/gin-gonic/gin"

func getUserId(c *gin.Context) string {
	v, _ := c.Get("user_id")
	return v.(string)
}
