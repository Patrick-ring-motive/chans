
# chans

### **Overview**
The `chans` package introduces a consistent and safer standard for handling Go channel operations. It addresses the unpredictable behavior of default channel operations (e.g., panics, indefinite blocking) by converting these into a reliable pattern of error handling.

The goal is to simplify the table of default channel behaviors into a more predictable and user-friendly standard.

---

### **Default Channel Behaviors**

By default, Go channels exhibit the following behaviors, which can be difficult to predict and handle in real-world applications:

| **Operation**           | **A Nil Channel**       | **A Closed Channel**          | **A Not-Closed Non-Nil Channel**      |
|-------------------------|-------------------------|--------------------------------|---------------------------------------|
| **Close**               | `panic`                | `panic`                       | Succeeds to close                     |
| **Send Value To**       | Blocks forever         | `panic`                       | Blocks or succeeds to send            |
| **Receive Value From**  | Blocks forever         | Never blocks, returns zero value | Blocks or succeeds to receive         |

---

### **The `chans` Package**

The `chans` package normalizes these operations into the following predictable behaviors:

| **Operation**           | **A Nil Channel**       | **A Closed Channel**          | **A Not-Closed Non-Nil Channel**      |
|-------------------------|-------------------------|--------------------------------|---------------------------------------|
| **Close**               | `return error`         | `return error`                | Succeeds to close, `return error(nil)` |
| **Send Value To**       | `return error`         | `return error`                | Blocks or succeeds to send, `return error(nil)` |
| **Receive Value From**  | `return zeroValue,error` | `return zeroValue,error`     | Blocks or succeeds to receive, `return value,error(nil)` |

---

### **Features**

1. **Consistent Error Handling**:
   - Panics are replaced with descriptive error messages.
   - Nil channels, closed channels, and valid channels are all handled gracefully.

2. **Predictable API**:
   - All operations return an error, allowing you to handle issues programmatically rather than relying on runtime panics.

3. **Support for Generics**:
   - The package supports typed and custom channel types through Go's generics (`C ~chan T`).

4. **Enhanced Debugging**:
   - Error messages include detailed context, such as the channel state, operation type, and any values involved.

---

### **Functions**

#### `Deref[T any](ptr *T) (T, error)`
Safely dereferences a pointer. Returns the value or an error if the pointer is nil or causes a panic.

```go
var ptr *int
val, err := Deref(ptr)
if err != nil {
    fmt.Printf("Deref failed: %v\n", err)
}
```

#### `Send[T any, C ~chan T](ch C, value T) error`
Safely sends a value to a channel. Returns an error for nil or closed channels, or if sending causes a panic.

```go
type MyChan chan int
var ch MyChan // nil channel
err := Send(ch, 42)
if err != nil {
    fmt.Printf("Send failed: %v\n", err)
}
```

#### `Receive[T any, C ~chan T](ch C) (T, error)`
Safely receives a value from a channel. Returns the zero value of `T` and an error for nil or closed channels, or if receiving causes a panic.

```go
ch := make(chan int, 1)
ch <- 42
close(ch)

val, err := Receive(ch)
if err != nil {
    fmt.Printf("Receive failed: %v\n", err)
} else {
    fmt.Println("Received value:", val)
}
```

#### `Close[T any, C ~chan T](ch C) error`
Safely closes a channel. Returns an error for nil or already-closed channels, or if closing causes a panic.

```go
var ch chan int // nil channel
err := Close(ch)
if err != nil {
    fmt.Printf("Close failed: %v\n", err)
}
```

---

### **Usage Example**

Here's an example that demonstrates the behavior of the `chans` package:

```go
package main

import (
    "fmt"
    "chans"
)

func main() {
    // Example 1: Send to a closed channel
    ch := make(chan int)
    close(ch)
    if err := chans.Send(ch, 42); err != nil {
        fmt.Printf("Send failed: %v\n", err)
    }

    // Example 2: Receive from a closed channel
    value, err := chans.Receive(ch)
    if err != nil {
        fmt.Printf("Receive failed: %v\n", err)
    } else {
        fmt.Println("Received value:", value)
    }

    // Example 3: Close a nil channel
    var nilCh chan int
    if err := chans.Close(nilCh); err != nil {
        fmt.Printf("Close failed: %v\n", err)
    }

    // Example 4: Dereference a nil pointer
    var ptr *int
    val, err := chans.Deref(ptr)
    if err != nil {
        fmt.Printf("Deref failed: %v\n", err)
    } else {
        fmt.Println("Dereferenced value:", val)
    }
}
```

---

### **Why Use This Package?**

1. **Simplifies Complex Behavior**:
   - Instead of dealing with panics, indefinite blocking, or inconsistent behavior, `chans` gives you clear, predictable outcomes.

2. **Improved Debugging**:
   - Error messages include all the context needed to diagnose issues effectively.

3. **Safe and Robust**:
   - The package is designed to handle edge cases gracefully, reducing the likelihood of runtime crashes.

---

### **Planned Improvements**

1. Add more utility functions for buffered channels (e.g., checking capacity, size).
2. Expand the test suite to cover all edge cases.
3. Include optional logging for error scenarios.

---

### **Contributing**

Contributions are welcome! Feel free to submit issues or pull requests to improve the `chans` package.

