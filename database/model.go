package database

type PayInfo struct {
	Wxid    string  `json:"wxid" gorm:"index;type:varchar(32) not null"`
	TransID string  `json:"trans_id" gorm:"unique;type:varchar(32) not null"`
	Money   float64 `json:"money" gorm:"type:decimal(8,2) not null"`
}

type Member struct {
	Wxid  string  `json:"wxid" gorm:"unique;type:varchar(32) not null"`
	Money float64 `json:"money" gorm:"type:decimal(8,2) not null"`
}
