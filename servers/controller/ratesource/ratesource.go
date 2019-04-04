package ratesource

import "github.com/briferz/usdmxn/model"

type Interface interface {
	GetSource() (model.Provider, error)
}
