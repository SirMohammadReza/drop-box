package http

import (
	"core/internal/file"
	"core/internal/uploader"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	uploaderProvider uploader.Provider
	fileService      file.Provider
}

func NewHandler(up uploader.Provider, fs file.Provider) *Handler {
	return &Handler{
		uploaderProvider: up,
		fileService:      fs,
	}
}

func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/upload-file", h.UploadFile)
}

func (h *Handler) UploadFile(c *echo.Context) error {
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()

	meta := uploader.Metadata{
		FileName: header.Filename,
		Size:     header.Size,
		UserUUID: "user-1",
	}

	res, err := h.uploaderProvider.UploadFile(c.Request().Context(), meta, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
