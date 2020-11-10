package alertscontroller

import (
	"context"
	"github.com/trustwallet/watchmarket/config"
	"github.com/trustwallet/watchmarket/db"
	"github.com/trustwallet/watchmarket/db/models"
	"github.com/trustwallet/watchmarket/services/controllers"
)

type Controller struct {
	database      db.Instance
	configuration config.Configuration
}

func NewController(
	database db.Instance,
	configuration config.Configuration,
) Controller {
	return Controller{
		database,
		configuration,
	}
}

func (c Controller) HandleAlertsRequest(ar controllers.AlertsRequest, ctx context.Context) (controllers.AlertsResponse, error) {
	alerts, err := c.database.GetAlertsByIntervalWithDifference(
		models.Interval(ar.Interval), c.configuration.RestAPI.Alerts.PriceDifference, ctx)
	if err != nil {
		return controllers.AlertsResponse{}, err
	}
	return normalizeAlerts(alerts), nil
}

func normalizeAlerts(alerts []models.Alert) controllers.AlertsResponse {
	var result controllers.AlertsResponse
	details := make(map[string]controllers.AlertsDetails)
	for _, a := range alerts {
		d := controllers.AlertsDetails{
			PriceDifference: a.Difference,
			UpdatedAt:       a.UpdatedAt.Unix(),
		}
		details[a.AssetID] = d
	}
	result.Assets = details
	return result
}
