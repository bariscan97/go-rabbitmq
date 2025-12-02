# Publish/Subscribe Pattern

This example demonstrates the **Publish/Subscribe (Pub/Sub)** messaging pattern using RabbitMQ's **fanout exchange**.

## Pattern Overview

In the Pub/Sub pattern:
- A **producer** publishes messages to an **exchange** (not directly to a queue)
- Multiple **consumers** can receive the same message
- Each consumer has its own **exclusive, temporary queue** bound to the exchange
- The **fanout exchange** broadcasts all messages to all bound queues

## Key Concepts

### Fanout Exchange
- Type: `fanout`
- Behavior: Broadcasts messages to **all** queues bound to it
- Routing key is **ignored** in fanout exchanges

### Exclusive Queues
- Each consumer creates a **temporary queue** with a unique name
- Queue is **exclusive** to that connection
- Queue is **automatically deleted** when the consumer disconnects
- Multiple consumers can run simultaneously, each receiving all messages

## How It Works

```
Producer → [logs exchange] → Queue 1 → Consumer 1
                           → Queue 2 → Consumer 2
                           → Queue 3 → Consumer 3
```

1. Producer declares the `logs` fanout exchange
2. Producer publishes messages to the `logs` exchange
3. Each consumer:
   - Declares the `logs` exchange
   - Creates a temporary, exclusive queue (RabbitMQ generates unique name)
   - Binds its queue to the `logs` exchange
   - Receives all messages published to the exchange

## Running the Example

### Start RabbitMQ
```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

### Start Multiple Consumers
Open 2-3 separate terminals and run:
```bash
# Terminal 1
go run receive_logs.go

# Terminal 2
go run receive_logs.go

# Terminal 3
go run receive_logs.go
```

### Send Messages
In another terminal:
```bash
go run emit_log.go "First log message"
go run emit_log.go "Another important log"
```

All running consumers will receive **all** messages!

## Key Differences from Work Queue Pattern

| Aspect | Work Queue | Pub/Sub |
|--------|-----------|---------|
| Exchange | Default (`""`) | Named fanout exchange |
| Queue | Named, durable | Temporary, exclusive |
| Message Distribution | Round-robin (one consumer gets each message) | Broadcast (all consumers get all messages) |
| Use Case | Task distribution | Event broadcasting (logging, notifications) |

## Use Cases

- **Logging systems**: Multiple log processors receive all log messages
- **Live updates**: Multiple clients receive real-time notifications
- **Cache invalidation**: All cache servers receive invalidation events
- **Fan-out architectures**: Distribute same data to multiple services
