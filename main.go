package main

import (
	chat_cpt_client "chat-gpt-product-spec/chat_gpt_client"
	"chat-gpt-product-spec/data_read"
	"chat-gpt-product-spec/db"
	"fmt"
	"os"
)

func main() {
	db := db.NewGptProductDatabase()
	gptClient := chat_cpt_client.NewChatGptClient(db, os.Args[2])
	dataParser := data_read.NewCsvParser()
	pdfExtractor := data_read.NewPdfExtractor()
	fmt.Println(os.Args[1])
	data := dataParser.Parse(os.Args[1])
	for i, dataItem := range data {
		if i == 0 {
			continue
		}
		// checking if its already been processed
		exists, err := db.DoesProductSpecExist(dataItem.Sku, nil)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exists {
			fmt.Println("Skipping item " + dataItem.Sku + " as it has already been processed")
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
			dataText = text
		} else if len(dataItem.CataloguePage) > 0 {
			text, err := pdfExtractor.Extract(dataItem.CataloguePage, dataItem.Sku+"cat")
			if err != nil {
				fmt.Println("could not extract data from pdf")
			}
			catText = text
		}
		err = gptClient.SendPrompt(
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
