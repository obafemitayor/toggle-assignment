package receiptprocessor

import "github.com/toggle-assignment/receipt_processing/processors"

func GetReceiptProcessor() ReceiptProcessor {
	return &processors.DocumentAIReceiptProcessor{}
}
