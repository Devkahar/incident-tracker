package router

import (
	"incident-tracker/config"
	"incident-tracker/errors"
	"runtime"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

type requestHandler[R any] func(applicationContext *config.ApplicationContext, c *gin.Context) (*R, error)

func HandleRequest[R any](applicationContext *config.ApplicationContext, handler requestHandler[R]) gin.HandlerFunc {
	log := applicationContext.Logger
	return func(c *gin.Context) {
		var response *R
		var handlerError error
		// Panic Recovery
		defer func() {
			if exception := recover(); exception != nil {
				stack := debug.Stack()
				log.Sugar().Errorf("Stack trace for exception : %v \n %v", exception, string(stack))
				if _, file, line, ok := runtime.Caller(1); ok {
					log.Sugar().Errorf("Recovered from panic in file %s at line %d: %v\n", file, line, exception)
				} else {
					log.Error("Recovered from panic but couldn't retrieve file name and line number")
				}
				c.AbortWithStatusJSON(500, exception)
				return
			}
			if handlerError != nil {
				apiError, ok := handlerError.(errors.APIError)
				if !ok {
					log.Sugar().Errorf("Api Request error %v", handlerError)
					c.JSON(500, handlerError)
					return
				}
				c.JSON(apiError.StatusCode(), apiError.Json())
				return
			}
			c.JSON(200, map[string]any{
				"data": *response,
			})
		}()
		response, handlerError = handler(applicationContext, c)
	}
}
