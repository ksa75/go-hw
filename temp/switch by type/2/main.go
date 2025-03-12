package main

import "fmt"

type MsgUserBalanceChanged struct {
	userID  string
	balance string
}

func (c *MsgUserBalanceChanged) String() string {
	return fmt.Sprintf("user %q balance was changed to %q", c.userID, c.balance)
}

func (c *MsgUserBalanceChanged) Do() error {
	return nil
}

type MsgEventChanged struct {
	eventID string
}

func (c *MsgEventChanged) String() string {
	return fmt.Sprintf("event %q was changed", c.eventID)
}

func (c *MsgEventChanged) Do() error {
	return nil
}

func processMessage(msg interface{}) {

	if obj, ok := msg.(fmt.Stringer); ok {
		fmt.Println(obj.String())
		return
	}

	fmt.Println(msg)
}

func do(obj interface{}) error {
	// приведение к анонимному интерфейсу
	if obj, ok := obj.(interface{ Do() error }); ok {
		return obj.Do()
	}
	return fmt.Errorf("")
}

/*
user "user-1" balance was changed to "1000"
event "event-1" was changed
unknown message: unknown
*/
func main() {
	processMessage(&MsgUserBalanceChanged{"user-1", "1000"})
	processMessage(&MsgEventChanged{"event-1"})
	processMessage("unknown")

	do(&MsgUserBalanceChanged{"user-1", "1000"})
	do(&MsgEventChanged{"event-1"})
	do("unknown")
}
