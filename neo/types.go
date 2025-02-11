package neo

import (
	"io"

	"github.com/yaoapp/yao/aigc"
	"github.com/yaoapp/yao/neo/conversation"
)

// DSL AI assistant
type DSL struct {
	ID                  string                    `json:"-" yaml:"-"`
	Name                string                    `json:"name,omitempty"`
	Use                 string                    `json:"use,omitempty"`
	Guard               string                    `json:"guard,omitempty"`
	Connector           string                    `json:"connector"`
	ConversationSetting conversation.Setting      `json:"conversation" yaml:"conversation"`
	Option              map[string]interface{}    `json:"option"`
	Prepare             string                    `json:"prepare,omitempty"`
	Prompts             []aigc.Prompt             `json:"prompts,omitempty"`
	Allows              []string                  `json:"allows,omitempty"`
	AI                  aigc.AI                   `json:"-" yaml:"-"`
	Conversation        conversation.Conversation `json:"-" yaml:"-"`
}

// Answer the answer interface
type Answer interface {
	Stream(func(w io.Writer) bool) bool
	Status(code int)
	Header(key, value string)
}
