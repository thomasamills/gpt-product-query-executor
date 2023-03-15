package main

import (
	chat_cpt_client "chat-gpt-product-spec/chat_gpt_client"
	"chat-gpt-product-spec/data_read"
	"chat-gpt-product-spec/db"
	"fmt"
)

func firstN(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}
	return s
}
func main() {
	db := db.NewGptProductDatabase()
	gptClient := chat_cpt_client.NewChatGptClient(db)
	dataParser := data_read.NewCsvParser()
	pdfExtractor := data_read.NewPdfExtractor()
	data := dataParser.Parse("./Allproducts-latest.csv")
	for i, dataItem := range data {
		if i == 0 {
			continue
		}
		fmt.Println("Processing item " + dataItem.Sku)
		dataText := ""
		catText := ""
		if len(dataItem.Datasheet) > 0 {
			text, err := pdfExtractor.Extract(dataItem.Datasheet, dataItem.Sku+"data")
			if err != nil {
				fmt.Println("could not extract data from pdf")
			}
			dataText = firstN(text, 3800)
		} else if len(dataItem.CataloguePage) > 0 {
			text, err := pdfExtractor.Extract(dataItem.CataloguePage, dataItem.Sku+"cat")
			if err != nil {
				fmt.Println("could not extract data from pdf")
			}
			catText = firstN(text, 3800)
		}
		err := gptClient.SendPrompt(
			dataItem.Sku,
			dataItem.Query,
			dataItem.MasterTradeURL,
			dataItem.Datasheet,
			dataItem.CataloguePage,
			dataText,
			catText,
		)
		if err != nil {
			fmt.Println(err)
		}
	}
}
