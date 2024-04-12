package DispenserStrategies

import (
	"errors"
	"mcs/TestDesign/Strategies"
)

type DispenserFactory struct{}

func (f *DispenserFactory) CreateStrategy(identifier string) (Strategies.Strategy, error) {
	switch identifier {
	case "v1":
		return &DispenserV1{}, nil
	case "v2":
		return &DispenserV2{}, nil
	case "v3":
		return &DispenserV3{}, nil
	case "v4":
		return &DispenserV4{}, nil
	default:
		return nil, errors.New("unknown compressor strategy")
	}
}
