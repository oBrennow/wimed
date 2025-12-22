package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"wimed/internal/application/ports"
	"wimed/internal/domain/availabilityDomain"
)

type GenerateSlots struct {
	TxManager 	ports.TxManager
	Slots		ports.SlotWriteRepoitory
	Now			func() time.Time
	NewID 		func(prefix string) string
}

type GenerateSlotsInput struct {
	DoctorID		string
	From			time.Time
	To				time.Time
	SessionMinutes	int
	WorkStartHour	int
	WorkEndHour		int
	TimeZone		*time.Location
}

type GenerateSlotsOutput struct {
	Created int `json:"created"`
}


func (uc *GenerateSlots) Execute(ctx context.Context, in GenerateSlotsInput) (*GenerateSlotsOutput, error) {
	if strings.TrimSpace(in.DoctorID) == "" {
		return nil, errors.New("doctor_id is required")
	}
	if in.From.IsZero() || in.To.IsZero() || !in.From.Before(in.To) {
		return nil, errors.New("invalid date range")
	}
	if in.SessionMinutes <= 0 {
		return nil, errors.New("session_minutes must be > 0")
	}
	if in.WorkStartHour < 0 || in.WorkStartHour > 23 || in.WorkEndHour < 1 || in.WorkEndHour > 24 || in.WorkStartHour >= in.WorkEndHour {
		return nil, errors.New("invalid work hours")
	}
		if in.TimeZone == nil {
		in.TimeZone = time.UTC
	}
	if uc.Now == nil {
		uc.Now = func() time.Time { return time.Now().UTC() }
	}
	if uc.NewID == nil {
		uc.NewID = func(prefix string) string {
			return prefix + "_" + time.Now().UTC().Format("20060102T150405.000000000")
		}
	}

	now := uc.Now()

	from := in.From.In(in.TimeZone)
	to := in.To.In(in.TimeZone)

	step := time.Duration(in.SessionMinutes) * time.Minute
	var slots []*availabilityDomain.SlotDomain

	for day := truncateToDay(from); day.Before(to); day = day.AddDate(0, 0, 1) {
		ws := time.Date(day.Year(), day.Month(), day.Day(), in.WorkStartHour, 0, 0, 0, in.TimeZone)
		we := time.Date(day.Year(), day.Month(), day.Day(), in.WorkEndHour, 0, 0, 0, in.TimeZone)
	

		start := maxTime(ws, from)
		end := minTime(we, to)

		for t := start; t.Add(step).Equal(end)|| t.Add(step).Before(end); t = t.Add(step) {
			id := uc.NewID("slot")

			s, err := availabilityDomain.CreateSlotDomain(
				id,
				in.DoctorID,
				t.UTC(),
				t.Add(step).UTC(),
				availabilityDomain.SlotAvailable,
				now.UTC(),
			)
			if err != nil {
				return nil, err
			}
			slots = append(slots, s)
		}
	}

	tx, err := uc.TxManager.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := uc.Slots.CreateBatch(ctx, tx, slots); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &GenerateSlotsOutput{Created: len(slots)}, nil
}

func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}