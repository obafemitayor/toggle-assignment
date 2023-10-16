package database

import "github.com/toggle-assignment/data/models"

type PostgreSQLDatabase struct{}

func (rp *PostgreSQLDatabase) AddReceipt(receiptUUID string) {
	// TODO: Implement
}

func (rp *PostgreSQLDatabase) GetReceipt(receiptUUID string) *models.Receipt {
	// TODO: Implement
	return &models.Receipt{}
}

func (rp *PostgreSQLDatabase) UpdateReceipt(receiptUUID string, receipt models.Receipt) {
	// TODO: Implement

}
