package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"russian-learns-english/internal/api/service"
)

type wordHandler struct {
	wordService service.WordService
}

func SetupWordHandler(router *gin.RouterGroup, wordService service.WordService) {
	wh := wordHandler{wordService}
	router.GET("", wh.GetAllWords)
	router.GET("/random", wh.RandomWord)
	router.POST("/check-translation", wh.CheckWordTranslation)
	router.POST("/upload", wh.UploadWords)
}

func (wh *wordHandler) RandomWord(c *gin.Context) {
	c.JSON(http.StatusOK, wh.wordService.GetRandomWord())
}

func (wh *wordHandler) GetAllWords(c *gin.Context) {
	c.JSON(http.StatusOK, wh.wordService.GetWordList())
}

func (wh *wordHandler) UploadWords(c *gin.Context) {
	h, err := c.FormFile("words")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if h.Header.Get("Content-Type") != "text/plain" {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "file must be in text/plain format"})
		return
	}

	f, err := h.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()

	totalLoaded, err := wh.wordService.LoadWordList(f)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Words loaded %d", totalLoaded)})
}

type checkWordTranslationRequest struct {
	ID           service.WordID `json:"id"`
	Translations []string       `json:"translations"`
}

func (wh *wordHandler) CheckWordTranslation(c *gin.Context) {
	req := checkWordTranslationRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := wh.wordService.CheckWordTranslation(req.ID, req.Translations)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": res})
}
