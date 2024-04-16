package CompressorStrategies

import (
	"errors"
)

type CompressorStrategyFactory struct{}

func (f *CompressorStrategyFactory) CreateStrategy(identifier string) (CompressorFunc, error) {
	switch identifier {
	case "v1":
		return CompressorV1, nil
	case "v2":
		return CompressorV2, nil
	case "v3":
		return CompressorV3, nil
	case "v4":
		return CompressorV4, nil
	// And so on for other strategies...
	default:
		return nil, errors.New("unknown compressor strategy")
	}
}
