package ratesource

import "github.com/briferz/usdmxn/servers/model"

type Interface interface {
	GetSource() (model.Provider, error)
}
