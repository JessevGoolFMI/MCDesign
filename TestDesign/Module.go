package TestDesign

import (
	"fmt"
	"time"
)

/*
This file defines the Module and SpecialModule structs, along with their associated methods and factory interfaces, in the context of a mediator pattern implementation. The Module struct represents a generic module that can subscribe to values, unsubscribe from values, and publish values. The SpecialModule extends the Module with additional functionality, showcasing how modules can be specialized for specific purposes.

Key Features:

    Module Struct: The Module struct encapsulates the core functionality of a module, including its identifier (id), a notifier callback (notifier), and a reference to the mediator (Mediator). This design allows modules to communicate with each other through the mediator, subscribing to and publishing values.

    NotificationCallback: A type alias for a function that takes a value name and a value, allowing modules to define custom behavior upon receiving notifications about value changes.

    NotifySubscriber Method: This method is used to notify the module about a value change. If a custom notifier callback is set, it is invoked; otherwise, a default message is printed to the console.

    Subscription and Publishing Methods: The Module provides methods to subscribe to and unsubscribe from values (SubscribeToTopic and UnsubscribeFromTopic), as well as to publish values (PublishToTopic). These methods utilize the mediator to send commands for subscription, unsubscription, and publication.

    SpecialModule Struct: The SpecialModule extends the Module with an additional field (specialValue), demonstrating how modules can be specialized for specific purposes. It inherits all methods from the Module struct, including subscription, unsubscription, and publishing methods.

    IModuleFactory Interface and DefaultModuleFactory Implementation: These components provide a factory pattern for creating instances of Module and SpecialModule. The IModuleFactory interface defines methods for creating modules, and the DefaultModuleFactory provides a default implementation of these methods. This design supports the creation of modules without directly instantiating the structs, promoting code flexibility and maintainability.

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

// Module struct
type Module struct {
	id       string
	notifier NotificationCallback
	Mediator IMediator
	state    State
	stopChan chan byte
}

func NewModule(id string, controller IMediator) *Module {
	module := &Module{
		id:       id,
		Mediator: controller,
		state:    InitState, // Initialize the state to InitState
		stopChan: make(chan byte),
	}
	return module
}

func (m *Module) GetId() string {
	return m.id
}

func (m *Module) GetState() State {
	return m.state
}

func (m *Module) SetState(state State) {
	m.state = state
}

func (m *Module) TransitionToRunning() {
	if m.state == InitState {
		fmt.Printf("Module %s is transitioning to idle state.\n", m.id)
		m.state = RunningState
		go m.startBackgroundProcess() // Start the background process in a goroutine
	} else {
		fmt.Printf("Module %s is already in idle state.\n", m.id)
	}
}

func (m *Module) ResolveError() {
	m.resolveErrorAndResume()
}

func (m *Module) startBackgroundProcess() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-m.stopChan:
			fmt.Printf("Module %s background process stopped.\n", m.id)
			return
		case <-ticker.C:
			// Check if the module is in ErrorState
			if m.state == ErrorState {
				fmt.Printf("Module %s is in error state, pausing background process.\n", m.id)
				// Pause the process by waiting for a signal to resume
				<-m.stopChan // This will block until the module is signaled to resume
				continue     // Skip the rest of the loop iteration
			}
			number, err := fetchRandomNumber()
			if err != nil {
				m.SetState(ErrorState) // Transition to ErrorState on error
				fmt.Printf("Module %s encountered an error: %v\n", m.id, err)
				return
			}
			m.PublishToTopic("randomInt", number)
		}
	}
}

func (m *Module) resolveErrorAndResume() {
	// Hypothetical error resolution logic here
	fmt.Printf("Module %s error resolved, resuming background process.\n", m.id)
	m.SetState(RunningState) // Transition back to RunningState
	m.stopChan <- 1          // Send a signal to resume the background process
}

func (m *Module) StopBackgroundProcess() {
	if m.state == RunningState {
		m.state = ShutdownState // Transition to Shutdown state
		m.stopChan <- 0         // Signal the background process to stop
	}
}

type NotificationCallback = func(valueName string, value interface{})

// NotifySubscriber notifies the module about a value change
func (m *Module) NotifySubscriber(valueName string, value interface{}) {
	if m.notifier != nil {
		m.notifier(valueName, value)
	} else {
		fmt.Printf("Module %s received %s value update: %v\n", m.id, valueName, value)
	}
}

func (m *Module) SetNotificationCallback(callback NotificationCallback) {
	m.notifier = callback
}

func (m *Module) SubscribeToTopic(topic string, target string) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&SubscribeCommand{subscriberID: m.id, publisherID: target, topic: topic}, target)
	}
}

func (m *Module) UnsubscribeFromTopic(topic string, target string) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&UnsubscribeCommand{subscriberID: m.id, publisherID: target, topic: topic}, target)
	}
}

func (m *Module) PublishToTopic(topic string, value interface{}) {
	if m.state != ErrorState {
		m.Mediator.SendCommand(&PublishValueCommand{publisherID: m.id, topic: topic, value: value}, m.id)
	}
}

type SpecialModule struct {
	*Module
	specialValue interface{}
}

func NewSpecialModule(id string, controller IMediator, specialValue interface{}) *SpecialModule {
	return &SpecialModule{
		Module:       NewModule(id, controller),
		specialValue: specialValue,
	}
}

// IModuleFactory interface
type IModuleFactory interface {
	CreateModule(id string, controller IMediator) *Module
	CreateSpecialModule(id string, controller IMediator, specialValue interface{}) *SpecialModule
}

// DefaultModuleFactory implementation
type DefaultModuleFactory struct{}

func (f *DefaultModuleFactory) CreateModule(id string, controller IMediator) *Module {
	return NewModule(id, controller)
}

func (f *DefaultModuleFactory) CreateSpecialModule(id string, controller IMediator, specialValue interface{}) *SpecialModule {
	return NewSpecialModule(id, controller, specialValue)
}
