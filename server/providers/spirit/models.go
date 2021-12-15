package spirit

type passengerTypes struct {
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type passengers struct {
	Types []passengerTypes `json:"types"`
}

type codes struct {
	Currency string `json:"currency"`
}

type fareFilters struct {
	Loyalty      string   `json:"loyalty"`
	Types        []string `json:"types"`
	ClassControl int      `json:"classControl"`
}

type searchModel struct {
	Criteria     []criteria  `json:"criteria"`
	Passengers   passengers  `json:"passengers"`
	Codes        codes       `json:"codes"`
	FareFilters  fareFilters `json:"fareFilters"`
	TaxesAndFees string      `json:"taxesAndFees"`
}

func NewSearchModel(srcAirport, destAirport, beginDate, endDate string, numPassengers int, roundTripFlight bool) *searchModel {
	s := new(searchModel)

	// Set the criteria
	srcCriteria := criteria{
		Stations: stations{
			OriginStationCodes:      []string{srcAirport},
			DestinationStationCodes: []string{destAirport},
		},
		Dates: dates{
			BeginDate: beginDate,
			EndDate:   beginDate,
		},
	}

	if roundTripFlight {
		destCritera := criteria{
			Stations: stations{
				OriginStationCodes:      []string{destAirport},
				DestinationStationCodes: []string{srcAirport},
			},
			Dates: dates{
				BeginDate: endDate,
				EndDate:   endDate,
			},
		}
		s.Criteria = []criteria{srcCriteria, destCritera}
	} else {
		srcCriteria.Dates.EndDate = endDate
		s.Criteria = []criteria{srcCriteria}
	}
	//Set passengers
	s.Passengers = passengers{Types: []passengerTypes{{
		Type:  "ADT",
		Count: numPassengers,
	}}}

	//Set codes
	s.Codes = codes{Currency: "USD"}

	s.FareFilters = fareFilters{
		Loyalty:      "MonetaryOnly",
		Types:        []string{},
		ClassControl: 1,
	}

	s.TaxesAndFees = "TaxesAndFees"

	return s
}

type stations struct {
	OriginStationCodes      []string `json:"originStationCodes"`
	DestinationStationCodes []string `json:"destinationStationCodes"`
}

type dates struct {
	BeginDate string `json:"beginDate"`
	EndDate   string `json:"endDate"`
}

type criteria struct {
	Stations stations `json:"stations"`
	Dates    dates    `json:"dates"`
}
