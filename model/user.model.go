package model

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type DeletedAt sql.NullTime

type Model struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type User struct {
	Model
	Name                string                `gorm:"unique" gorm:"not null" json:"name"`
	Password            string                `gorm:"not null" json:"-"`
	Rooms               []*Room               `gorm:"many2many:users_rooms" json:"rooms"`
	PaymentDetails      []*PaymentDetail      `gorm:"foreignKey:PayerId" json:"payment_details"`
	PaymentDestinations []*PaymentDestination `gorm:"foreignKey:DestinationUserId" json:"payment_destinations"`
}

type Room struct {
	Model
	Name                string                `gorm:"unique" json:"name"`
	Password            string                `json:"-"`
	Users               []*User               `gorm:"many2many:users_rooms" json:"users"`
	PaymentDestinations []*PaymentDestination `gorm:"foreignKey:RoomId" json:"payment_destinations"`
}

// PaymentDestination 支払い先
type PaymentDestination struct {
	Model
	TotalAmount       int              `gorm:"not null" json:"total_amount"`                        // いくらか
	DestinationUserId uint             `gorm:"not null" json:"destination_user_id"`                 // 誰に支払うか
	DestinationUser   User             `gorm:"refecence:DestinationUserId" json:"destination_user"` // 誰に支払うか
	CreatedUserId     uint             `gorm:"not null" json:"created_user_id"`                     // 誰が作ったか
	CreatedUser       User             `gorm:"reference:CreatedUserId" json:"create_user"`          // 誰が作ったか
	RoomId            uint             `gorm:"not null" json:"room_id"`                             // どの部屋の支払いか
	Room              Room             `json:"room"`                                                // どの部屋の支払いか
	Name              string           `gorm:"not null" json:"name"`                                // 何のための支払いか
	Description       string           `json:"description"`                                         // 説明
	PaymentDetails    []*PaymentDetail `json:"payment_details"`                                     // 支払い詳細
	CompletedAt       *gorm.DeletedAt  `json:"completed_at"`                                        // 支払い完了日
}

// PaymentDetail 誰がどの支払いに対していくら支払うか/支払ったか
type PaymentDetail struct {
	Model
	PayerId              uint          `gorm:"not null" json:"payer_id"`             // 誰が支払うか
	Payer                User          `gorm:"reference:PayerId" json:"payer"`       // 誰が支払うか
	Amount               int           `gorm:"not null" json:"amount"`               // いくらか
	Status               PaymentStatus `gorm:"default:'pending'" json:"status"`      // 支払い状況
	Description          string        `json:"description"`                          // 説明
	PaymentDestinationId uint          `gorm:"not null" json:"paymentDestinationId"` // どの支払いに対する支払いか
	CompletedAt          time.Time
}

type PaymentStatus string

const (
	Pending  PaymentStatus = "pending"
	Complete PaymentStatus = "complete"
)
