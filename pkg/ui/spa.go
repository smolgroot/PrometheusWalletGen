package ui

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"prometheus/pkg/utils"

	"github.com/gin-gonic/gin"
)

//go:embed dist
var Content embed.FS

const rootPath = "/"

// Handler represents a HTTP handler for serving a Single Page Application (SPA).
// It encapsulates the necessary components for serving static files and handling requests.
type Handler struct {
	fileSystem       http.FileSystem
	fileServer       http.Handler
	once             sync.Once
	serveSPAFuncName string
}

// openInBrowser opens the specified URL in the default web browser.
func OpenInBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// NewHandler returns a new instance of the Handler struct, which handles serving static files for the web application.
func NewHandler() *Handler {
	// Open the dist folder.
	authorizeSite, err := fs.Sub(Content, "dist")
	utils.Check(err)

	fileSystem := http.FS(authorizeSite)
	return &Handler{
		fileSystem: fileSystem,
		fileServer: http.StripPrefix(rootPath, http.FileServer(fileSystem)),
	}
}

// Register registers the ServeSPA middleware with the provided Gin engine.
// This middleware is responsible for serving the Single Page Application (SPA)
// for all routes that are not matched by other handlers.
func (h *Handler) Register(e *gin.Engine) {
	e.Use(h.ServeSPA)
}

// ServeSPA serves the Single Page Application (SPA) for the web3vault application.
// It handles the HTTP GET and HEAD requests and serves the appropriate files from the file system.
// If the requested file does not exist, it redirects the request to the root path.
// This function is called by the Gin framework's router.
//
// Parameters:
// - c: The Gin context object representing the current HTTP request and response.
//
// Returns: None
func (h *Handler) ServeSPA(c *gin.Context) {
	h.once.Do(func() {
		h.serveSPAFuncName = runtime.FuncForPC(reflect.ValueOf(h.ServeSPA).Pointer()).Name()
	})

	if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
		c.Next()
		return
	}

	names := c.HandlerNames()
	if len(names) != 0 && names[len(names)-1] != h.serveSPAFuncName {
		c.Next()
		return
	}

	filePath := c.Request.URL.Path

	if filePath == rootPath {
		h.fileServer.ServeHTTP(c.Writer, c.Request)
		return
	}

	f, err := h.fileSystem.Open(strings.TrimPrefix(path.Clean(filePath), "/"))
	if err == nil {
		_ = f.Close()
	}

	var req = c.Request
	if os.IsNotExist(err) {
		req = c.Request.Clone(c.Request.Context())
		req.URL.Path = rootPath
	}

	h.fileServer.ServeHTTP(c.Writer, req)
}
