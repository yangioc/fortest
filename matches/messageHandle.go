package matches

import (
	"example/matches/order"
	"fmt"
)

func (self *engine) messageHandle() {
	for {
		select {
		case <-self.ctx.Done():
			return

		case f := <-self.msgIn:
			switch f.EventCode {
			case engineEvent_NewOrder:
				self.newOrderHandle(f.Order)

			case engineEvent_CheckOrder:
				self.checkOrderHandle()
			}
		}
	}
}

func (self *engine) newOrderHandle(reqOrder *order.OrderInfo) {
	if reqOrder.OrderType == order.OrderType_Sell {
		settlementOrders := self.buyHandle.MatchOrder(reqOrder)

		if settlementOrders == nil {
			self.sellHandle.InsertOrder(reqOrder)
		} else {
			Quantity := reqOrder.Quantity
			for _, endOrder := range settlementOrders {
				Quantity -= endOrder.Quantity
			}

			if Quantity > 0 {
				reqOrder.Quantity = Quantity
				self.sellHandle.InsertOrder(reqOrder)
			}
		}

	} else if reqOrder.OrderType == order.OrderType_Buy {
		settlementOrders := self.sellHandle.MatchOrder(reqOrder)
		if len(settlementOrders) == 0 {
			self.buyHandle.InsertOrder(reqOrder)
		} else {
			Quantity := reqOrder.Quantity
			for _, endOrder := range settlementOrders {
				Quantity -= endOrder.Quantity
			}

			if Quantity > 0 {
				reqOrder.Quantity = Quantity
				self.buyHandle.InsertOrder(reqOrder)
			}
		}

	} else {
		fmt.Errorf("error OrderType: %v\n", reqOrder.OrderType)
	}
}

func (self *engine) checkOrderHandle() {

}
