package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name                string                `gorm:"unique" gorm:"not null" json:"name"`
	Password            string                `gorm:"not null" json:"-"`
	Rooms               []*Room               `gorm:"many2many:users_rooms" json:"rooms"`
	PaymentDetails      []*PaymentDetail      `gorm:"foreignKey:PayerId" json:"payment_details"`
	PaymentDestinations []*PaymentDestination `gorm:"foreignKey:DestinationUserId" json:"payment_destinations"`
}

type Room struct {
	gorm.Model
	Name                string `gorm:"unique" json:"name"`
	Password            string
	Users               []*User               `gorm:"many2many:users_rooms" json:"users"`
	PaymentDestinations []*PaymentDestination `gorm:"foreignKey:RoomId" json:"payment_destinations"`
}

// PaymentDestination 支払い先
type PaymentDestination struct {
	gorm.Model
	TotalAmount       int              `gorm:"not null"`                    // いくらか
	DestinationUserId uint             `gorm:"not null"`                    // 誰に支払うか
	DestinationUser   User             `gorm:"refecence:DestinationUserId"` // 誰に支払うか
	CreatedUserId     uint             `gorm:"not null"`                    // 誰が作ったか
	CreatedUser       User             `gorm:"reference:CreatedUserId"`     // 誰が作ったか
	RoomId            uint             `gorm:"not null"`                    // どの部屋の支払いか
	Room              Room             // どの部屋の支払いか
	Name              string           `gorm:"not null"` // 何のための支払いか
	Description       string           // 説明
	PaymentDetails    []*PaymentDetail // 支払い詳細
	CompletedAt       *gorm.DeletedAt  // 支払い完了日
}

// PaymentDetail 誰がどの支払いに対していくら支払うか/支払ったか
type PaymentDetail struct {
	gorm.Model
	PayerId              uint          `gorm:"not null"`          // 誰が支払うか
	Payer                User          `gorm:"reference:PayerId"` // 誰が支払うか
	Amount               int           `gorm:"not null"`          // いくらか
	Status               PaymentStatus `gorm:"default:'pending'"` // 支払い状況
	Description          string        // 説明
	PaymentDestinationId uint          `gorm:"not null"` // どの支払いに対する支払いか
	CompletedAt          time.Time
}

type PaymentStatus string

const (
	Pending  PaymentStatus = "pending"
	Complete PaymentStatus = "complete"
)
