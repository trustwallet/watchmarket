package api

import (
	"github.com/gin-gonic/gin"
	"github.com/trustwallet/blockatlas/pkg/ginutils"
)

func GetRoot(c *gin.Context) {
	ginutils.RenderSuccess(c,
		`Welcome to the Block Atlas API!

Don't know how you landed here?
Visit https://trustwallet.com to get back to the main page.

If you know what you're doing:
 - Visit /v1/ to list platforms
 - Source: https://github.com/trustwallet/watchmarket
 - Any questions? https://t.me/walletcore
`)
}
