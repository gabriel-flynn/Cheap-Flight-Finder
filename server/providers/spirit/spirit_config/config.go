package spirit_config

import "sync"

type ApiConfig struct {
	sync.RWMutex
	AuthToken string
}

var ApiInfo *ApiConfig

type Requirements struct {
	IsSaversClub bool
}

var FlightRequirements *Requirements

func init() {
	FlightRequirements = &Requirements{
		IsSaversClub: false,
	}
}
