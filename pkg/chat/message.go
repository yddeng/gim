package gim

// 消息队列
type MessageQueue struct {
	ID      uint64 // 全局唯一ID
	Counter uint64 // 消息计数器
}

type MessageEntity struct {
	MID  uint64 // 消息ID
	Text string // 消息
}
