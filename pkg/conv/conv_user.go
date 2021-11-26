package conv

type CMember struct {
	ID       string
	ConvID   int64
	UserID   string
	Nickname string
	CreateAt int64
	Mute     int // 禁言
	Role     int // 会话角色
}
