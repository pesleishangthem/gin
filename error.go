package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/pesleishangthem/apperror"
	"github.com/rs/zerolog"
)

func WriteError(log zerolog.Logger, c *gin.Context, err error) {
	log.Err(err)
	if appErr, ok := err.(*apperror.AppError); ok {
		switch appErr.Code {
		case apperror.ErrValidation:
			c.JSON(400, gin.H{"error": appErr.Message})

		case apperror.ErrNotFound:
			c.JSON(404, gin.H{"error": appErr.Message})

		case apperror.ErrConflict:
			c.JSON(409, gin.H{"error": appErr.Message})

		default:
			c.JSON(500, gin.H{"error": "internal error"})
		}
		return
	}
	c.JSON(500, gin.H{"error": "internal error"})
}
