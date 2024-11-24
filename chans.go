package chans

import (
	"fmt"
)

// Deref safely dereferences a pointer and returns the value or an error
// if the pointer is nil or dereferencing causes a panic.
func Deref[T any](ptr *T) (T, error) {
	var x T
	var err error
	if ptr == nil {
		return x, fmt.Errorf("[dereference on nil pointer] %v %v", ptr, x)
	}
	(func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[dereference panic] %v %v %v", r, ptr, x)
			}
		}()
		x = *ptr
	})()
	return x, err
}

// Send safely sends a value to a channel, returning an error for nil or
// closed channels, or if sending causes a panic.
func Send[T any, C ~chan T](ch C, value T) error {
	if ch == nil {
		return fmt.Errorf("[send on nil channel] channel(%v) value(%v)", ch, value)
	}
	var err error
	(func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[channel send panic] %v channel(%v) value(%v)", r, ch, value)
			}
		}()
		c := (chan T)(ch)
		if c == nil {
			err = fmt.Errorf("[send on nil channel] channel(%v) value(%v)", ch, value)
			return
		}
		c <- value
	})()
	return err
}

// Receive safely receives a value from a channel, returning an error
// for nil or closed channels, or if receiving causes a panic.
func Receive[T any, C ~chan T](ch C) (T, error) {
	var result T
	ok := true
	if ch == nil {
		return result, fmt.Errorf("[receive on nil channel] channel(%v) result(%v)", ch, result)
	}
	var err error
	(func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[channel receive panic] %v channel(%v) result(%v)", r, ch, result)
			}
		}()
		c := (chan T)(ch)
		if c == nil {
			err = fmt.Errorf("[receive on nil channel] channel(%v) result(%v)", ch, result)
			return
		}
		result, ok = <-c
		if !ok {
			err = fmt.Errorf("[receive on closed channel] channel(%v) result(%v)", ch, result)
		}
	})()
	return result, err
}

// Close safely closes a channel, returning an error for nil or
// already-closed channels, or if closing causes a panic.
func Close[T any, C ~chan T](ch C) error {
	if ch == nil {
		return fmt.Errorf("[close on nil channel] channel(%v)", ch)
	}
	var err error
	(func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("[channel close panic] %v channel(%v)", r, ch)
			}
		}()
		c := (chan T)(ch)
		if c == nil {
			err = fmt.Errorf("[close on nil channel] channel(%v)", ch)
			return
		}
		close(c)
	})()
	return err
}

func main() {
	ch := make(chan int)
	close(ch)

	if err := Send(ch, 42); err != nil {
		fmt.Printf("Send failed [%v]\n", err)
	} else {
		fmt.Println("Value sent successfully")
	}
	value, err := Receive(ch)
	if err != nil {
		fmt.Printf("Receive failed [%v]\n", err)
	} else {
		fmt.Println("Value received successfully ", value)
	}

	var cc chan int
	if err := Close(cc); err != nil {
		fmt.Printf("Close failed [%v]\n", err)
	} else {
		fmt.Println("Channel closed successfully")
	}

}
