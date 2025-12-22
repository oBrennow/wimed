package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"wimed/internal/application/ports"
)

type ListAvailableSlots struct {
	TxManager ports.TxManager
	Slots     ports.SlotReadRepository

	Now func() time.Time
}

type ListAvailableSlotsInput struct {
	DoctorID string
	From     time.Time
	To       time.Time
	Limit    int
}

type SlotItem struct {
	ID   	   	string
	StartedAt	time.Time
	EndedAt   	time.Time
}

type ListAvailableSlotsOutput struct {
	DoctorID string
	Slots    []SlotItem
}

func (uc *ListAvailableSlots) Execute(ctx context.Context, in ListAvailableSlotsInput) (*ListAvailableSlotsOutput, error) {
	if strings.TrimSpace(in.DoctorID) == "" {
		return nil, errors.New("doctor_id is required")
	}
	if in.From.IsZero() || in.To.IsZero() || !in.From.Before(in.To) {
		return nil, errors.New("invalid date range")
	}

	tx, err := uc.TxManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	slots, err := uc.Slots.ListAvailableByDoctor(ctx , tx , in.DoctorID, in.From, in.To, in.Limit)
	if err != nil {
		return nil, err
	}

	out := make([]SlotItem, 0, len(slots))
	for _, s := range slots {
		out = append(out, SlotItem{
			ID: 		s.ID(),
			StartedAt: 	s.StartedAt(),
			EndedAt: 	s.EndedAt(),
		})
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &ListAvailableSlotsOutput{
		DoctorID: 	in.DoctorID,
		Slots: 		out,
	}, nil
}