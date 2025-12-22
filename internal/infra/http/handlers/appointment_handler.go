package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"wimed/internal/application/dto"
	"wimed/internal/application/usecase"
	"wimed/internal/infra/http/restError"
)

type AppointmentHandler struct {
	BookUC *usecase.BookAppointment
}

type bookAppointmentRequest struct {
	AppointmentID   string `json:"appointment_id"`
	PaymentID       string `json:"payment_id"`
	SlotID          string `json:"slot_id"`
	DoctorID        string `json:"doctor_id"`
	PatientID       string `json:"patient_id"`
	PriceCents      int64  `json:"price_cents"`
	PaymentProvider string `json:"payment_provider"`
	ExternalRef     string `json:"external_ref"`
}

type bookAppointmentResponse struct {
	AppointmentID string `json:"appointment_id"`
	PaymentID     string `json:"payment_id"`
	Status        string `json:"status"`
}

func NewAppointmentHandler(book *usecase.BookAppointment) *AppointmentHandler {
	return &AppointmentHandler{BookUC : book}
}

func (h *AppointmentHandler) Book(c *gin.Context) {
	var req bookAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		re := restError.NewBadRequestError("invalid json body")
		c.JSON(re.Code, re)
		return
	}

	causes := make([]restError.Causes, 0)

	if strings.TrimSpace(req.AppointmentID) == "" {
		causes = append(causes, restError.Causes{Field: "appointment_id", Message: "is required"})
	}
	if strings.TrimSpace(req.PaymentID) == "" {
		causes = append(causes, restError.Causes{Field: "payment_id", Message: "is required"})
	}
	if strings.TrimSpace(req.SlotID) == "" {
		causes = append(causes, restError.Causes{Field: "slot_id", Message: "is required"})
	}
	if strings.TrimSpace(req.DoctorID) == "" {
		causes = append(causes, restError.Causes{Field: "doctor_id", Message: "is required"})
	}
	if strings.TrimSpace(req.PatientID) == "" {
		causes = append(causes, restError.Causes{Field: "patient_id", Message: "is required"})
	}
	if req.PriceCents < 0 {
		re := restError.NewBadRequestValidationError("validation error", causes)
		c.JSON(re.Code, re)
		return
	}

	out, err := h.BookUC.Execute(c.Request.Context(), dto.BookAppointmentInput{
		AppointmentID: 		req.AppointmentID,
		PaymentID: 			req.PaymentID,
		SlotID: 			req.SlotID,
		DoctorID: 			req.DoctorID,
		PatientID: 			req.PatientID,
		PriceCents: 		req.PriceCents,
		PaymentProvider: 	req.PaymentProvider,
		ExternalRef: 		req.ExternalRef,
	})
	if err != nil {
		re := mapBookError(err)
		c.JSON(re.Code, re)
		return
	}

	c.JSON(http.StatusCreated, bookAppointmentResponse{
		AppointmentID: 	out.AppointmentID,
		PaymentID: 		out.PaymentID,
		Status: 		out.Status,
	})
}

func mapBookError(err error) *restError.RestErr {
	msg := err.Error()

	if strings.Contains(msg, "not found") || strings.Contains(msg, "patient not found") || strings.Contains(msg, "slot not found") {
		return restError.NewNotFoundError(msg)
	}

	if strings.Contains(msg, "only book an available slot") ||
		strings.Contains(msg, "slot is not available") ||
		strings.Contains(msg, "already exists for this slot") {
		return restError.NewRestErr(msg, "conflict", http.StatusConflict, nil)
	}

	if errors.Is(err, contextCanceled()) {
		return restError.NewRestErr("request canceled", "canceled", http.StatusRequestTimeout, nil)
	}

	return restError.NewInternalServerError("internal server error")
}

func contextCanceled() error {
	return errors.New("Context canceled")
}