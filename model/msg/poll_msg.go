package msg

/**
 * 轮询请求消息
 */
type PollMsg struct {
	DNodeId   int64
	DNodeAddr string
	Offset    int
	Session   string //会话标识
	Token     string //认证信息
}
