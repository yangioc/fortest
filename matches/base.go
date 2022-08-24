package matches

import (
	"context"
	"example/matches/order"
)

const (
	engineEvent_None       = 0
	engineEvent_NewOrder   = 1
	engineEvent_CheckOrder = 2
)

type IMatches interface {
	NewOrder(order *order.OrderInfo) // 發起新訂單
}

type engineFlow struct {
	EventCode uint8
	Order     *order.OrderInfo
}

type engine struct {
	ctx        context.Context
	sellHandle *orderHandle // 賣單管理
	buyHandle  *orderHandle // 買單管理
	msgIn      chan *engineFlow
}

func New(ctx context.Context) IMatches {
	handle := &engine{
		msgIn:      make(chan *engineFlow, 1024),
		buyHandle:  newOrderHandle(order.OrderType_Buy),
		sellHandle: newOrderHandle(order.OrderType_Sell),
		ctx:        ctx,
	}

	go handle.messageHandle()
	return handle
}

type orderHandle struct {
	priceSlice       []float64                      // 有訂單的價格陣列
	orderMap         map[float64][]*order.OrderInfo // 訂單表格
	marketPriceOrder []*order.OrderInfo             // 未成交市價單
	orderType        order.OrderType
}

func newOrderHandle(orderType order.OrderType) *orderHandle {
	return &orderHandle{
		priceSlice:       make([]float64, 0),
		orderMap:         make(map[float64][]*order.OrderInfo),
		marketPriceOrder: make([]*order.OrderInfo, 0),
		orderType:        orderType,
	}
}
