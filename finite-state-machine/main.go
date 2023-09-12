package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/qmuntal/stateless"
)

const (
	triggerCallDialed             = "CallDialed"
	triggerCallConnected          = "CallConnected"
	triggerLeftMessage            = "LeftMessage"
	triggerPlacedOnHold           = "PlacedOnHold"
	triggerTakenOffHold           = "TakenOffHold"
	triggerPhoneHurledAgainstWall = "PhoneHurledAgainstWall"
	triggerMuteMicrophone         = "MuteMicrophone"
	triggerUnmuteMicrophone       = "UnmuteMicrophone"
	triggerSetVolume              = "SetVolume"
)

const (
	stateOffHook        = "OffHook"
	stateRinging        = "Ringing"
	stateConnected      = "Connected"
	stateOnHold         = "OnHold"
	statePhoneDestroyed = "PhoneDestroyed"
)

func main() {
	phoneCall := stateless.NewStateMachine(stateOffHook)
	phoneCall.SetTriggerParameters(triggerSetVolume, reflect.TypeOf(0))
	phoneCall.SetTriggerParameters(triggerCallDialed, reflect.TypeOf(""))

	phoneCall.Configure(stateOffHook).
		Permit(triggerCallDialed, stateRinging)

	phoneCall.Configure(stateRinging).
		OnEntryFrom(triggerCallDialed, func(_ context.Context, args ...any) error {
			onDialed(args[0].(string))
			return nil
		}).
		Permit(triggerCallConnected, stateConnected)

	phoneCall.Configure(stateConnected).
		OnEntry(startCallTimer).
		OnExit(func(_ context.Context, _ ...any) error {
			stopCallTimer()
			return nil
		}).
		InternalTransition(triggerMuteMicrophone, func(_ context.Context, _ ...any) error {
			onMute()
			return nil
		}).
		InternalTransition(triggerUnmuteMicrophone, func(_ context.Context, _ ...any) error {
			onUnmute()
			return nil
		}).
		InternalTransition(triggerSetVolume, func(_ context.Context, args ...any) error {
			onSetVolume(args[0].(int))
			return nil
		}).
		Permit(triggerLeftMessage, stateOffHook).
		Permit(triggerPlacedOnHold, stateOnHold)

	phoneCall.Configure(stateOnHold).
		SubstateOf(stateConnected).
		Permit(triggerTakenOffHold, stateConnected).
		Permit(triggerPhoneHurledAgainstWall, statePhoneDestroyed)

	phoneGraph := phoneCall.ToGraph()
	err := stateMachineGenerator(phoneGraph)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	phoneCall.Fire(triggerCallDialed, "qmuntal")
	phoneCall.Fire(triggerCallConnected)
	phoneCall.Fire(triggerSetVolume, 2)
	phoneCall.Fire(triggerPlacedOnHold)
	phoneCall.Fire(triggerMuteMicrophone)
	phoneCall.Fire(triggerUnmuteMicrophone)
	phoneCall.Fire(triggerTakenOffHold)
	phoneCall.Fire(triggerSetVolume, 11)
	phoneCall.Fire(triggerPlacedOnHold)
	phoneCall.Fire(triggerPhoneHurledAgainstWall)
	fmt.Printf("State is %v\n", phoneCall.MustState())

}

func stateMachineGenerator(graph string) error {
	output := []byte(graph)
	path, _ := filepath.Abs("./main.dot")
	if err := os.WriteFile(path, output, 0666); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func onSetVolume(volume int) {
	fmt.Printf("Volume set to %d!\n", volume)
}

func onUnmute() {
	fmt.Println("Microphone unmuted!")
}

func onMute() {
	fmt.Println("Microphone muted!")
}

func onDialed(callee string) {
	fmt.Printf("[Phone Call] placed for : [%s]\n", callee)
}

func startCallTimer(_ context.Context, _ ...any) error {
	fmt.Println("[Timer:] Call started at 11:00am")
	return nil
}

func stopCallTimer() {
	fmt.Println("[Timer:] Call ended at 11:30am")
}

hello
