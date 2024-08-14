package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"prometheus/pkg/generator"
	"prometheus/pkg/ui"
	"prometheus/pkg/utils"
	"strconv"
	"syscall"

	_ "github.com/charmbracelet/bubbletea"
	"github.com/gin-contrib/cors"

	// _ "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// _ "github.com/gin-gonic/gin"
	_ "github.com/spf13/cobra"
)

func main() {
	// _ = wallet.InitPass()
	// log.Println(wallet.Pass)

	// Catch SIGINT when Ctrl+C is pressed, and exit gracefully
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Exiting prometheus...")
		os.Exit(1)
	}()

	// gin.SetMode(gin.ReleaseMoœde)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.SetTrustedProxies(nil)
	router.Use(cors.Default())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// React SPA Middleware (must be last middleware declared)
	// router.Use(ui.NewHandler().ServeSPA)
	// log.Printf("UI available at http://prometheus.localhost:%d/\n", 8081)

	// API Router
	api := router.Group("/api/v0")
	api.GET("/", utils.Ping)
	api.GET("/ping", utils.Ping)
	api.POST("/generate", generator.Generate)
	// api.GET("/balance", wallet.GetBalance)
	// api.GET("/checkAuth", wallet.IsUnlocked)
	// api.GET("/channels", generator.GetChannels)
	// api.GET("/newChannel/:room", generator.CreateChannel)
	// api.GET("/trust/:addr", generator.TrustAddress)
	// api.GET("/channel/:room/messages", generator.GetMessages)
	// api.GET("/channel/:room/clear", generator.ClearMessages)
	// api.POST("/channel/:room/message", generator.SendMessage)
	// api.POST("/channel/:room/gps", generator.SendGPSCoordinates)
	// api.POST("/channel/:room/picture", generator.SendPicture)
	// api.POST("/unlock", wallet.Unlock)
	// api.GET("/quit", wallet.Quit)

	// React SPA Middleware
	// It must be last middleware declared
	router.Use(ui.NewHandler().ServeSPA)
	uiUrl := fmt.Sprintf("http://localhost:%d/", 8081)
	fmt.Println(uiUrl)
	log.Printf(utils.Green+"Prometheus UI available at %v"+utils.Reset, uiUrl)
	ui.OpenInBrowser(uiUrl)

	// Run with HTTP
	err := router.Run("0.0.0.0:" + strconv.FormatUint(uint64(8081), 10))
	utils.Check(err)
}
