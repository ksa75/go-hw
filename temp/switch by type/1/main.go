package main

import "fmt"

type MsgUserBalanceChanged struct {
	userID  string
	balance string
}

type MsgEventChanged struct {
	eventID string
}

func processMessage(msg interface{}) {
	switch t := msg.(type) {
	case MsgUserBalanceChanged:
		fmt.Printf("user %q balance was changed to %q\n", t.userID, t.balance)
	case MsgEventChanged:
		fmt.Printf("event %q was changed\n", t.eventID)
	default:
		fmt.Printf("unknown message: %v\n", msg)
	}
}

/*
user "user-1" balance was changed to "1000"
event "event-1" was changed
unknown message: unknown
*/
func main() {
	processMessage(MsgUserBalanceChanged{"user-1", "1000"})
	processMessage(MsgEventChanged{"event-1"})
	processMessage("unknown")
}
