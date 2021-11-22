package gim

type ConversationType int

const (
	ConversationType_Default   = iota // 普通对话
	ConversationType_Transient        // 临时对话
	ConversationType_System           // 系统对话
)

type Conversation struct {
	Type     ConversationType // 对话类型
	ID       uint64           // 全局唯一ID
	Creator  uint64           // 对话创建者
	CreateAt int64            // 创建时间戳 秒
	Members  []uint64         // 成员列表
	Name     string           // 会话名
}
