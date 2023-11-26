package room

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"waroka/services"
)

type PaymentController struct {
	userService    services.IUserService
	roomService    services.IRoomService
	paymentService services.IPaymentService
}

func NewPaymentController(
	u services.IUserService,
	r services.IRoomService,
	p services.IPaymentService,
) *PaymentController {
	return &PaymentController{
		userService:    u,
		roomService:    r,
		paymentService: p,
	}
}

func (r *PaymentController) RegisterRoutes(e *echo.Echo) {
	pay := e.Group("/payment")

	pay.GET("/:id", r.getById)
}

func (r *PaymentController) getById(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
