package matches

import (
	"example/matches/order"
	"example/util"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func Test_OrderHandleCase1(t *testing.T) {
	buyOrderHandle := newOrderHandle(order.OrderType_Buy)
	sellOrderHandle := newOrderHandle(order.OrderType_Sell)

	testCase1(buyOrderHandle, sellOrderHandle)
}
func Test_OrderHandleCase2(t *testing.T) {
	rand.Seed(12345)
	buyOrderHandle := newOrderHandle(order.OrderType_Buy)
	sellOrderHandle := newOrderHandle(order.OrderType_Sell)

	testCase2(buyOrderHandle, sellOrderHandle)
}

// TestCase 1
// 賣 1.5 Market
// 賣 1.5 Limit
// 買 1.5 Limit
// 買 1.5 Market
// 買 1.5 Limit
func testCase1(buyOrderHandle, sellOrderHandle *orderHandle) {
	reqOrder := &order.OrderInfo{
		OrderType:       order.OrderType_Sell,
		Price:           1.5,
		TransactionType: order.TransactionType_MarketPrice,
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        100,
	}

	settlementOrders := buyOrderHandle.MatchOrder(reqOrder)
	if len(settlementOrders) == 0 {
		sellOrderHandle.InsertOrder(reqOrder)
	} else {
		Quantity := reqOrder.Quantity
		for _, endOrder := range settlementOrders {
			Quantity -= endOrder.Quantity
		}

		if Quantity > 0 {
			reqOrder.Quantity = Quantity
			sellOrderHandle.InsertOrder(reqOrder)
		}
	}
	time.Sleep(time.Microsecond)

	reqOrder = &order.OrderInfo{
		OrderType:       order.OrderType_Sell,
		Price:           1.5,
		TransactionType: order.TransactionType_LimitPrice,
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        100,
	}

	settlementOrders = buyOrderHandle.MatchOrder(reqOrder)
	if len(settlementOrders) == 0 {
		sellOrderHandle.InsertOrder(reqOrder)
	} else {
		Quantity := reqOrder.Quantity
		for _, endOrder := range settlementOrders {
			Quantity -= endOrder.Quantity
		}

		if Quantity > 0 {
			reqOrder.Quantity = Quantity
			sellOrderHandle.InsertOrder(reqOrder)
		}
	}

	reqOrder = &order.OrderInfo{
		OrderType:       order.OrderType_Buy,
		Price:           1.5,
		TransactionType: order.TransactionType_LimitPrice,
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        150,
	}
	settlementOrders = sellOrderHandle.MatchOrder(reqOrder)
	if len(settlementOrders) == 0 {
		buyOrderHandle.InsertOrder(reqOrder)
	} else {
		Quantity := reqOrder.Quantity
		for _, endOrder := range settlementOrders {
			Quantity -= endOrder.Quantity
		}

		if Quantity > 0 {
			reqOrder.Quantity = Quantity
			sellOrderHandle.InsertOrder(reqOrder)
		}
	}

	reqOrder = &order.OrderInfo{
		OrderType:       order.OrderType_Buy,
		Price:           1.5,
		TransactionType: order.TransactionType_MarketPrice,
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        20,
	}
	settlementOrders = sellOrderHandle.MatchOrder(reqOrder)
	if len(settlementOrders) == 0 {
		buyOrderHandle.InsertOrder(reqOrder)
	} else {
		Quantity := reqOrder.Quantity
		for _, endOrder := range settlementOrders {
			Quantity -= endOrder.Quantity
		}

		if Quantity > 0 {
			reqOrder.Quantity = Quantity
			buyOrderHandle.InsertOrder(reqOrder)
		}
	}

	reqOrder = &order.OrderInfo{
		OrderType:       order.OrderType_Buy,
		Price:           1.5,
		TransactionType: order.TransactionType_LimitPrice,
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        300,
	}
	settlementOrders = sellOrderHandle.MatchOrder(reqOrder)
	if len(settlementOrders) == 0 {
		buyOrderHandle.InsertOrder(reqOrder)
	} else {
		Quantity := reqOrder.Quantity
		for _, endOrder := range settlementOrders {
			Quantity -= endOrder.Quantity
		}

		if Quantity > 0 {
			reqOrder.Quantity = Quantity
			buyOrderHandle.InsertOrder(reqOrder)
		}
	}
}

// TestCase 2
// 亂數測試
func testCase2(buyOrderHandle, sellOrderHandle *orderHandle) {
	midPrice := 155
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

		if reqOrder.OrderType == order.OrderType_Sell {
			fmt.Println("reqOrder sell:", reqOrder.TransactionType, reqOrder.Price)
			settlementOrders := buyOrderHandle.MatchOrder(&reqOrder)
			if settlementOrders == nil {
				sellOrderHandle.InsertOrder(&reqOrder)
			} else {
				prices := []float64{}
				Quantity := reqOrder.Quantity
				for _, endOrder := range settlementOrders {
					Quantity -= endOrder.Quantity
					prices = append(prices, endOrder.Price)
				}

				if Quantity > 0 {
					reqOrder.Quantity = Quantity
					sellOrderHandle.InsertOrder(&reqOrder)
				}
				fmt.Println("Match: sell", reqOrder.TransactionType, reqOrder.Price, prices)
			}
			fmt.Println("SellPrice:", sellOrderHandle.priceSlice, "noprice:", len(sellOrderHandle.marketPriceOrder))
			fmt.Println("BuyPrice:", buyOrderHandle.priceSlice, "noprice:", len(buyOrderHandle.marketPriceOrder))
		} else {
			fmt.Println("reqOrder buy:", reqOrder.TransactionType, reqOrder.Price)
			settlementOrders := sellOrderHandle.MatchOrder(&reqOrder)
			if len(settlementOrders) == 0 {
				buyOrderHandle.InsertOrder(&reqOrder)
			} else {
				prices := []float64{}
				Quantity := reqOrder.Quantity
				for _, endOrder := range settlementOrders {
					Quantity -= endOrder.Quantity
					prices = append(prices, endOrder.Price)
				}

				if Quantity > 0 {
					reqOrder.Quantity = Quantity
					buyOrderHandle.InsertOrder(&reqOrder)
				}
				fmt.Println("Match: buy", reqOrder.TransactionType, reqOrder.Price, prices)
			}
			fmt.Println("SellPrice:", sellOrderHandle.priceSlice, "noprice:", len(sellOrderHandle.marketPriceOrder))
			fmt.Println("BuyPrice:", buyOrderHandle.priceSlice, "noprice:", len(buyOrderHandle.marketPriceOrder))
		}
		fmt.Println("----------------------------------")
	}
}
