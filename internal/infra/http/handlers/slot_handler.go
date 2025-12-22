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
	ListUC 		*usecase.ListAvailableSlots
	GenerateUC 	*usecase.GenerateSlots
}

func NewSlotHandler(list *usecase.ListAvailableSlots, gen *usecase.GenerateSlots) *SlotHandler {
	return &SlotHandler{ListUC: list, GenerateUC: gen}
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

type generateSlotsRequest struct {
	From           string `json:"from"`
	To             string `json:"to"`
	SessionMinutes int    `json:"session_minutes"`
	WorkStartHour  int    `json:"work_start_hour"`
	WorkEndHour    int    `json:"work_end_hour"`
	TimeZone       string `json:"timezone"`
}

func (h *SlotHandler) Generate(c *gin.Context) {
	doctorID := c.Param("doctor_id")

	var req generateSlotsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		re := restError.NewBadRequestError("invalid json body")
		c.JSON(re.Code, re)
		return
	}

	from, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		re := restError.NewBadRequestError("invalid 'from' (RFC3339 required)")
		c.JSON(re.Code, re)
		return
	}

	to, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		re := restError.NewBadRequestError("invalid 'to' (RFC3339 required)")
		c.JSON(re.Code, re)
		return
	}

	loc := time.UTC
	if req.TimeZone != "" {
		l, err := time.LoadLocation(req.TimeZone)
		if err != nil {
			re := restError.NewBadRequestError("invalid timezone")
			c.JSON(re.Code, re)
			return
		}
		loc = l	
	}
	
	out, err := h.GenerateUC.Execute(c.Request.Context(), usecase.GenerateSlotsInput{
		DoctorID: 		doctorID,
		From: 			from,
		To:				to,
		SessionMinutes: req.SessionMinutes,
		WorkStartHour: 	req.WorkStartHour,
		WorkEndHour: 	req.WorkEndHour,
		TimeZone: 		loc,
	})
	if err != nil {
		re := restError.NewBadRequestError(err.Error())
		c.JSON(re.Code, re)
		return
	}

	c.JSON(http.StatusCreated, out)
}