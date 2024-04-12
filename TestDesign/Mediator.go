package TestDesign

import (
	"errors"
	"fmt"
	"sync"
)

/*
This file contains the implementation of a simple mediator pattern in Go, specifically designed for managing and executing commands in a thread-safe manner. The MasterController struct serves as the central mediator, facilitating communication between different modules (Module and CompressorModule) by managing subscriptions, executing commands, and notifying subscribers about value changes.

Key Features:

    Mediator Pattern: The MasterController acts as the mediator, centralizing the communication between modules. This design pattern helps in decoupling the modules, allowing them to communicate indirectly through the mediator.

    Command Pattern: Commands are encapsulated as objects, allowing for the execution of specific actions. The MasterController processes these commands concurrently using a worker pool pattern, improving performance and scalability.

    Thread Safety with Mutex: To ensure thread safety when accessing shared resources, such as the modules map, a mutex (sync.Mutex) is employed. This mechanism prevents race conditions by allowing only one goroutine to access the shared resource at a time.

    Concurrency and Synchronization: The MasterController uses a combination of goroutines and a wait group (sync.WaitGroup) to process commands concurrently. This approach enhances the application's performance by leveraging Go's concurrency model.

    Subscription Management: The mediator supports subscribing and unsubscribing modules to values, allowing for a flexible and dynamic communication system. It maintains a map of subscriptions to manage these relationships.

    Command Queue: Commands are queued for processing, ensuring that they are executed in the order they are received. This design helps in managing the flow of commands and maintaining the integrity of the system.

This implementation demonstrates a practical application of the mediator and command patterns in Go, showcasing how to manage complex interactions between objects in a structured and efficient manner. The use of a mutex ensures that the system remains robust and thread-safe, making it suitable for use in concurrent environments.
*/

// IMediator interface
type IMediator interface {
	SendCommand(command ICommand, targetID string)
	GetModule(id string) *Module
	Subscribe(subscriberID, publisherID, valueName string)
	Unsubscribe(subscriberID, publisherID, valueName string)
	NotifySubscribers(publisherID, valueName string, value interface{})
}

// MasterController struct
type MasterController struct {
	modules       map[string]*Module
	subscriptions map[string]map[string]bool
	commandQueue  chan commandWithTargetID // Command queue channel
	wg            sync.WaitGroup
	mu            sync.Mutex
}
type commandWithTargetID struct {
	command  ICommand
	targetID string
}

func NewMasterController() *MasterController {
	mc := &MasterController{
		modules:       make(map[string]*Module),
		subscriptions: make(map[string]map[string]bool),
		commandQueue:  make(chan commandWithTargetID), // Initialize the command queue channel
	}

	// Start a goroutine to process commands from the queue
	numWorkers := 5 // Adjust based on your needs
	for i := 0; i < numWorkers; i++ {
		mc.wg.Add(1)
		go mc.processCommands()
	}

	return mc
}

func (mc *MasterController) GetModules() map[string]*Module {
	return mc.modules
}

func (mc *MasterController) processCommands() {
	defer mc.wg.Done()
	for command := range mc.commandQueue {
		mc.mu.Lock()
		if targetModule := mc.modules[command.targetID]; targetModule != nil {
			mc.mu.Unlock()
			if err := command.command.Execute(targetModule); err != nil {
				// Handle error, e.g., log it
				fmt.Printf("Error executing command: %v\n", err)
			}
		} else {
			mc.mu.Unlock()
			// Optionally, handle the case where the target module is not found
		}
	}
}

func (mc *MasterController) DisplaySubscriptions() {
	fmt.Println("--------------------")
	for s, m := range mc.subscriptions {
		fmt.Printf("Subscription %v %v\n", s, m)
	}
	fmt.Println("--------------------")
}

// Subscribe allows us to subscribe to any value. These values don't have to exist at the time of subscription
func (mc *MasterController) Subscribe(subscriberID, publisherID, valueName string) {
	key := publisherID + ":" + valueName
	mc.mu.Lock() // Lock for writing to the subscriptions map
	if _, exists := mc.subscriptions[key]; !exists {
		mc.subscriptions[key] = make(map[string]bool)
	}
	mc.subscriptions[key][subscriberID] = true
	mc.mu.Unlock() // Unlock after writing
}

func (mc *MasterController) Wait() {
	mc.wg.Wait()
}

func (mc *MasterController) Unsubscribe(subscriberID, publisherID, valueName string) {
	key := publisherID + ":" + valueName
	if subscribers, exists := mc.subscriptions[key]; exists {
		// Check if the subscriber is actually subscribed
		if _, subscribed := subscribers[subscriberID]; subscribed {
			// If the subscriber is subscribed, remove them from the list
			delete(subscribers, subscriberID)
			fmt.Printf("Subscriber %s unsubscribed from %s:%s\n", subscriberID, publisherID, valueName)
		} else {
			fmt.Printf("Subscriber %s is not subscribed to %s:%s\n", subscriberID, publisherID, valueName)
		}
	} else {
		fmt.Printf("No subscribers found for %s:%s\n", publisherID, valueName)
	}
}

func (mc *MasterController) NotifySubscribers(publisherID, valueName string, value interface{}) {
	key := publisherID + ":" + valueName
	if subscribers, exists := mc.subscriptions[key]; exists {
		for subscriberID := range subscribers {
			if module := mc.GetModule(subscriberID); module != nil && module.GetState() != ErrorState {
				module.NotifySubscriber(valueName, value)
			}
		}
	}
}

func (mc *MasterController) RegisterModule(module interface{}) error {
	if m, ok := module.(*Module); ok {
		mc.modules[m.id] = m
	} else if sm, ok := module.(*CompressorModule); ok {
		mc.modules[sm.id] = sm.Module
	} else if sm, ok := module.(*DispenserModule); ok {
		mc.modules[sm.id] = sm.Module
	} else {
		return errors.New("module not supported")
	}
	return nil
}

func (mc *MasterController) UnregisterModule(moduleId string) error {
	//Can add something in here that notifies modules if their subscribed module unregisters
	if mc.modules[moduleId] != nil {
		delete(mc.modules, moduleId)
	} else {
		return errors.New("module id not found")
	}
	return nil
}

func (mc *MasterController) SendCommand(command ICommand, targetID string) {
	// Send the command with its target ID to the commandQueue channel
	mc.commandQueue <- commandWithTargetID{command: command, targetID: targetID}
}

func (mc *MasterController) GetModule(id string) *Module {
	return mc.modules[id]
}
