package shared

type personBlock struct {
	Data struct {
		Persons []Person `json:"persons"`
	} `json:"data"`
}

type Person struct {
	ID                 int                `json:"id"`
	FullName           string             `json:"fullName"`
	FirstName          string             `json:"firstName"`
	LastName           string             `json:"lastName"`
	UnavailabilityCode UnavailabilityCode `json:"unavailabilityCode"`
}

type UnavailabilityCode struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Color       string `json:"color"`
	IsAvailable bool   `json:"isAvailable"`
}

type interventionBlock struct {
	Data []Intervention `json:"data"`
}

type Intervention struct {
	ID       int `json:"id"`
	Status   int `json:"status"`
	Priority int `json:"priority"`
	Type     struct {
		Code                 string `json:"code"`
		Description          string `json:"description"`
		Fastest              bool   `json:"fastest"`
		Qualified            bool   `json:"qualified"`
		Urgency              int    `json:"urgency"`
		AllowLocalAllocation bool   `json:"allowLocalAllocation"`
		Category             struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Icon    string `json:"icon"`
			IsRound bool   `json:"isRound"`
		}
		Icon string `json:"icon"`
	} `json:"type"`
}

type vehicleBlock struct {
	Data []Vehicle `json:"data"`
}

type Vehicle struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Station struct {
		ID   int    `json:"id"`
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"station"`
	Vehicle struct {
		ID           int    `json:"id"`
		Code         string `json:"code"`
		LicensePlate string `json:"licensePlate"`
	} `json:"vehicle"`
}
