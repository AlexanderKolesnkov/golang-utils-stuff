package errs

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func BindError(c *gin.Context, logger *zap.Logger, err error) {
	logger.Error("error", zap.Error(err))
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func SendBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}
