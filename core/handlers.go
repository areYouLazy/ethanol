package core

import (
	"net/http"

	"github.com/areYouLazy/ethanol/plugins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func getStatusHandler(c *gin.Context) {
	c.String(http.StatusOK, "ONLINE")
}

func getStatusJSONHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ONLINE",
	})
}

func getWebSocketHandler(c *gin.Context) {
	serveWebSocket(webSocketHub, c.Writer, c.Request)
}

func getAPIV1UISearchHandler(c *gin.Context) {
	// get q=query parameter from request url
	query := c.Query("q")

	// traw error on empty queries
	if query == "" {
		logrus.WithFields(logrus.Fields{
			"query": query,
		}).Error("empty query parameter")

		c.JSON(422, gin.H{
			"error": "query parameter cannot be empty",
			"query": "",
		})
	} else {
		logrus.WithFields(logrus.Fields{
			"query": query,
		}).Debug("executing bulk search")

		// execute bulk search
		results, err := plugins.BulkSearch(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.Status(http.StatusOK)
		c.Writer.Write([]byte(renderHTMLResultsFromJSON(results)))
	}
}

func getAPIV1SearchHandler(c *gin.Context) {
	// get q=query parameter from request url
	query := c.Query("q")

	// traw error on empty queries
	if query == "" {
		logrus.WithFields(logrus.Fields{
			"query": query,
		}).Error("empty query parameter")

		c.JSON(422, gin.H{
			"error": "query parameter cannot be empty",
			"query": "",
		})
	} else {
		logrus.WithFields(logrus.Fields{
			"query": query,
		}).Debug("executing bulk search")

		// execute bulk search
		results, err := plugins.BulkSearch(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(http.StatusOK, results)
	}
}

func uixHandler(c *gin.Context) {
	c.HTML(http.StatusOK, c.Param("template"), nil)
}
