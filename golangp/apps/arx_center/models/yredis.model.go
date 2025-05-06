package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/toheart/functrace"
	"gorm.io/gorm"
)

type UUIDArray []uuid.UUID

func (a *UUIDArray) Scan(src interface{}) error {
	defer functrace.Trace([]interface {
	}{a, src})()
	return pq.Array(a).Scan(src)
}

func (a UUIDArray) Value() (driver.Value, error) {
	defer functrace.Trace([]interface {
	}{a})()
	return pq.Array(a).Value()
}

type Yroom struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key"  json:"id"`
	Name         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	CreateBy     uuid.UUID      `gorm:"type:uuid;" json:"create_by"`
	RW           UUIDArray      `gorm:"type:uuid[];" json:"rw"`
	R            UUIDArray      `gorm:"type:uuid[];" json:"r"`
	IsPublic     bool           `gorm:"not null" json:"is_public"`
	CreateAt     time.Time      `gorm:"not null" json:"create_at"`
	UpdateAt     time.Time      `gorm:"not null" json:"update_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	PositionData string         `gorm:"type:jsonb" json:"position_data,omitempty"`
}

func (y *Yroom) AddKeyIfNotExists(key string) error {
	defer functrace.Trace([]interface {
	}{y, key})()
	var data map[string]interface{}

	if y.PositionData == "" {
		data = make(map[string]interface{})
	} else {

		if err := json.Unmarshal([]byte(y.PositionData), &data); err != nil {
			return err
		}
	}

	if _, exists := data[key]; !exists {
		data[key] = len(data) + 1
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	y.PositionData = string(jsonData)
	return nil
}

func (y *Yroom) PositionDataToJSON() (string, error) {
	defer functrace.Trace([]interface {
	}{y})()
	data, err := json.Marshal(y.PositionData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

type ShareProjectUserAddInput struct {
	ShareEmail  string `json:"email"  binding:"required"`
	ShareLink   string `json:"share_link"  binding:"required"`
	ProjectName string `json:"project_name"  binding:"required"`
	Access      string `json:"access"  binding:"required"`
}

type ShareProjectAddInput struct {
	ProjectName string `json:"project_name"  binding:"required"`
}

type ShareProjectUserRemoveInput struct {
	RemoveEmail string `json:"email"  binding:"required"`
	ProjectName string `json:"project_name"  binding:"required"`
}

type ShareProjectCloseInput struct {
	ProjectName string `json:"project_name"  binding:"required"`
}

type YRedisRoomShareGetResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatorID uuid.UUID `json:"creator_id"`
	IsPublic  bool      `json:"is_public"`
	CreateAt  time.Time `json:"create_at"`
	UpdateAt  time.Time `json:"update_at"`
	IsClosed  bool      `json:"is_closed"`

	Users []YRedisRoomShareGetResponseUserPermission `json:"users"`
}

type YRedisRoomShareGetResponseUserPermission struct {
	ID     uuid.UUID `json:"id"`
	Email  string    `json:"email"`
	Access string    `json:"access"`
	Name   string    `json:"name"`
}
