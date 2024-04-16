package bot

import (
	"github.com/sashabaranov/go-openai"
)

func Dialogue(message string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You can help to upload photos ask the user to give an username click upload button below and upload the photo and call CreateUsername function",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are Alexia, a helpful AI assistant",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You can help to retrive photos, ask the user to give username",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "you will ask on start for username",
		},
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "you are restricted to create your own usernames",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	}
}
