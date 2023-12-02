package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"net/http"
	"strconv"
	dto "waroka/controllers/DTO"
	"waroka/middleware"
	"waroka/model"
	"waroka/services"
)

type RoomController struct {
	userService    services.IUserService
	roomService    services.IRoomService
	paymentService services.IPaymentService
}

func NewRoomController(
	u services.IUserService,
	r services.IRoomService,
	p services.IPaymentService,
) *RoomController {
	return &RoomController{
		userService:    u,
		roomService:    r,
		paymentService: p,
	}
}

func (r *RoomController) RegisterRoutes(e *echo.Echo) {
	room := e.Group("/room")

	room.Use(middleware.AuthenticateMiddleware)

	room.POST("/create", r.createRoom)
	room.POST("/join", r.joinRoom)
	room.GET("/list", r.listAll)
	room.GET("/:id", r.byId)
	room.GET("/:id/users", r.users)
}

type RoomDetail struct {
	Id          uint   `json:"id"`
	RoomName    string `json:"room_name"`
	TotalAmount int    `json:"total_amount"`
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
			"message": "bad" + err.Error(),
		})
	}

	roomList := lo.Map(rooms, func(rm model.Room, idx int) RoomDetail {
		totalAmount := lo.Reduce(rm.PaymentDestinations, func(acc int, dest *model.PaymentDestination, idx int) int {
			return acc + dest.TotalAmount
		}, 0)

		rd := RoomDetail{
			Id:          rm.ID,
			RoomName:    rm.Name,
			TotalAmount: totalAmount,
		}
		return rd
	})

	return c.JSON(http.StatusOK, roomList)
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
			"message": err,
		})
	}
	println(room.Users)

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

	var userDtos []dto.UserDto
	for _, user := range room.Users {
		userDtos = append(userDtos, dto.UserDto{
			Id:   user.ID,
			Name: user.Name,
		})
	}

	// すべての支払い金額を合計
	var roomResponse dto.RoomResponse
	roomResponse.Id = room.ID
	roomResponse.Name = room.Name
	roomResponse.Users = userDtos
	roomResponse.PaymentDestinations = room.PaymentDestinations
	roomResponse.TotalAmount = totalAmount
	roomResponse.TotalActiveAmount = totalActiveAmount

	return c.JSON(http.StatusOK, roomResponse)
}

func (r *RoomController) users(c echo.Context) error {
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"roomid": roomId})
}
