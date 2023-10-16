package data

import "github.com/toggle-assignment/data/database"

func GetDatabase() ReceiptDatabase {
	return &database.MockDatabase{}
}
