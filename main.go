package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func main() {
	// Chat Completions URL to proxy
	proxyURL := os.Getenv("CHAT_URL")
	if proxyURL == "" {
		panic("CHAT_URL env var is missing")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Initialize Sentry
	dsn := os.Getenv("SENTRY_DSN")
	if dsn != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn: dsn,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Sentry initialization failed: %v\n", err)
		}
	}

	r := gin.Default()

	// Health check endpoint
	r.GET("/v1/healthcheck", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Proxy endpoint
	r.POST("/v1/chat/completions", func(c *gin.Context) {
		proxyReq, err := http.NewRequest(http.MethodPost, proxyURL, c.Request.Body)
		if err != nil {
			if dsn != "" {
				sentry.CaptureException(err)
			}
			c.String(http.StatusInternalServerError, "Error creating request")
			return
		}

		proxyReq.Header = c.Request.Header

		client := &http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			if dsn != "" {
				sentry.CaptureException(err)
			}
			c.String(http.StatusInternalServerError, "Error sending request")
			return
		}
		defer resp.Body.Close()

		// Set headers for SSE
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")

		// Stream the response body directly to the client
		io.Copy(c.Writer, resp.Body)
	})

	// Start the server
	r.Run(":" + port)
}
