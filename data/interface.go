package data

import "github.com/toggle-assignment/data/models"

type ReceiptDatabase interface {
	AddReceipt(receiptUUID string)
	GetReceipt(receiptUUID string) *models.Receipt
	UpdateReceipt(receiptUUID string, receipt models.Receipt)
}
