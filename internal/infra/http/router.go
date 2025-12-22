package http

import (
	"wimed/internal/infra/http/handlers"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *handlers.AppointmentHandler, slots *handlers.SlotHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.POST("/appointments/book", h.Book)
	r.GET("/doctors/:doctor_id/slots", slots.ListAvailableByDoctor)

	return r
}