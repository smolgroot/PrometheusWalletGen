package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Display(err error) {
	if err != nil {
		log.Println(err)
	}
}

func checkFileExists(filePath string) bool {
	_, error := os.Open(filePath) // For read access.
	return error == nil
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

// Coloring the terminal
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"
