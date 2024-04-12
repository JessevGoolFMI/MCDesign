package Strategies

type IStrategyFactory interface {
	CreateStrategy(identifier string) (IStrategy, error)
}
