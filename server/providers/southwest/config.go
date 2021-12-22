package southwest

import "sync"

type ApiConfig struct {
	sync.RWMutex
	Headers map[string]string
}

var ApiInfo = &ApiConfig{
	RWMutex: sync.RWMutex{},
	Headers: nil,
}

