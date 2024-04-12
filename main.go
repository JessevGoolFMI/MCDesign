package main

import (
	"fmt"
	"math/rand"
	"mcs/TestDesign"
	"mcs/TestDesign/Strategies/CompressorStrategies"
	"sync"
	"time"
)

/*
This Go program demonstrates the use of a mediator pattern to manage interactions between modules in a Go application. It showcases the creation, registration, and dynamic management of modules through a central mediator, the MasterController. The program also introduces a specialized module, CompressorModule, which can subscribe to values and set a custom notification callback.

Key Components and Flow:

    Initialization: The program initializes a MasterController and a DefaultModuleFactory. It then creates and registers several modules (module1, module2, module3) and a specialized module (specialModule) with the controller. This setup establishes a mediated environment where modules can communicate indirectly through the mediator.

    Module Creation and Registration: The program creates three modules (module1, module2, module3) and a specialized module (specialModule) using the DefaultModuleFactory. Each module is registered with the MasterController.

    Subscription and Unsubscription: The example demonstrates how modules can subscribe to and unsubscribe from values. Module1 subscribes to Module2's "x" value updates and later unsubscribes. CompressorModule subscribes to Module2's "x" value updates and sets a custom notification callback to handle value updates.

    Value Publishing: A goroutine simulates Module2's "x" value changing every 500 milliseconds. When Module2 publishes a new value, the mediator notifies all subscribers, including Module1 and CompressorModule.

    Specialized Module: The CompressorModule is a specialized version of the Module that includes an additional field (specialValue). It demonstrates how modules can be extended to provide additional functionality, such as setting a custom notification callback.

    Dynamic Subscription Management: The example includes dynamic subscription management, where Module1 unsubscribes from Module2's "x" value updates and then resubscribes after a delay. This showcases the flexibility of the mediator pattern in managing subscriptions.

    Concurrency and Synchronization: The use of goroutines and the time.Sleep function to simulate asynchronous behavior and delays highlights the concurrency model of Go. It also demonstrates how the mediator pattern can manage concurrent operations, such as value updates and notifications.

Conclusion:

This code serves as a practical example of how the mediator and command patterns can be used to manage complex interactions between modules in a Go application. It demonstrates the decoupling of modules, the encapsulation of actions as commands, and the use of a mediator to manage subscriptions and notifications. The example also introduces a specialized module to show how the system can be extended to support additional functionalities.
*/

func main() {
	subscriptions()
	//patterns.Main()
}

func subscriptions() {
	controller := TestDesign.NewMasterController()
	factory := &TestDesign.DefaultModuleFactory{}

	module1 := factory.CreateModule("module1", controller)
	err := controller.RegisterModule(module1)
	if err != nil {
		return
	}
	module1.TransitionToRunning()
	module2 := factory.CreateModule("module2", controller)
	err = controller.RegisterModule(module2)
	if err != nil {
		return
	}
	module3 := factory.CreateModule("module3", controller)
	err = controller.RegisterModule(module3)
	if err != nil {
		return
	}
	compressorModule := factory.CreateCompressorModule("compressorModule", controller, "Special value")
	err = controller.RegisterModule(compressorModule)
	compressorModule.TransitionToRunning()
	if err != nil {
		return
	}
	compressorModule.SubscribeToTopic("x", "module2")
	compressorModule.SetNotificationCallback(func(valueName string, value any) {
		fmt.Printf("This is a message from the callback in compressorModule %v %v\n", valueName, value)
	})

	// Module1 requests to subscribe to Module2's "x" value updates
	module1.SubscribeToTopic("x", "module2")
	module1.UnsubscribeFromTopic("x", "module3")
	module2.SubscribeToTopic("randomInt", "compressorModule")
	module2.SubscribeToTopic("randomInt", "module1")

	// Simulate Module2's "x" value changing every 500 ms
	go func() {
		x := 0
		for {
			x++
			// Share the updated "x" value with subscribers
			module2.PublishToTopic("x", x)
			time.Sleep(500 * time.Millisecond)
		}
	}()
	time.Sleep(time.Second * 20)
	testSubscriptions(controller, module1)
	time.Sleep(time.Second * 10)
	testUnregisterModule(controller, module2)
	time.Sleep(time.Second * 10)
	testErrorState(module1)
	time.Sleep(time.Second * 10)
	testErrorState(module2)
	time.Sleep(time.Second * 10)

	unregisterRandomModules(controller, 10)
	time.Sleep(time.Second * 10)
	unregisterRandomModulesAsync(controller, 1000)
	//Keep main running for a few more seconds before shutting down
	time.Sleep(time.Second * 10)

	testCompressorModule(compressorModule)

}

func testCompressorModule(module *TestDesign.CompressorModule) {
	strategyFactory := CompressorStrategies.CompressorFactory{}
	fmt.Println("------------------------------------")
	if strategy, err := strategyFactory.CreateStrategy("v1"); err == nil {
		module.SetCompressorStrategy(strategy)
		compressed, err := module.Compress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(compressed)
	} else {
		fmt.Println(err)
	}

	if strategy, err := strategyFactory.CreateStrategy("v2"); err == nil {
		module.SetCompressorStrategy(strategy)
		compressed, err := module.Compress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(compressed)
	} else {
		fmt.Println(err)
	}

	if strategy, err := strategyFactory.CreateStrategy("v3"); err == nil {
		module.SetCompressorStrategy(strategy)
		compressed, err := module.Compress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(compressed)
	} else {
		fmt.Println(err)
	}

	if strategy, err := strategyFactory.CreateStrategy("v4"); err == nil {
		module.SetCompressorStrategy(strategy)
		compressed, err := module.Compress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(compressed)
	} else {
		fmt.Println(err)
	}

	if strategy, err := strategyFactory.CreateStrategy("v5"); err == nil {
		module.SetCompressorStrategy(strategy)
		compressed, err := module.Compress()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(compressed)
	} else {
		fmt.Println(err)
	}
	fmt.Println("------------------------------------")
}

func testSubscriptions(controller *TestDesign.MasterController, module *TestDesign.Module) {
	fmt.Println("------------------------------------")
	fmt.Printf("Unsubscribing %v from x on module 2\n", module.GetId())
	module.UnsubscribeFromTopic("x", "module2")
	controller.DisplaySubscriptions()
	fmt.Printf("Subscribing %v to x on module 2\n", module.GetId())
	module.SubscribeToTopic("x", "module2")
	controller.DisplaySubscriptions()
}

func testUnregisterModule(controller *TestDesign.MasterController, module *TestDesign.Module) {
	fmt.Println("------------------------------------")
	moduleId := module.GetId()
	fmt.Printf("Unregistering %s\n", moduleId)
	err := controller.UnregisterModule(moduleId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Unregistered  %s\n", moduleId)
	fmt.Println("------------------------------------")
	time.Sleep(time.Second * 10)

	fmt.Printf("registering  %s\n", moduleId)
	err = controller.RegisterModule(module)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("registered  %s\n", moduleId)
	fmt.Println("------------------------------------")
}

func testErrorState(module *TestDesign.Module) {
	fmt.Println("------------------------------------")
	fmt.Printf("Putting %v in error state\n", module.GetId())
	module.SetState(TestDesign.ErrorState)
	fmt.Printf("%v should be in error state\n", module.GetId())
	fmt.Println("------------------------------------")
	time.Sleep(time.Second * 10)
	fmt.Println("------------------------------------")
	fmt.Printf("Putting %v in running state\n", module.GetId())
	module.SetState(TestDesign.RunningState)
	fmt.Printf("%v should be in running state\n", module.GetId())
	fmt.Println("------------------------------------")
}

func unregisterRandomModule(controller *TestDesign.MasterController, module *TestDesign.Module) {
	fmt.Println("Unregistering module:", module.GetId())
	err := controller.UnregisterModule(module.GetId())
	if err != nil {
		fmt.Println("Error unregistering module:", err)
		return
	}
	fmt.Printf("Unregistered %v\n", module.GetId())
	time.Sleep(time.Second * 10)
	fmt.Println("Registering module:", module.GetId())
	err = controller.RegisterModule(module)
	if err != nil {
		fmt.Println("Error registering module:", err)
		return
	}
	fmt.Printf("Registered %v\n", module.GetId())

}

// Helper function to select random modules to unregister
func selectRandomModules(controller *TestDesign.MasterController, amount int) []*TestDesign.Module {
	controllerModules := controller.GetModules()
	keys := make([]string, 0, len(controllerModules))
	for key := range controllerModules {
		keys = append(keys, key)
	}

	modulesToUnregister := make([]*TestDesign.Module, 0, amount)
	for i := 0; i < amount; i++ {
		if len(keys) == 0 {
			// If there are no more modules to select, break the loop.
			break
		}
		randomIndex := rand.Intn(len(keys))
		randomKey := keys[randomIndex]
		modulesToUnregister = append(modulesToUnregister, controllerModules[randomKey])
	}
	return modulesToUnregister
}

func unregisterRandomModules(controller *TestDesign.MasterController, amount int) {
	fmt.Println("------------------------------------")
	fmt.Printf("Unregistering %d modules synchronously\n", amount)
	modulesToUnregister := selectRandomModules(controller, amount)
	for _, module := range modulesToUnregister {
		unregisterRandomModule(controller, module)
	}
}

func unregisterRandomModulesAsync(controller *TestDesign.MasterController, amount int) {
	fmt.Println("------------------------------------")
	fmt.Printf("Unregistering %d modules asynchronously\n", amount)
	modulesToUnregister := selectRandomModules(controller, amount)
	var wg sync.WaitGroup

	for _, module := range modulesToUnregister {
		wg.Add(1)
		go func(m *TestDesign.Module) {
			defer wg.Done()
			unregisterRandomModule(controller, m)
		}(module)
	}

	wg.Wait()
}
