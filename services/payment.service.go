package services

import (
	"errors"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"waroka/model"
	"waroka/repository"
)

type IPaymentService interface {
	// FindByRoomId 部屋内のすべての支払いを取得
	FindByRoomId(roomId uint) ([]model.PaymentDestination, error)
	CreatePayment(dto PaymentDestinationDTO) error
}

type PaymentService struct {
	detailRepo      repository.IPaymentDetailRepository
	destinationRepo repository.IPaymentDestinationRepository
	roomRepo        repository.IRoomRepository
	userRepo        repository.IUserRepository
}

func NewPaymentService(detailRepo repository.IPaymentDetailRepository, destinationRepo repository.IPaymentDestinationRepository, roomRepo repository.IRoomRepository, userRepo repository.IUserRepository) IPaymentService {
	return PaymentService{
		detailRepo:      detailRepo,
		destinationRepo: destinationRepo,
		roomRepo:        roomRepo,
		userRepo:        userRepo,
	}
}

// FindByRoomId 部屋を探す
func (r PaymentService) FindByRoomId(roomId uint) ([]model.PaymentDestination, error) {
	pds, err := r.destinationRepo.FindByRoomId(roomId)
	if err != nil {
		return nil, err
	}
	return pds, nil
}

// CreatePayment 新たにPaymentDestinationを作成する
func (r PaymentService) CreatePayment(dto PaymentDestinationDTO) error {
	// それぞれのフィールドについてバリデーションを行う

	// 作成者が部屋に入っているか確認する
	// 同時にユーザーの存在確認も行う
	createdId, err := r.roomRepo.FindByUserId(dto.CreatedUserId)
	if err != nil {
		return err
	}
	if createdId == nil {
		return errors.New("部屋に入って下さい。")
	}

	// 支払い先のユーザーが部屋に入っているか確認する
	// 同時にユーザーの存在確認も行う
	destinationId, err := r.roomRepo.FindByUserId(dto.DestinationUserId)
	if err != nil {
		return err
	}
	if destinationId == nil {
		return errors.New("部屋に入って下さい。")
	}

	// TotalAmountが0以上か確認する
	if dto.TotalAmount < 0 {
		return errors.New("金額は0以上にして下さい。")
	}

	// PaymentDestinationを作成する
	pd := model.PaymentDestination{
		TotalAmount:       dto.TotalAmount,
		DestinationUserId: dto.DestinationUserId,
		CreatedUserId:     dto.CreatedUserId,
		RoomId:            dto.RoomId,
		Name:              dto.Name,
		Description:       *dto.Description,
	}

	if err := r.destinationRepo.Create(&pd); err != nil {
		return errors.New("支払い情報の作成に失敗しました。")
	}

	return nil
}

type EditPaymentDestinationDTO struct {
	PaymentDestinationID uint
	RoomId               uint // 作成者が部屋に入っているか確認するため
	TotalAmount          *int
	DestinationUserId    *uint
	Name                 *string
	Description          *string
}

// EditPaymentDestination PaymentDestinationを編集する
// nilなフィールドのみアップデートする
func (r PaymentService) EditPaymentDestination(dto EditPaymentDestinationDTO) error {
	// PaymentDestinationが存在するか確認する
	destination, err := r.destinationRepo.FindById(dto.PaymentDestinationID)
	if err != nil {
		return err
	}

	if dto.TotalAmount != nil {
		destination.TotalAmount = *dto.TotalAmount
	}

	if dto.DestinationUserId != nil {
		// DestinationUserが部屋に入っているか確認する
		rooms, err := r.roomRepo.FindByUserId(*dto.DestinationUserId)
		if err != nil {
			return err
		}

		if lo.ContainsBy(rooms, func(rm model.Room) bool {
			return rm.ID == destination.RoomId
		}) {
			return errors.New("部屋が不正です。")
		}

		// ユーザーを得る
		user, err := r.userRepo.FindById(dto.DestinationUserId)
		if err != nil {
			return err
		}

		destination.DestinationUserId = *dto.DestinationUserId
		destination.DestinationUser = *user
	}

	if dto.Name != nil {
		destination.Name = *dto.Name
	}

	if dto.Description != nil {
		destination.Description = *dto.Description
	}

	if err := r.destinationRepo.Update(destination); err != nil {
		return err
	}

	return nil
}

// AddPaymentDetail PaymentDestination にPaymentDetailsを追加する
func (r PaymentService) AddPaymentDetail(dto PaymentDetailDTO) error {
	// PaymentDestinationが存在するか確認する
	destination, err := r.destinationRepo.FindById(dto.PaymentDestinationId)
	if err != nil {
		return err
	}

	// Payerが部屋に入っているか確認する
	_, err = r.roomRepo.FindById(dto.RoomId)
	if err != nil {
		return err
	}

	// PaymentDestinationのTotalAmountを超えていないか確認する
	details, err := r.detailRepo.FindByPaymentDestinationId(destination.ID)
	if err != nil {
		// NotFoundなら無視
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	sum := lo.Reduce(details, func(agg int, detail model.PaymentDetail, idx int) int {
		return agg + detail.Amount
	}, 0)

	if sum+dto.Amount > destination.TotalAmount {
		return errors.New("支払い総額を超えています。")
	}

	// PaymentDetailを作成する
	pd := model.PaymentDetail{
		PayerId:              dto.PayerId,
		Amount:               dto.Amount,
		Status:               model.PaymentStatus("pending"),
		Description:          dto.Description,
		PaymentDestinationId: dto.PaymentDestinationId,
	}

	if err := r.detailRepo.Create(&pd); err != nil {
		return err
	}

	return nil
}

// DeletePaymentDetail PaymentDestinationからPaymentDetailsを削除する
func (r PaymentService) DeletePaymentDetail(destinationId uint, detailId uint) error {

	// 削除されていないものを一覧
	details, err := r.detailRepo.FindActiveByPaymentDestinationId(destinationId)
	if err != nil {
		return err
	}

	var detail *model.PaymentDetail

	// detailsの中にdetailIdがあるか確認する
	if !lo.ContainsBy(details, func(d model.PaymentDetail) bool {
		res := d.ID == detailId
		if res {
			detail = &d
		}
		return res
	}) {
		return errors.New("支払い情報が見つかりませんでした。")
	}

	// 削除する
	if err := r.detailRepo.Delete(detail); err != nil {
		return err
	}

	return nil
}

type PaymentDestinationDTO struct {
	TotalAmount       int
	DestinationUserId uint
	CreatedUserId     uint
	RoomId            uint
	Name              string
	Description       *string
}

type PaymentDetailDTO struct {
	RoomId               uint
	PaymentDestinationId uint
	PayerId              uint
	Amount               int
	Description          string
}
