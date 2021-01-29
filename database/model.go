package database

type PayInfo struct {
	Wxid    string  `json:"wxid" gorm:"index"`
	TransID string  `json:"trans_id" gorm:"index"`
	Money   float32 `json:"money"`
}
