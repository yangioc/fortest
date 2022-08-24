package order

type OrderType = uint8
type TransactionType = uint8

// OrderType 訂單類型
const (
	OrderType_None OrderType = 0 // 預設值
	OrderType_Buy  OrderType = 1 // 買單
	OrderType_Sell OrderType = 2 // 賣單
)

// TransactionType 交易方式
const (
	TransactionType_None        TransactionType = 0 // 預設值
	TransactionType_LimitPrice  TransactionType = 1 // 限價
	TransactionType_MarketPrice TransactionType = 2 // 市價
)

type OrderInfo struct {
	OrderType       OrderType
	TransactionType TransactionType
	Price           float64
	Quantity        int64
	OrderStartTime  int64
}
