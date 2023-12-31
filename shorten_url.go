package main

import (
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "Home",
		},
	)
}

func shorturl(c *gin.Context) {
	var url URL
	_ = c.ShouldBind(&url)

	id, _ := (*urlDAO).save(url)
	shortURL := idToShortURL(id, chars)

	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		fmt.Println("HOSTNAME environment variable not set")
		return
	}
	protocol := "http"

	// Construct your URL using the retrieved hostname
	domain := net.JoinHostPort(hostname, serverPort)

	surlLink := fmt.Sprintf("%s://%s/u/%s", protocol, domain, shortURL)

	c.HTML(
		http.StatusOK,
		"url_shorten_summary.html",
		gin.H{
			"title":     "Home",
			"url":       url.URL,
			"short_url": shortURL,
			"domain":    domain,
			"surl_link": template.URL(surlLink),
		},
	)
}

func changeLink(c *gin.Context) {
	var url URLChange
	_ = c.ShouldBind(&url)

	URLID := shortURLToID(url.ShortURL, chars)

	newURL := URL{
		URL: url.NewURL,
	}

	_, err := (*urlDAO).update(URLID, newURL)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.HTML(
		http.StatusOK,
		"surl_linkchanged.html",
		gin.H{
			"title":     "surl - link changed",
			"from_link": url.ShortURL,
			"to_link":   url.NewURL,
		},
	)
}

func redirectShortURL(c *gin.Context) {
	shortURLParam := c.Param("url")
	id := shortURLToID(shortURLParam, chars)

	urlFromDB, err := (*urlDAO).findByID(id)
	if err != nil {
		c.HTML(
			http.StatusInternalServerError,
			"error5xx.html",
			gin.H{
				"title":             "Error",
				"error_description": fmt.Sprintf(`Error redirecting to: %s`, shortURLParam),
			},
		)
	} else {
		c.Redirect(http.StatusMovedPermanently, urlFromDB.URL)
	}
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present.
func render(c *gin.Context, data gin.H, templateName string) {
	loggedInInterface, _ := c.Get("is_logged_in")
	data["is_logged_in"] = loggedInInterface.(bool)

	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

func viewURLs(c *gin.Context) {
	urls, err := (*urlDAO).findAllByUser()
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, urls)
}
