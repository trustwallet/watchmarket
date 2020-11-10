package endpoint

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/trustwallet/watchmarket/services/controllers"
	"net/http"
)

func GetAlertsHandler(controller controllers.AlertsController) func(c *gin.Context) {
	return func(c *gin.Context) {
		interval := c.Query("interval")
		ar := controllers.AlertsRequest{Interval: interval}
		response, err := controller.HandleAlertsRequest(ar, c.Request.Context())
		if err != nil {
			log.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, nil)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}
