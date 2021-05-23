package server

import (
	"github.com/gin-gonic/gin"
	"russian-learns-english/internal/api/delivery/http/handler"
	"russian-learns-english/internal/api/service"
)

func New(wordService service.WordService) *gin.Engine {
	router := gin.Default()
	handler.SetupWordHandler(router.Group("/words"), wordService)
	return router
}
