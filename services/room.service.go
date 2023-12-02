package services

import (
	"errors"
	"waroka/model"
	"waroka/repository"
	"waroka/utils"
)

type IRoomService interface {
	CreateRoom(room CreateRoomReq, user *model.User) (*model.Room, error)
	JoinRoom(roomReq JoinRoomReq, user *model.User) (*model.Room, error)
	ListRoom(user *model.User) ([]model.Room, error)
	ListPayment(room *model.Room) ([]model.PaymentDestination, error)
	FindById(roomId uint) (*model.Room, error)
}

type RoomService struct {
	roomRepo repository.IRoomRepository
	pdnRepo  repository.IPaymentDestinationRepository
}

func NewRoomService(repo repository.IRoomRepository) IRoomService {
	return &RoomService{
		roomRepo: repo,
	}
}

// CreateRoomReq POST /rooms/create時に使用するDTO
type CreateRoomReq struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (r RoomService) CreateRoom(roomReq CreateRoomReq, user *model.User) (*model.Room, error) {
	if roomReq.Name == "" {
		return nil, errors.New("部屋名を指定してください")
	}

	var room model.Room
	if roomReq.Password != "" {
		hashed := utils.Sign(roomReq.Password)
		if hashed == "" {
			return nil, errors.New("部屋の作成に失敗しました")
		}
		room.Password = hashed
	}
	room.Name = roomReq.Name
	room.Users = append(room.Users, user)

	created, err := r.roomRepo.CreateRoom(&room)
	if err != nil {
		return nil, errors.New("部屋の作成に失敗しました")
	}

	return created, nil
}

type JoinRoomReq struct {
	Id       uint   `json:"id"`
	Password string `json:"password"`
}

func (r RoomService) JoinRoom(roomReq JoinRoomReq, user *model.User) (*model.Room, error) {
	// idが存在するか確認
	room, err := r.roomRepo.FindById(roomReq.Id)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, errors.New("部屋が見つかりません")
	}

	// 存在したならばパスワードがnilか調べ、存在するならVerify
	if room.Password != "" {
		if !utils.Verify(roomReq.Password, room.Password) {
			return nil, errors.New("合言葉が違います")
		}
	}

	// Idが存在すればr.roomRepository.Join()
	err = r.roomRepo.JoinRoom(room, user)
	if err != nil {
		return nil, err
	}

	// joinできれば部屋を返す
	return room, nil
}

func (r RoomService) ListRoom(user *model.User) ([]model.Room, error) {
	rooms, err := r.roomRepo.FindByUserId(user.ID)
	if err != nil {
		println(err.Error())
		return nil, errors.New("部屋一覧の取得に失敗しました。")
	}
	return rooms, nil
}

func (r RoomService) ListPayment(room *model.Room) ([]model.PaymentDestination, error) {
	paymentDestinations, err := r.pdnRepo.FindByRoomId(room.ID)
	if err != nil {
		return nil, errors.New("支払い先一覧の取得に失敗しました。")
	}
	return paymentDestinations, nil
}

func (r RoomService) FindById(roomId uint) (*model.Room, error) {
	room, err := r.roomRepo.FindById(roomId)
	if err != nil {
		return nil, err
	}
	return room, nil
}
