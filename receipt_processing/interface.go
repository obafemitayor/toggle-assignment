package receiptprocessor

type ReceiptProcessor interface {
	ExtractDetailsFromReceipt(filePath string) string
}
