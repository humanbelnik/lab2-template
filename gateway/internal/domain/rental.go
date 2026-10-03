package domain

type Rental struct {
	RentalUID  string
	Username   string
	PaymentUID string
	CarUID     string
	DateFrom   string
	DateTo     string
	Status     string
}

type RentalDetails struct {
	RentalUID string
	Status    string
	DateFrom  string
	DateTo    string
	Car       Car
	Payment   Payment
}

type CreateRentalInput struct {
	Username string
	CarUID   string
	DateFrom string
	DateTo   string
}

type CreateRentalResult struct {
	RentalUID string
	Status    string
	CarUID    string
	DateFrom  string
	DateTo    string
	Payment   Payment
}
