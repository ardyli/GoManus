package schema

// AgentState 表示代理的状态
type AgentState int

const (
	// StateIdle 表示代理处于空闲状态
	StateIdle AgentState = iota
	
	// StateRunning 表示代理正在运行
	StateRunning
	
	// StateFinished 表示代理已完成任务
	StateFinished
	
	// StateError 表示代理遇到错误
	StateError
)

// String 返回代理状态的字符串表示
func (s AgentState) String() string {
	switch s {
	case StateIdle:
		return "idle"
	case StateRunning:
		return "running"
	case StateFinished:
		return "finished"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}
