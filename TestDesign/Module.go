package TestDesign

import (
	"errors"
	"fmt"
	"mcs/TestDesign/Strategies"
	"mcs/TestDesign/Strategies/CompressorStrategies"
	"time"
)

/*
This file defines the BaseModule and CompressorModule structs, along with their associated methods and factory interfaces, in the context of a mediator pattern implementation. The BaseModule struct represents a generic module that can subscribe to values, unsubscribe from values, and publish values. The CompressorModule extends the BaseModule with additional functionality, showcasing how modules can be specialized for specific purposes.

Key Features:

    BaseModule Struct: The BaseModule struct encapsulates the core functionality of a module, including its identifier (id), a notifier callback (notifier), and a reference to the mediator (Mediator). This design allows modules to communicate with each other through the mediator, subscribing to and publishing values.

    NotificationCallback: A type alias for a function that takes a value name and a value, allowing modules to define custom behavior upon receiving notifications about value changes.

    NotifySubscriber Method: This method is used to notify the module about a value change. If a custom notifier callback is set, it is invoked; otherwise, a default message is printed to the console.

    Subscription and Publishing Methods: The BaseModule provides methods to subscribe to and unsubscribe from values (SubscribeToTopic and UnsubscribeFromTopic), as well as to publish values (PublishToTopic). These methods utilize the mediator to send commands for subscription, unsubscription, and publication.

    CompressorModule Struct: The CompressorModule extends the BaseModule with an additional field (specialValue), demonstrating how modules can be specialized for specific purposes. It inherits all methods from the BaseModule struct, including subscription, unsubscription, and publishing methods.

    IModuleFactory Interface and DefaultModuleFactory Implementation: These components provide a factory pattern for creating instances of BaseModule and CompressorModule. The IModuleFactory interface defines methods for creating modules, and the DefaultModuleFactory provides a default implementation of these methods. This design supports the creation of modules without directly instantiating the structs, promoting code flexibility and maintainability.

This file exemplifies the use of the mediator pattern in Go, focusing on the creation and interaction of modules within a system. It demonstrates how modules can be specialized for specific purposes while maintaining a common interface for communication and value management. The factory pattern for module creation further enhances the design, making it easier to instantiate modules and specialized modules as needed.
*/

type State int

const (
	// InitState is the initial state of the module.
	InitState State = iota
	// RunningState is the state after the module has transitioned from the initial state.
	RunningState
	ShutdownState
	ErrorState
)

type Module interface {
	Execute() (interface{}, error)
}

// BaseModule struct
type BaseModule struct {
	id       string
	notifier NotificationCallback
	Mediator IMediator
	state    State
	stopChan chan byte
}

func NewModule(id string, controller IMediator) *BaseModule {
	module := &BaseModule{
		id:       id,
		Mediator: controller,
		state:    InitState, // Initialize the state to InitState
		stopChan: make(chan byte),
	}
	return module
}

func (m *BaseModule) GetId() string {
	return m.id
}

func (m *BaseModule) GetState() State {
	return m.state
}

func (m *BaseModule) SetState(state State) {
	m.state = state
}

func (m *BaseModule) TransitionToRunning() {
	if m.state == InitState {
		fmt.Printf("BaseModule %s is transitioning to idle state.\n", m.id)
		m.state = RunningState
		go m.startBackgroundProcess() // Start the background process in a goroutine
	} else {
		fmt.Printf("BaseModule %s is already in idle state.\n", m.id)
	}
}

func (m *BaseModule) ResolveError() {
	m.resolveErrorAndResume()
}

func (m *BaseModule) startBackgroundProcess() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopChan:
			fmt.Printf("BaseModule %s background process stopped.\n", m.id)
			return
		case <-ticker.C:
			// Check if the module is in ErrorState
			if m.state == ErrorState {
				fmt.Printf("BaseModule %s is in error state, pausing background process.\n", m.id)
				// Pause the process by waiting for a signal to resume
				<-m.stopChan // This will block until the module is signaled to resume
				continue     // Skip the rest of the loop iteration
			}
			number, err := fetchRandomNumber()
			if err != nil {
				m.SetState(ErrorState) // Transition to ErrorState on error
				fmt.Printf("BaseModule %s encountered an error: %v\n", m.id, err)
				return
			}
			m.PublishToTopic("randomInt", number)
		}
	}
}

func (m *BaseModule) resolveErrorAndResume() {
	// Hypothetical error resolution logic here
	fmt.Printf("BaseModule %s error resolved, resuming background process.\n", m.id)
	m.SetState(RunningState) // Transition back to RunningState
	m.stopChan <- 1          // Send a signal to resume the background process
}

func (m *BaseModule) StopBackgroundProcess() {
	if m.state == RunningState {
		m.state = ShutdownState // Transition to Shutdown state
		m.stopChan <- 0         // Signal the background process to stop
	}
}

type NotificationCallback = func(valueName string, value interface{})

// NotifySubscriber notifies the module about a value change
func (m *BaseModule) NotifySubscriber(valueName string, value interface{}) {
	if m.notifier != nil {
		m.notifier(valueName, value)
	} else {
		fmt.Printf("BaseModule %s received %s value update: %v\n", m.id, valueName, value)
	}
}

func (m *BaseModule) SetNotificationCallback(callback NotificationCallback) {
	m.notifier = callback
}

func (m *BaseModule) SubscribeToTopic(topic string, target string) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&SubscribeCommand{subscriberID: m.id, publisherID: target, topic: topic}, target)
	}
}

func (m *BaseModule) UnsubscribeFromTopic(topic string, target string) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&UnsubscribeCommand{subscriberID: m.id, publisherID: target, topic: topic}, target)
	}
}

func (m *BaseModule) PublishToTopic(topic string, value interface{}) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&PublishValueCommand{publisherID: m.id, topic: topic, value: value}, m.id)
	}
}

type CompressorModule struct {
	*BaseModule
	specialValue interface{}
	compressor   func(interface{}) (interface{}, error)
}

func (cm *CompressorModule) Execute() (interface{}, error) {
	if cm.compressor == nil {
		return nil, errors.New("execute wasn't set correctly")
	}
	return cm.compressor(cm.specialValue)
}

func (cm *CompressorModule) SetStrategy(strategy CompressorStrategies.CompressorFunc) {
	cm.compressor = strategy
}

func NewCompressorModule(id string, controller IMediator, specialValue interface{}) *CompressorModule {
	return &CompressorModule{
		BaseModule:   NewModule(id, controller),
		specialValue: specialValue,
	}
}

type DispenserModule struct {
	*BaseModule
	specialValue interface{}
	dispenser    Strategies.Strategy
}

func (dm *DispenserModule) Execute() (interface{}, error) {
	if dm.dispenser == nil {
		return nil, errors.New("execute wasn't set correctly")
	}
	return dm.dispenser.Execute(dm.specialValue)
}

func (dm *DispenserModule) SetStrategy(strategy Strategies.Strategy) {
	dm.dispenser = strategy
}

func NewDispenserModule(id string, controller IMediator, specialValue interface{}) *DispenserModule {
	return &DispenserModule{
		BaseModule:   NewModule(id, controller),
		specialValue: specialValue,
	}
}

// IModuleFactory interface
type IModuleFactory interface {
	CreateModule(id string, controller IMediator) *BaseModule
	CreateCompressorModule(id string, controller IMediator, specialValue interface{}) *CompressorModule
	CreateDispenserModule(id string, controller IMediator, specialValue interface{}) *DispenserModule
}

// DefaultModuleFactory implementation
type DefaultModuleFactory struct{}

func (f *DefaultModuleFactory) CreateModule(id string, controller IMediator) *BaseModule {
	return NewModule(id, controller)
}

func (f *DefaultModuleFactory) CreateCompressorModule(id string, controller IMediator, specialValue interface{}) *CompressorModule {
	return NewCompressorModule(id, controller, specialValue)
}

func (f *DefaultModuleFactory) CreateDispenserModule(id string, controller IMediator, specialValue interface{}) *DispenserModule {
	return NewDispenserModule(id, controller, specialValue)
}
