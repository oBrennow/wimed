package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"wimed/internal/application/usecase"
	"wimed/internal/infra/http/restError"
)

type SlotHandler struct {
	ListUC *usecase.ListAvailableSlots
}

func NewSlotHandler(list *usecase.ListAvailableSlots) *SlotHandler {
	return &SlotHandler{ListUC: list}
}

func (h *SlotHandler) ListAvailableByDoctor(c *gin.Context) {
	doctorID := c.Param("doctor_id")

	fromStr := c.Query("from")
	toStr := c.Query("to")
	limitStr := c.DefaultQuery("limit", "50")

	from, err := time.Parse(time.RFC3339, fromStr)
	if err != nil {
		re := restError.NewBadRequestError("invalid 'from' (RFC3339 required)")
		c.JSON(re.Code, re)
		return
	}
	to, err := time.Parse(time.RFC3339, toStr)
	if err != nil {
		re := restError.NewBadRequestError("invalid 'to' (RFC3339 required)")
		c.JSON(re.Code, re)
		return
	}

	limit, _ := strconv.Atoi(limitStr)

	out, err := h.ListUC.Execute(c.Request.Context(), usecase.ListAvailableSlotsInput{
		DoctorID: doctorID,
		From:     from,
		To:       to,
		Limit:    limit,
	})
	if err != nil {
		re := restError.NewBadRequestError(err.Error())
		c.JSON(re.Code, re)
		return
	}

	c.JSON(http.StatusOK, out)
}
