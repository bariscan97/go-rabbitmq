# go-rabbitmq

A **minimal, idiomatic Go wrapper** around the official [`amqp091-go`](https://github.com/rabbitmq/amqp091-go) client that makes it trivial to publish and consume events through RabbitMQ topic exchanges.\
The goal is to offer just enough abstraction to cover the 80 % use-case—**no reconnection loops, pooling, or opinionated frameworks—just plain Go interfaces** you can embed anywhere.

---

## ✨ Features

| Feature                                    | Notes                                                                                                 |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------- |
| **Tiny API surface**                       | Two public helpers—`Emitter` and `Consumer`—expose exactly the methods you need: `Push` and `Listen`. |
| **Topic exchange out-of-the-box**          | Exchange name is hard-coded to `logs_topic`, declared automatically on first use.                     |
| **Random, exclusive queues for consumers** | Every `Consumer` gets a fresh, auto-deleted queue, perfect for fan-out pub/sub.                       |
| **100 % Go**                               | Only dependency is `github.com/rabbitmq/amqp091-go v1.10.0`.                                          |
| **Example CLIs & unit tests included**     | Ready-to-run `emitter/` and `consumer/` binaries plus table-driven tests.                             |

---

## 🚀 Quick Start

### 1. Prerequisites

- **Go ≥ 1.22** (see `go.mod`)
- **RabbitMQ 3.13+** running locally or in Docker:

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### 2. Installation

```bash
go get github.com/bariscan97/go-rabbitmq
```

(You’ll import `github.com/bariscan97/go-rabbitmq/event` in your code.)

---

## 🛠️ Usage

### Publish events

```go
package main

import (
    "fmt"
    "log"

    "github.com/bariscan97/go-rabbitmq/event"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
    defer conn.Close()

    emitter, _ := event.NewEventEmitter(conn)
    for i := 1; i <= 10; i++ {
        if err := emitter.Push(
            fmt.Sprintf("[%d] – hello world", i),
            "INFO",
        ); err != nil {
            log.Fatal(err)
        }
    }
}
```

### Consume events

```go
package main

import (
    "github.com/bariscan97/go-rabbitmq/event"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672")
    defer conn.Close()

    consumer, _ := event.NewConsumer(conn)
    // Listen to multiple severities
    consumer.Listen([]string{"INFO", "ERROR"})
}
```

### CLI demos

```bash
# one terminal
go run consumer/main.go "INFO"

# second terminal
go run emitter/main.go  "INFO"
```

Watch the consumer terminal receive and print each message.

---

## 🗂 Project layout

```
go-rabbitmq/
├── event/          # Library code (Emitter, Consumer, helpers)
│   ├── emitter.go
│   ├── consumer.go
│   ├── event.go
│   └── event_test.go
├── emitter/        # Example publisher CLI
│   └── sender.go
├── consumer/       # Example subscriber CLI
│   └── consumer.go
└── go.mod
```
