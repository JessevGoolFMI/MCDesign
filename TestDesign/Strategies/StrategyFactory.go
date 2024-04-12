package Strategies

type StrategyFactory interface {
	CreateStrategy(identifier string) (Strategy, error)
}
