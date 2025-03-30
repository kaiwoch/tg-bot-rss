package summary

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sheeiavellie/go-yandexgpt"
)

type YandexGPTAISummarizer struct {
	client    *yandexgpt.YandexGPTClient
	prompt    string
	model     string
	catalogID string
	enabled   bool
	mu        sync.Mutex
}

func NewYandexGPTSummarizer(apiKey, model, prompt, catalogID string) *YandexGPTAISummarizer {

	y := &YandexGPTAISummarizer{
		client:    yandexgpt.NewYandexGPTClientWithIAMToken(apiKey),
		prompt:    prompt,
		model:     model,
		catalogID: catalogID,
	}

	log.Printf("yandexGPT summarizer is enabled: %v", apiKey != "")
	fmt.Println(apiKey)
	if apiKey != "" {
		y.enabled = true
	}

	return y
}

func (y *YandexGPTAISummarizer) Summarize(text string) (string, error) {
	y.mu.Lock()
	defer y.mu.Unlock()

	if !y.enabled {
		return "", fmt.Errorf("openai summarizer is disabled")
	}

	request := yandexgpt.YandexGPTRequest{
		ModelURI: yandexgpt.MakeModelURI(y.catalogID, yandexgpt.YandexGPT4Model),
		CompletionOptions: yandexgpt.YandexGPTCompletionOptions{
			Stream:      false,
			Temperature: 0.6,
			MaxTokens:   2000,
		},
		Messages: []yandexgpt.YandexGPTMessage{
			{
				Role: yandexgpt.YandexGPTMessageRoleSystem,
				Text: y.prompt,
			},
			{
				Role: yandexgpt.YandexGPTMessageRoleUser,
				Text: text,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	resp, err := y.client.GetCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(resp.Result.Alternatives) == 0 {
		return "", errors.New("no choices in openai response")
	}

	rawSummary := strings.TrimSpace(resp.Result.Alternatives[0].Message.Text)
	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	// cut all after the last ".":
	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}
