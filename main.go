package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURLLength = 6
)

type ShortURL struct {
	ID       string `json:"id"`
	Original string `json:"original"`
}

var urlMap = make(map[string]ShortURL)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	router := gin.Default()

	router.POST("/shorten", shortenURL)
	router.GET("/:shortURL", redirectURL)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Error starting the Server:", err)
	}
}

func generateShortURL() string {
	b := make([]byte, shortURLLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func shortenURL(c *gin.Context) {
	var requestBody struct {
		Original string `json:"original"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	shortURL := generateShortURL()
	urlMap[shortURL] = ShortURL{ID: shortURL, Original: requestBody.Original}

	c.JSON(http.StatusOK, gin.H{"short.url": shortURL})
}

func redirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	url, exists := urlMap[shortURL]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
	}

	c.Redirect(http.StatusFound, url.Original)
}
