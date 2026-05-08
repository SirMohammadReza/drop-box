package http

import (
	"core/internal/auth"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	authProvider auth.Provider
}

func NewHandler(ap auth.Provider) *Handler {
	return &Handler{
		authProvider: ap,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (h *Handler) Login(c *echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	res, err := h.authProvider.Login(
		c.Request().Context(),
		req.PhoneNumber,
		req.Password,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (h *Handler) Register(c *echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	res, err := h.authProvider.Register(
		c.Request().Context(),
		req.PhoneNumber,
		req.Name,
		req.Password,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
