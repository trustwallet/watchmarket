package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
)

func GetRoot(c *gin.Context) {
	ginutils.RenderSuccess(c,
		`Watchmarket API`)
}
