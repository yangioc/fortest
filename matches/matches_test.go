package matches

import (
	"context"
	"example/matches/order"
	"example/util"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestMatches(t *testing.T) {
	midPrice := 150
	engine := New(context.TODO())

	for i := 0; i < 200000; i++ {
		reqOrder := order.OrderInfo{
			OrderType:       uint8(rand.Intn(2) + 1),
			TransactionType: 1, // uint8(rand.Intn(2) + 1),
			OrderStartTime:  time.Now().UnixNano(),
			Quantity:        int64(rand.Intn(100)),
		}

		if reqOrder.TransactionType == order.TransactionType_LimitPrice {
			if reqOrder.OrderType == order.OrderType_Sell {
				reqOrder.Price = util.RoundFloor(float64(midPrice-10)+float64(rand.Intn(1000))/100, 2)
			} else {
				reqOrder.Price = util.RoundFloor(float64(midPrice+10)-float64(rand.Intn(1000))/100, 2)
			}
		}

		go engine.NewOrder(&reqOrder)
	}

	fmt.Println("---------------------------------\n OrderEnd\n---------------------------------")
	select {}
}
