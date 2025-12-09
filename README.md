# Go RabbitMQ Examples

Professional examples demonstrating various RabbitMQ patterns in Go. This project includes a robust event library and multiple examples showcasing different messaging patterns.

## 🚀 Features

- **Robust Connection Management**: Auto-reconnect logic with graceful shutdown
- **Publisher Confirms**: Reliable message publishing with acknowledgments
- **Manual Acknowledgments**: Consumer-side message handling with Ack/Nack
- **Multiple Patterns**: Pub/Sub, Routing, Topics, and RPC examples
- **Professional Code**: Clean architecture, error handling, and logging

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for RabbitMQ)
- RabbitMQ server (via Docker or local installation)

## 🎯 Examples

### 1. Publish/Subscribe (Fanout Exchange)

**Pattern**: All consumers receive all messages (broadcasting).

**Start receivers** (in separate terminals):
```bash
go run examples/pubsub/receive_logs.go
go run examples/pubsub/receive_logs.go
```

**Send messages**:
```bash
go run examples/pubsub/emit_log.go "Hello, World!"
go run examples/pubsub/emit_log.go "This is a broadcast message"
```

### 2. Routing (Direct Exchange)

**Pattern**: Messages are routed based on exact routing key matches.

**Start receivers**:
```bash
# Terminal 1: Listen for error messages only
go run examples/routing/receive_logs_direct.go error

# Terminal 2: Listen for info and warning messages
go run examples/routing/receive_logs_direct.go info warning
```

**Send messages**:
```bash
go run examples/routing/emit_log_direct.go error "Error occurred!"
go run examples/routing/emit_log_direct.go info "Just an info message"
go run examples/routing/emit_log_direct.go warning "Warning: disk space low"
```

### 3. Topics (Topic Exchange)

**Pattern**: Messages are routed based on wildcard patterns.

**Routing key format**: `<source>.<severity>`

**Start receivers**:
```bash
# Terminal 1: Receive all kernel messages
go run examples/topics/receive_logs_topic.go "kern.*"

# Terminal 2: Receive all critical messages from any source
go run examples/topics/receive_logs_topic.go "*.critical"

# Terminal 3: Receive all messages
go run examples/topics/receive_logs_topic.go "#"
```

**Send messages**:
```bash
go run examples/topics/emit_log_topic.go "kern.critical" "A critical kernel error!"
go run examples/topics/emit_log_topic.go "kern.info" "Kernel info message"
go run examples/topics/emit_log_topic.go "user.critical" "Critical user error"
```

### 4. RPC (Request/Reply)

**Pattern**: Client sends a request and waits for a response.

**Start the RPC server**:
```bash
go run examples/rpc/rpc_server.go
```

**Make RPC calls** (in separate terminal):
```bash
go run examples/rpc/rpc_client.go 10
go run examples/rpc/rpc_client.go 30
```

The server computes the Fibonacci number and returns the result.

## 🔧 Event Library API

### Connection

```go
conn := event.NewConnection("amqp://guest:guest@localhost:5672")
defer conn.Close()
```

### Producer

```go
producer := event.NewProducer(conn)
err := producer.Publish("exchange_name", "routing_key", []byte("message"))
```

### Consumer

```go
consumer := event.NewConsumer(conn)
err := consumer.Listen("exchange_name", "exchange_type", "queue_name", "routing_key", 
    func(body []byte) error {
        log.Printf("Received: %s", body)
        return nil
    })
```

### Exchange Types

- `fanout`: Broadcasts to all bound queues
- `direct`: Routes based on exact routing key
- `topic`: Routes based on pattern matching
- `headers`: Routes based on message headers

## 🎨 Features

### Connection Manager
- Automatic reconnection on connection failures
- Graceful shutdown support
- Connection health monitoring

### Producer
- Publisher Confirms for reliable delivery
- Timeout handling
- Persistent message delivery

### Consumer
- Manual acknowledgments (Ack/Nack)
- Error handling with requeue support
- QoS (prefetch) support

## 📖 RabbitMQ Concepts

### Exchanges
- **Fanout**: Routes messages to all bound queues (Pub/Sub)
- **Direct**: Routes based on exact routing key match (Routing)
- **Topic**: Routes based on wildcard patterns (Topics)
- **Headers**: Routes based on message headers

### Queues
- Store messages until consumed
- Can be durable or transient
- Can be exclusive to a connection

### Bindings
- Link exchanges to queues with routing keys
- Define message routing rules

## 🧪 Testing

To test the examples, ensure RabbitMQ is running and execute the commands in the Examples section.

Monitor the RabbitMQ Management UI to see:
- Connections and channels
- Exchanges and their types
- Queues and their messages
- Message rates

## 🐛 Troubleshooting

**Docker fails to start**:
- Ensure Docker Desktop is running
- Check if port 5672 and 15672 are available

**Connection errors**:
- Verify RabbitMQ is running: `docker ps`
- Check RabbitMQ logs: `docker logs rabbitmq`

**Messages not received**:
- Ensure consumers are started before sending messages
- Check exchange and queue bindings in Management UI
