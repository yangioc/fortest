package matches

import (
	"example/matches/order"
	"example/util"
	"fmt"
	"sort"
)

func (self *orderHandle) InsertOrder(reqOrder *order.OrderInfo) error {
	if reqOrder.TransactionType == order.TransactionType_MarketPrice { // 新增 市價訂單
		self.marketPriceOrder = append(self.marketPriceOrder, reqOrder)
		return nil

	} else if reqOrder.TransactionType == order.TransactionType_LimitPrice { // 新增 限價訂單
		if _, exist := self.orderMap[reqOrder.Price]; exist {
			self.orderMap[reqOrder.Price] = append(self.orderMap[reqOrder.Price], reqOrder)

		} else {
			if _, ok := self.orderMap[reqOrder.Price]; !ok {
				self.priceSlice = append(self.priceSlice, reqOrder.Price)
				if self.orderType == order.OrderType_Sell { // 更新 賣單類型排序
					sort.Slice(self.priceSlice, func(i, j int) bool {
						return self.priceSlice[i] < self.priceSlice[j] || (util.IsNaN(self.priceSlice[i]) && !util.IsNaN(self.priceSlice[j]))
					})
				} else { // 更新 買單類型排序
					sort.Float64s(self.priceSlice)
				}
			}

			self.orderMap[reqOrder.Price] = make([]*order.OrderInfo, 0)
			self.orderMap[reqOrder.Price] = append(self.orderMap[reqOrder.Price], reqOrder)
		}
		return nil

	} else {
		return fmt.Errorf("[order][InsertOrder] TransactionType Error: %v", reqOrder.TransactionType)
	}
}

// 尋找符合的訂單
func (self *orderHandle) MatchOrder(reqOrder *order.OrderInfo) []*order.OrderInfo {
	settlementOrders := []*order.OrderInfo{}
	if reqOrder.TransactionType == order.TransactionType_MarketPrice {
		settlementOrders = self.matchMarketPriceOrder(reqOrder)
	} else if reqOrder.TransactionType == order.TransactionType_LimitPrice {
		settlementOrders = self.matchLimitPriceOrder(reqOrder)
	} else {
		fmt.Println("[order][MatchOrder] TransactionType Error:", reqOrder.TransactionType)
	}
	return settlementOrders
}

// 匹配限價訂單
// Q: 匹配問題 部分未完成的訂單如何處理?
func (self *orderHandle) matchLimitPriceOrder(reqOrder *order.OrderInfo) []*order.OrderInfo {

	settlementOrders := []*order.OrderInfo{}
	var limitPriceOrder, marketPriceOrder *order.OrderInfo

	for Quantity, isNext := reqOrder.Quantity, true; Quantity > 0 && isNext; {
		limitPriceOrder = nil
		marketPriceOrder = nil
		price := reqOrder.Price

		// 取得限價訂單
		if len(self.priceSlice) > 0 && self.priceSlice[0] < price { // 存在更優價格
			// 取得價格最接近訂單
			price = self.priceSlice[0]
		}

		if orders, ok := self.orderMap[price]; ok && len(orders) > 0 {
			limitPriceOrder = self.orderMap[price][0]
		}

		// 取得時價訂單
		if len(self.marketPriceOrder) > 0 {
			marketPriceOrder = self.marketPriceOrder[0]
		}

		// 比對訂單
		var settlementOrder *order.OrderInfo
		switch {
		case limitPriceOrder != nil && marketPriceOrder != nil: // 兩種訂單都有
			// 判斷訂單優先度
			if limitPriceOrder.OrderStartTime > marketPriceOrder.OrderStartTime { // 市價單優先
				if marketPriceOrder.Quantity > Quantity {
					settlementOrder = self.settlementMarketPriceOrder(marketPriceOrder.Price, Quantity)
				} else {
					settlementOrder = self.settlementMarketPriceOrder(marketPriceOrder.Price, marketPriceOrder.Quantity)
				}
			} else { // 限價單優先
				if limitPriceOrder.Quantity > Quantity {
					settlementOrder = self.settlementPriceOrder(limitPriceOrder.Price, Quantity)
				} else {
					settlementOrder = self.settlementPriceOrder(limitPriceOrder.Price, limitPriceOrder.Quantity)
				}
			}

		case limitPriceOrder != nil: // 只存在限價訂單
			if limitPriceOrder.Quantity > Quantity {
				settlementOrder = self.settlementPriceOrder(limitPriceOrder.Price, Quantity)
			} else {
				settlementOrder = self.settlementPriceOrder(limitPriceOrder.Price, limitPriceOrder.Quantity)
			}

		case marketPriceOrder != nil: // 只存在時價訂單
			if marketPriceOrder.Quantity > Quantity {
				settlementOrder = self.settlementMarketPriceOrder(marketPriceOrder.Price, Quantity)
			} else {
				settlementOrder = self.settlementMarketPriceOrder(marketPriceOrder.Price, marketPriceOrder.Quantity)
			}

		default:
			isNext = false
		}

		if isNext {
			settlementOrders = append(settlementOrders, settlementOrder)
			Quantity -= settlementOrder.Quantity
		}
	}
	return settlementOrders
}

func (self *orderHandle) matchMarketPriceOrder(reqOrder *order.OrderInfo) []*order.OrderInfo {
	// 未存在任何訂單
	if len(self.priceSlice) <= 0 {
		return nil
	}

	// 調整訂單資料
	settlementOrders := []*order.OrderInfo{}
	var price float64
	for Quantity := reqOrder.Quantity; Quantity > 0; {
		if len(self.priceSlice) <= 0 {
			break
		}
		price = self.priceSlice[0]
		settlementOrder := self.settlementPriceOrder(price, Quantity)
		settlementOrders = append(settlementOrders, settlementOrder)
		Quantity -= settlementOrder.Quantity
	}
	return settlementOrders
}

func (self *orderHandle) settlementPriceOrder(price float64, quantity int64) *order.OrderInfo {
	var settlementOrder order.OrderInfo
	// 調整訂單資料
	orders := self.orderMap[price]
	if len(orders) > 1 {
		if orders[0].Quantity <= quantity {
			settlementOrder = *self.orderMap[price][0]
			self.orderMap[price] = self.orderMap[price][1:]
		} else {
			settlementOrder = *self.orderMap[price][0]
			settlementOrder.Quantity = quantity
			self.orderMap[price][0].Quantity -= quantity
		}

	} else {
		if orders[0].Quantity <= quantity {
			settlementOrder = *self.orderMap[price][0]
			self.priceSlice = self.priceSlice[1:]
			delete(self.orderMap, price)
		} else {
			settlementOrder = *self.orderMap[price][0]
			settlementOrder.Quantity = quantity
			self.orderMap[price][0].Quantity -= quantity
		}
	}
	return &settlementOrder
}

func (self *orderHandle) settlementMarketPriceOrder(price float64, quantity int64) *order.OrderInfo {
	var settlementOrder order.OrderInfo
	if self.marketPriceOrder[0].Quantity <= quantity {
		settlementOrder = *self.marketPriceOrder[0]
		self.marketPriceOrder = self.marketPriceOrder[1:]
	} else {
		settlementOrder = *self.marketPriceOrder[0]
		settlementOrder.Quantity = quantity
		self.marketPriceOrder[0].Quantity -= quantity
	}
	return &settlementOrder
}
