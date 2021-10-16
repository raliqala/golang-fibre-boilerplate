package Users

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id          uint    `gorm:"primaryKey;AUTO_INCREMENT"`
	Name        string  `gorm:"type:varchar(300);NOT NULL"`
	Email       *string `gorm:"type:varchar(100);NOT NULL;unique_index"`
	Password    string  `gorm:"NOT NULL"`
	ActivatedAt sql.NullTime
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
