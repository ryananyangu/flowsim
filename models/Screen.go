package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/zohari-tech/flowsim/database"
	"github.com/zohari-tech/flowsim/utils"
)

type Metadata map[string]interface{}
type ExternalHandlerFunc func(map[string]interface{}, *Screen) IScreen

type AuditInfo struct {
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy string    `json:"createdby" gorm:"column:createdby"`
	UpdatedBy string    `json:"updatedby" gorm:"column:updatedby"`
}
type Menu struct {
	ID          uint   `json:"id" gorm:"column:id"`
	Telco       string `json:"telco" gorm:"column:telco"`
	Shortcode   string `json:"shortcode" gorm:"column:shortcode"`
	CountryCode string `json:"country_code" gorm:"column:country_code"`
	Description string `json:"description" gorm:"column:description"`
	AuditInfo
}

type Screen struct {
	ID          uint             `json:"id" gorm:"column:id"`
	MenuID      uint             `json:"menu_id" gorm:"column:menu_id"`
	Name        string           `json:"name" gorm:"column:name"`
	ScreenType  utils.ScreenType `json:"screen_type" gorm:"column:screen_type"`
	Language    string           `json:"language" gorm:"column:language"`
	IsEnd       bool             `json:"is_end" gorm:"column:is_end"`
	Location    uint             `json:"location" gorm:"column:location"` // NOTE: This is tracked within the menu
	BackEnabled bool             `json:"back_enabled" gorm:"column:back_enabled"`
	ExitEnabled bool             `json:"exit_enabled" gorm:"column:exit_enabled"`
	Details     Metadata         `json:"details" gorm:"column:details"`
	AuditInfo
}

type Message struct {
	ID                uint            `json:"id" gorm:"column:id"`
	ScreenID          uint            `json:"screen_id" gorm:"column:screen_id"`
	ConversationID    string          `json:"conversation_id" gorm:"column:conversation_id"`
	Content           string          `json:"content" gorm:"column:content"`
	Source            string          `json:"source" gorm:"column:source"`
	Destination       string          `json:"destination" gorm:"column:destination"`
	Direction         utils.Direction `json:"direction" gorm:"column:direction"`
	Status            string          `json:"status" gorm:"column:status"`
	StatusDescription string          `json:"status_description" gorm:"column:status_description"`
	AuditInfo
	SessionData Metadata `json:"session_data" gorm:"column:session_data"`
}

func (msg *Message) GetLastScreen() (screen Screen, start bool) {
	// NOTE: This is where we retrieve the previous screen from cache
	iscreen, active := utils.CacheInstance.Get(msg.ConversationID)
	if active {
		screen = iscreen.(Screen)
		return
	}
	resultset := database.Db.Table("messages").Where("conversation_id = ?", msg.ConversationID).Order("created_at DESC").First(msg)
	if resultset.RowsAffected == 0 {
		start = true
		return
	}
	_ = database.Db.Table("screens").Where("id = ?", msg.ScreenID).First(&screen).Error
	return
}

// Value Marshal
func (a Metadata) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Scan Unmarshal
func (a *Metadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}

