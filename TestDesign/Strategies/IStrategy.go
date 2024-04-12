package Strategies

type IStrategy interface {
	Execute(value interface{}) (interface{}, error)
}
