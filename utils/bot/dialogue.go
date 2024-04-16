package bot

import (
	"github.com/sashabaranov/go-openai"
)

func Dialogue(message string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "in starting you will ask for username ?",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You can help to upload photos ask the user to give an username click upload button below and upload ?",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "you will ask the user if they want to upload photos or retrieve them ?",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are Alexia, a helpful AI assistant",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	}
}
