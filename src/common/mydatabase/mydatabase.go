package mydatabase

import "github.com/jinzhu/gorm"

type Tx struct {
	GormTx *gorm.DB
}

