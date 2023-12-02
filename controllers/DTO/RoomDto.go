package dto

import (
	"gorm.io/gorm"
	"waroka/model"
)

type UserDto struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type PaymentDestinationDto struct {
	Id                  uint            `json:"id"`
	TotalAmount         int             `json:"total_amount"`
	DestinationUserId   uint            `json:"destination_user_id"`
	DestinationUserName string          `json:"destination_user_name"`
	CreatedUserId       uint            `json:"created_user_id"`
	CreatedUserName     string          `json:"created_user_name"`
	RoomId              uint            `json:"room_id"`
	RoomName            string          `json:"room_name"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	CompletedAt         *gorm.DeletedAt `json:"completed_at"`
}

type PaymentDetailDto struct {
	Id                   uint                `json:"id"`
	PayerId              uint                `json:"payer_id"`
	PayerName            string              `json:"payer_name"`
	Amount               int                 `json:"amount"`
	Status               model.PaymentStatus `json:"status"`
	Description          string              `json:"description"`
	PaymentDestinationId uint                `json:"payment_destination_id"`
	CompletedAt          *gorm.DeletedAt     `json:"completed_at"`
}

type RoomResponse struct {
	Name                string                      `json:"name"`
	Id                  uint                        `json:"id"`
	Users               []UserDto                   `json:"users"`
	PaymentDestinations []*model.PaymentDestination `json:"payment_destinations"`
	TotalAmount         int                         `json:"total_amount"`
	TotalActiveAmount   int                         `json:"total_active_amount"`
}
