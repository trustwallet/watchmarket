package postgres

//import (
//	"github.com/trustwallet/watchmarket/db/models"
//)
//
//func (i *Instance) AddSubscription(subscriptionID, coinID uint, token string, price float64, condition models.Condition) error {
//	sub := models.Subscription{
//		Coin:           coinID,
//		Token:          token,
//		Condition:      condition,
//		SubscriptionId: subscriptionID,
//		Price:          price,
//	}
//
//	err := i.Gorm.Set("gorm:insert_option", "ON CONFLICT (subscription_id) DO NOTHING").Create(&sub).Error
//	if err != nil {
//		return err
//	}
//	return nil
//}
