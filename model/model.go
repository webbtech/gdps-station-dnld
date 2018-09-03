package model

import "time"

// FuelTypes var
var FuelTypes = [4]string{"NL", "SNL", "DSL", "CDSL"}

var x = [5]float64{98, 93, 77, 82, 83}

// RequestInput struct
type RequestInput struct {
	Date      string `json:"date"`
	StationID string `json:"stationID"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Status  int    `json:"status"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Request struct
type Request struct {
	Date      time.Time
	StationID string
}

// ======================== Qraphql Structs ================================ //

// FuelSales struct
type FuelSales struct {
	Date   time.Time
	Report struct {
		StationSales []struct {
			Date  int64
			Sales map[string]float64
		}
		SalesSummary map[string]float64
		SalesTotal   float64
	} `json:"fuelSaleMonth"`
	Station struct {
		ID   string
		Name string
	}
	FuelTypes []string
}

// FuelDelivery struct
type FuelDelivery struct {
	Date   time.Time
	Report struct {
		Deliveries []struct {
			Date int64
			Data map[string]int32
		}
		DeliverySummary map[string]float64
		FuelTypes       []string
	} `json:"fuelDeliveryReport"`
	Station struct {
		ID   string
		Name string
	}
}

// OverShortMonth struct
type OverShortMonth struct {
	Date   time.Time
	Report struct {
		OverShort []struct {
			Date int64
			Data map[string]struct {
				TankLitres float64
				OverShort  float64
			}
		}
		OverShortSummary map[string]float64
		FuelTypes        []string
	} `json:"dipOSMonthReport"`
	Station struct {
		ID   string
		Name string
	}
}

// OverShortAnnual struct
type OverShortAnnual struct {
	Date   time.Time
	Report struct {
		Months    map[string]map[string]float64
		Summary   map[string]float64
		FuelTypes []string
		Year      int
	} `json:"dipOSAnnualReport"`
	Station struct {
		ID   string
		Name string
	}
}
