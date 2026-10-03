package domain

type Car struct {
	CarUID             string
	Brand              string
	Model              string
	RegistrationNumber string
	Power              *int
	Type               string
	Price              int
	Available          bool
}

type CarPage struct {
	Page          int
	PageSize      int
	TotalElements int
	Items         []Car
}
