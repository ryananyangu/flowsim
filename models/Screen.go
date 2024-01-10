package models

import "time"

type Metadata map[string]interface{}

type Direction string

const (
	INBOX  Direction = "INBOX"
	OUTBOX Direction = "OUTBOX"
)

type AuditInfo struct {
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	CreatedBy string    `json:"createdby" gorm:"column:createdby"`
	UpdatedBy string    `json:"updatedby" gorm:"column:updatedby"`
}
type Menu_ struct {
	ID          uint   `json:"id" gorm:"column:id"`
	Shortcode   string `json:"shortcode" gorm:"column:shortcode"`
	CountryCode string `json:"country_code" gorm:"column:country_code"`
	Description string `json:"description" gorm:"column:description"`
	AuditInfo
}

type Screen_ struct {
	ID          uint     `json:"id" gorm:"column:id"`
	MenuID      string   `json:"menu_id" gorm:"column:menu_id"`
	Name        string   `json:"name" gorm:"column:name"`
	Language    string   `json:"language" gorm:"column:language"`
	IsEnd       string   `json:"is_end" gorm:"column:is_end"`
	BackEnabled string   `json:"back_enabled" gorm:"column:back_enabled"`
	ExitEnabled string   `json:"exit_enabled" gorm:"column:exit_enabled"`
	Details     Metadata `json:"details" gorm:"column:details"`
	AuditInfo
}

type Message_ struct {
	ID                uint      `json:"id" gorm:"column:id"`
	ScreenID          uint      `json:"screen_id" gorm:"column:screen_id"`
	ConversationID    string    `json:"conversation_id" gorm:"column:conversation_id"`
	Content           string    `json:"content" gorm:"column:content"`
	Source            string    `json:"source" gorm:"column:source"`
	Destination       string    `json:"destination" gorm:"column:destination"`
	Direction         Direction `json:"direction" gorm:"column:direction"`
	Status            string    `json:"status" gorm:"column:status"`
	StatusDescription string    `json:"status_description" gorm:"column:status_description"`
	AuditInfo
	SessionData Metadata `json:"session_data" gorm:"column:session_data"`
}
