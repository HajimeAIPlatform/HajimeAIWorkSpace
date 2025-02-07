package agent_adapter

type Adapter interface {
	CreateConversation(message string) (Conversation, error)
}

type Conversation interface {
	SendMessage(message string) error
}
