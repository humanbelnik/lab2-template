package domain

import "errors"

type CarType string

const (
	CarTypeSedan    CarType = "SEDAN"
	CarTypeSUV      CarType = "SUV"
	CarTypeMinivan  CarType = "MINIVAN"
	CarTypeRoadster CarType = "ROADSTER"
)

type Car struct {
	ID                 int64
	CarUID             string
	Brand              string
	Model              string
	RegistrationNumber string
	Power              *int
	Price              int
	Type               CarType
	Available          bool
}

var ErrCarNotFound = errors.New("car not found")

type CarPage struct {
	Page          int
	PageSize      int
	TotalElements int
	Items         []Car
}
