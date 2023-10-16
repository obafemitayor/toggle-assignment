package database

import "github.com/toggle-assignment/data/models"

type MockDatabase struct{}

var receipts = []models.Receipt{
	{ID: "1", UUID: "5a2ac852-1489-47ca-9487-82194087f934", Details: "Mock Receipt Details", IsProcessing: true},
}

func (rp *MockDatabase) AddReceipt(receiptUUID string) {
	var newReceipt models.Receipt
	newReceipt.UUID = receiptUUID
	newReceipt.IsProcessing = true
	receipts = append(receipts, newReceipt)
}

func (rp *MockDatabase) GetReceipt(receiptUUID string) *models.Receipt {
	for i := range receipts {
		if receipts[i].UUID == receiptUUID {
			return &receipts[i]
		}
	}
	return nil
}

func (rp *MockDatabase) UpdateReceipt(receiptUUID string, receipt models.Receipt) {
	for i := range receipts {
		if receipts[i].UUID == receiptUUID {
			receipts[i].Details = receipt.Details
			receipts[i].IsProcessing = receipt.IsProcessing
		}
	}
}
