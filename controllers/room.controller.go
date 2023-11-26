package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	"waroka/middleware"
	"waroka/model"
	"waroka/services"
)

type RoomController struct {
	userService services.IUserService
	roomService services.IRoomService
}

func NewRoomController(
	u services.IUserService,
	r services.IRoomService,
	p services.IPaymentService,
) *RoomController {
	return &RoomController{
		userService: u,
		roomService: r,
	}
}

func (r *RoomController) RegisterRoutes(e *echo.Echo) {
	room := e.Group("/room")

	room.Use(middleware.AuthenticateMiddleware)

	room.POST("/create", r.createRoom)
	room.POST("/join", r.joinRoom)
	room.GET("/list", r.listAll)
	room.GET("/:id", r.byId)
}

type RoomDetail struct {
	Id       uint   `json:"id"`
	RoomName string `json:"room_name"`
}

func (r *RoomController) createRoom(c echo.Context) error {
	var req services.CreateRoomReq
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "リクエストの形式が不正です。",
		})
	}

	userId := c.Get("user").(uint)
	user, err := r.userService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ログイン中のユーザーを取得できませんでした。",
		})
	}

	room, err := r.roomService.CreateRoom(req, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err,
		})
	}

	return c.JSON(http.StatusCreated, RoomDetail{
		Id:       room.ID,
		RoomName: room.Name,
	})
}

func (r *RoomController) joinRoom(c echo.Context) error {
	var req services.JoinRoomReq
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "リクエストの形式が不正です。",
		})
	}

	// ログイン中のユーザーを取得
	userId := c.Get("user").(uint)
	user, err := r.userService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ログイン中のユーザーを取得できませんでした。",
		})
	}

	// 入室させる
	room, err := r.roomService.JoinRoom(req, user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err,
		})
	}

	return c.JSON(http.StatusCreated, RoomDetail{
		Id:       room.ID,
		RoomName: room.Name,
	})
}

func (r *RoomController) listAll(c echo.Context) error {
	// ログイン中のユーザーを取得
	userId := c.Get("user").(uint)
	user, err := r.userService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ログインしてください。",
		})
	}

	rooms, err := r.roomService.ListRoom(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err,
		})
	}

	roomList := lo.Map(rooms, func(rm model.Room, idx int) RoomDetail {
		rd := RoomDetail{
			Id:       rm.ID,
			RoomName: rm.Name,
		}
		return rd
	})

	return c.JSON(http.StatusOK, roomList)
}

type RoomResponse struct {
	Name                string
	Id                  uint
	Users               []*model.User
	PaymentDestinations []*model.PaymentDestination
	TotalAmount         int
	TotalActiveAmount   int
}

func (r *RoomController) byId(c echo.Context) error {
	userId := c.Get("user").(uint)
	_, err := r.userService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "ログインしてください。",
		})
	}

	roomId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "IDを読み取れませんでした。",
		})
	}

	room, err := r.roomService.FindById(uint(roomId))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "IDを読み取れませんでした。",
		})
	}

	destinations := room.PaymentDestinations
	var totalAmount int
	var totalActiveAmount int
	// すべての支払い先を取得
	for _, destination := range destinations {
		if destination.DeletedAt.Valid {
			details := destination.PaymentDetails
			for _, detail := range details {
				if detail.DeletedAt.Valid {
					totalActiveAmount += detail.Amount
				}
			}
		}
		totalAmount += destination.TotalAmount
	}

	// すべての支払い金額を合計
	var roomResponse RoomResponse
	roomResponse.Id = room.ID
	roomResponse.Name = room.Name
	roomResponse.Users = room.Users
	roomResponse.PaymentDestinations = room.PaymentDestinations
	roomResponse.TotalAmount = totalAmount
	roomResponse.TotalActiveAmount = totalActiveAmount

	return c.JSON(http.StatusOK, roomResponse)
}
