package southwest

//func NewSearchModel(srcAirport, destAirport, beginDate, endDate string, numPassengers int, roundTripFlight bool) *searchModel {
//	s := new(searchModel)
//
//	// Set the criteria
//	srcCriteria := criteria{
//		Stations: stations{
//			OriginStationCodes:      []string{srcAirport},
//			DestinationStationCodes: []string{destAirport},
//		},
//		Dates: dates{
//			BeginDate: beginDate,
//			EndDate:   beginDate,
//		},
//	}
//
//	if roundTripFlight {
//		destCritera := criteria{
//			Stations: stations{
//				OriginStationCodes:      []string{destAirport},
//				DestinationStationCodes: []string{srcAirport},
//			},
//			Dates: dates{
//				BeginDate: endDate,
//				EndDate:   endDate,
//			},
//		}
//		s.Criteria = []criteria{srcCriteria, destCritera}
//	} else {
//		srcCriteria.Dates.EndDate = endDate
//		s.Criteria = []criteria{srcCriteria}
//	}
//	//Set passengers
//	s.Passengers = passengers{Types: []passengerTypes{{
//		Type:  "ADT",
//		Count: numPassengers,
//	}}}
//
//	//Set codes
//	s.Codes = codes{Currency: "USD"}
//
//	s.FareFilters = fareFilters{
//		Loyalty:      "MonetaryOnly",
//		Types:        []string{},
//		ClassControl: 1,
//	}
//
//	s.TaxesAndFees = "TaxesAndFees"
//
//	return s
//}