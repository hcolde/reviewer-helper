package database

type PublishMsg struct {
	Wxid string `json:"wxid"`
	Msg  string `json:"msg"`
}

type PayInfo struct {
	Wxid    string  `json:"wxid" gorm:"index;type:varchar(32) not null"`
	TransID string  `json:"trans_id" gorm:"primaryKey;type:varchar(32) not null"`
	Money   float64 `json:"money" gorm:"type:decimal(8,2) not null"`
}

type Member struct {
	Wxid  string  `json:"wxid" gorm:"primaryKey;type:varchar(32) not null"`
	Money float64 `json:"money" gorm:"type:decimal(8,2) not null"`
	Vip   int64   `json:"vip" gorm:"type:bigint not null default 0"`
}

type VipInfo struct {
	Wxid  string  `json:"wxid" gorm:"index;type:varchar(32) not null"`
	Money float64 `json:"money" gorm:"type:decimal(8,2) not null"`
}
