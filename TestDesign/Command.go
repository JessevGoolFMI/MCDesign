package TestDesign

/*
This file defines the ICommand interface and its concrete implementations (SubscribeCommand, UnsubscribeCommand, and PublishValueCommand) within the context of a mediator pattern implementation. These command structs encapsulate specific actions that modules can execute, such as subscribing to, unsubscribing from, or publishing values.

Key Features:

    ICommand Interface: The ICommand interface defines a single method, Execute, which takes a pointer to a BaseModule as its argument. This method is intended to be implemented by concrete command structs to perform specific actions.

    SubscribeCommand Struct: The SubscribeCommand struct represents a command to subscribe a module to a value. It implements the Execute method by calling the Subscribe method on the module's mediator, passing the subscriber ID, publisher ID, and value name.

    UnsubscribeCommand Struct: The UnsubscribeCommand struct represents a command to unsubscribe a module from a value. It implements the Execute method by calling the Unsubscribe method on the module's mediator, passing the subscriber ID, publisher ID, and value name.

    PublishValueCommand Struct: The PublishValueCommand struct represents a command to publish a value from a module. It implements the Execute method by first checking if the module's ID matches the publisher ID. If the IDs match, it calls the NotifySubscribers method on the module's mediator, passing the publisher ID, value name, and value. This ensures that only the intended publisher can publish values.

    GetTargetID Method: The GetTargetID method is implemented in the SubscribeCommand struct to return the publisher ID. This method could be used in scenarios where the target ID of a command is needed, such as when routing commands to the correct module.

This file showcases the command pattern in Go, focusing on how commands encapsulate actions that modules can execute. By implementing the ICommand interface, each command struct can be executed by a module, promoting a clean separation of concerns and making the system more modular and easier to extend. The use of a mediator within the commands allows for decoupled communication between modules, adhering to the principles of the mediator pattern.
*/

// ICommand interface
type ICommand interface {
	Execute(module *BaseModule) error
}

// SubscribeCommand struct
type SubscribeCommand struct {
	subscriberID string
	publisherID  string
	topic        string
}

func (sc *SubscribeCommand) Execute(module *BaseModule) error {
	module.Mediator.Subscribe(sc.subscriberID, sc.publisherID, sc.topic)
	return nil
}

func (sc *SubscribeCommand) GetTargetID() string {
	return sc.publisherID
}

type UnsubscribeCommand struct {
	subscriberID string
	publisherID  string
	topic        string
}

func (uc *UnsubscribeCommand) Execute(module *BaseModule) error {
	module.Mediator.Unsubscribe(uc.subscriberID, uc.publisherID, uc.topic)
	return nil
}

// PublishValueCommand struct
type PublishValueCommand struct {
	publisherID string
	topic       string
	value       interface{}
}

func (svc *PublishValueCommand) Execute(module *BaseModule) error {
	if module.id == svc.publisherID {
		module.Mediator.NotifySubscribers(svc.publisherID, svc.topic, svc.value)
	}
	return nil
}
