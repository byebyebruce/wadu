package vlm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/byebyebruce/wadu/pkg/imagex"

	"github.com/sashabaranov/go-openai"
)

func ChatImage(ctx context.Context, cli *openai.Client, model string, prompt string, image []byte, jsonFormat bool) (string, error) {
	if imagex.IsPNG(image) {
		var err error
		image, err = imagex.ConvertPNGtoJPEG(image)
		if err != nil {
			return "", err
		}
	}
	if !imagex.IsJPEG(image) {
		return "", fmt.Errorf("unsupported image, only support jpeg and png")
	}
	imageBase64 := base64.StdEncoding.EncodeToString(image)
	imageURL := fmt.Sprintf("data:image/jpeg;base64,%s", imageBase64)

	m := openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser,
		MultiContent: []openai.ChatMessagePart{
			{Type: openai.ChatMessagePartTypeText, Text: prompt},
			{Type: openai.ChatMessagePartTypeImageURL, ImageURL: &openai.ChatMessageImageURL{URL: imageURL}},
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: []openai.ChatCompletionMessage{m},
	}
	if jsonFormat {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject}
	}

	resp, err := cli.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", err
	}
	ret := resp.Choices[0].Message.Content
	return ret, nil
}

const maxRetry = 3

func ChatImageJSON[T any](ctx context.Context, cli *openai.Client, model string, prompt string, image []byte) (*T, error) {
	for i := 0; i < maxRetry; i++ {
		str, err := ChatImage(ctx, cli, model, prompt, image, true)
		if err != nil {
			return nil, err
		}
		var ret T
		err = json.Unmarshal([]byte(str), &ret)
		if err != nil {
			if i == maxRetry-1 {
				return nil, err
			}
			continue
		}
		return &ret, nil
	}
	return nil, fmt.Errorf("failed to unmarshal json")
}
