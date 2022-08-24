package matches

import (
	"example/matches/order"
)

func (self *engine) NewOrder(order *order.OrderInfo) {
	// fmt.Println("NewOrder:", order.OrderType, order.TransactionType, order.Price, order.Quantity)
	self.msgIn <- &engineFlow{
		EventCode: engineEvent_NewOrder,
		Order:     order,
	}
}

func (self *engine) CheckOrder() {
	self.msgIn <- &engineFlow{
		EventCode: engineEvent_CheckOrder,
	}
}
