package internal

import "time"

type currentAvailability struct {
	Unavailability struct {
		UnavailabilityCode struct {
			Code          string `json:"code"`
			Name          string `json:"name"`
			Color         string `json:"color"`
			IsUnavailable bool   `json:"isUnavailable"`
			TextColor     string `json:"textColor"`
		} `json:"unavailabilityCode"`
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	} `json:"unavailability"`
	Location struct {
		Company struct {
			Code    string `json:"code"`
			Name    string `json:"name"`
			Station struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"station"`
		} `json:"company"`
		StartDate time.Time `json:"startDate"`
		EndDate   time.Time `json:"endDate"`
	} `json:"location"`
}

type intervention struct {
	InterventionCode struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"interventionCode"`
}

type availabiltyStat struct {
	StartDate        time.Time         `json:"startDate"`
	EndDate          time.Time         `json:"endDate"`
	PercentAvailable float32           `json:"percentAvailable"`
	Periods          []availabiltyStat `json:"periods"`
	AverageMin       int               `json:"avarageMin"`
	AverageMax       int               `json:"averageMax"`
	Level            int               `json:"level"`
}
