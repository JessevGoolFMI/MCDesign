package CompressorStrategies

import (
	"errors"
	"mcs/TestDesign/Strategies"
)

type CompressorFactory struct{}

func (f *CompressorFactory) CreateStrategy(identifier string) (Strategies.Strategy, error) {
	switch identifier {
	case "v1":
		return &CompressorV1{}, nil
	case "v2":
		return &CompressorV2{}, nil
	case "v3":
		return &CompressorV3{}, nil
	case "v4":
		return &CompressorV4{}, nil
	default:
		return nil, errors.New("unknown compressor strategy")
	}
}
