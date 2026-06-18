package proxy

import (
	"io"
	"log"
	"net/http"
	"sociul-api-gateway/internal/config"

	"github.com/gin-gonic/gin"
)

type Proxy struct {
	client         *http.Client
	authServiceUrl string
}

func NewProxy(client *http.Client, cfg *config.Config) *Proxy {
	return &Proxy{client: client, authServiceUrl: cfg.AuthServiceURL}
}

// Hop-by-hop headers - they are made only for one http call (upstream -> api gateway)
// Need to strip them before passing upstream response to client
var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"TE",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

// Check if a header is hop by hop
func isHopbyHop(header string) bool {
	for _, hopHeader := range hopByHopHeaders {
		if hopHeader == header {
			return true
		}
	}
	return false
}

// Create url for auth routes and serve request
func (p *Proxy) Auth(c *gin.Context) {
	path := c.Param("path")
	url := p.authServiceUrl + "/auth" + path
	p.do(c, url)
}

func (p *Proxy) do(c *gin.Context, url string) {
	// Create request
	request, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, url, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Send to internal service
	resp, err := p.client.Do(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	defer resp.Body.Close()

	// Copy headers
	for key, values := range resp.Header {
		if isHopbyHop(key) { // skip hop-by-hop headers
			continue
		}
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code and stream response
	c.Status(resp.StatusCode)

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		log.Println("Error streaming response - ", err)
		c.Status(500)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
