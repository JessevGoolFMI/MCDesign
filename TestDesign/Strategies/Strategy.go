package Strategies

type Strategy interface {
	Execute(value interface{}) (interface{}, error)
}
