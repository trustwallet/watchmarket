package chartscontroller

type (
	chartsNormalizedRequest struct {
		Coin            uint
		Token, Currency string
		TimeStart       int64
		MaxItems        int
	}
)
