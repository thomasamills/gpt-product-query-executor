package chat_cpt_client

import (
	"chat-gpt-product-spec/db"
	"errors"
	"fmt"
	gpt35 "github.com/AlmazDelDiablo/gpt3-5-turbo-go"
)

type ChatGptClient interface {
	SendPrompt(
		sku string,
		query string,
		url string,
		datasheet string,
		catalogue string,
		dataSheetContent string,
		catalogueContent string,
	) error
}

type ChatGptClientImpl struct {
	gpt35Client *gpt35.Client
	db          db.GptProductDatabase
}

func NewChatGptClient(db db.GptProductDatabase, apiKey string) ChatGptClient {
	return &ChatGptClientImpl{
		gpt35Client: gpt35.NewClient(apiKey),
		db:          db,
	}
}

func (c *ChatGptClientImpl) SendPrompt(
	sku string,
	query string,
	url string,
	datasheet string,
	catalogue string,
	dataSheetContent string,
	catalogueContent string,
) error {
	prompt := query
	if len(dataSheetContent) > 0 {
		prompt += " to aid you, here is some data from a website: " + dataSheetContent
	}
	if len(catalogueContent) > 0 {
		prompt += " to aid you, here is some data from a website: " + catalogueContent
	}
	req := &gpt35.Request{
		Model: gpt35.ModelGpt35Turbo,
		Messages: []*gpt35.Message{
			{
				Role:    gpt35.RoleUser,
				Content: prompt,
			},
		},
		MaxTokens: 1000,
	}
	resp, err := c.gpt35Client.GetChat(req)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		fmt.Println(resp.Error)
		return errors.New(resp.Error.Message)
	}
	content := ""
	for _, choice := range resp.Choices {
		if choice.Message != nil {
			content = choice.Message.Content
			break
		}
	}
	// now save to db
	fmt.Println(content)
	return c.db.CreateProductSpec(sku, query, url, datasheet, catalogue, content, nil)
}
