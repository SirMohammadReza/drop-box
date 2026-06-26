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
	requestFile, header, err := c.Request().FormFile("file")
	if err != nil {
		return err
	}
	defer requestFile.Close()

	buffer := make([]byte, 512)
	_, err = requestFile.Read(buffer)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "failed to read file")
	}

	fileDTO := file.NewFileInputs{
		Name:     header.Filename,
		Size:     header.Size,
		MimeType: http.DetectContentType(buffer),
	}

	newFile, err := h.fileService.NewFile(c.Request().Context(), fileDTO)

	meta := uploader.Metadata{
		FileID:   newFile.ID.Hex(),
		FileName: newFile.Name,
		Size:     newFile.Size,
		UserUUID: "user-1",
	}

	res, err := h.uploaderProvider.UploadFile(c.Request().Context(), meta, requestFile)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}
