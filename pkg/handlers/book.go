package handlers

import (
	"net/http"
	"strconv"

	"github.com/2hard4me/pkg/models"
	"github.com/gin-gonic/gin"
)

type getAllResponse struct {
	Data []models.Books `json:"data"`
}

func (h *Handler) GetBatch(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	batch, err := h.services.Book.GetBatch(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllResponse{
		Data: batch,
	})
}