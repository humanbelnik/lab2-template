package domain

import (
	"errors"
	"time"
)

type RentalStatus string

const (
	RentalStatusInProgress RentalStatus = "IN_PROGRESS"
	RentalStatusFinished   RentalStatus = "FINISHED"
	RentalStatusCanceled   RentalStatus = "CANCELED"
)

type Rental struct {
	ID         int64
	RentalUID  string
	Username   string
	PaymentUID string
	CarUID     string
	DateFrom   time.Time
	DateTo     time.Time
	Status     RentalStatus
}

var (
	ErrRentalNotFound = errors.New("rental not found")
)

const DateLayout = "2006-01-02"
