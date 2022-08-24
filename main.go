package main

import (
	"context"
	"example/matches"
	"example/matches/order"
	"example/util"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	st := time.Now()
	midPrice := 100
	ctx := context.TODO()
	handle := matches.New(ctx)

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
		handle.NewOrder(&reqOrder)
	}

	fmt.Println("-----------------Req send finish:", time.Since(st))
	select {}
}

func mai2() {
	ctx := context.TODO()
	handle := matches.New(ctx)

	reqOrder := order.OrderInfo{
		OrderType:       uint8(rand.Intn(2) + 1),
		TransactionType: uint8(rand.Intn(2) + 1),
		OrderStartTime:  time.Now().UnixNano(),
		Quantity:        int64(rand.Intn(100)),
		Price:           float64(rand.Intn(1000)),
	}

	handle.NewOrder(&reqOrder)
}
