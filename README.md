# Redis Server

This is a simple Redis server implementation written in Go.

## Features

- Supports basic Redis commands: SET, GET, PING, and ECHO.
- Handles concurrent client connections over TCP/IP.
- Uses a synchronized map and locks for data storage to prevent race conditions.
- Supports expiration time for keys.

## Getting Started

### Prerequisites

- Go 1.15 or higher

### Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/theakhandpatel/redis-server.git
   ```

2. Change to the project directory:

   ```shell
   cd redis-server
   ```



3. Run the server:

   ```shell
   ./spawn_redis_server.sh
   ```

   By default, the server listens on port 6379.

## Usage

Connect to the Redis server using a Redis client and send commands. The supported commands are:

- `SET`: Set the value of a key.
- `GET`: Get the value of a key.
- `PING`: Check if the server is running.
- `ECHO`: Return the given message.

Example usage:

```shell
$ redis-cli
127.0.0.1:6379> SET mykey hello
OK
127.0.0.1:6379> GET mykey
"hello"
127.0.0.1:6379> SET yourkey lost PX 500ms
OK
127.0.0.1:6379> GET yourkey
nil
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> ECHO hello world
"hello world"
```

## Concurrent Clients and Race Condition Prevention

The Redis server implementation handles concurrent client connections using Goroutines. Each client connection is handled in a separate Goroutine, allowing multiple clients to interact with the server simultaneously.

To prevent race conditions when accessing and modifying the data storage (`sync.Map`), the implementation uses locks. The `sync.Map` type provides built-in synchronization for concurrent access, but locks are used when necessary to ensure exclusive access to critical sections of the code.

By utilizing Goroutines and locks, the Redis server ensures concurrent execution and prevents race conditions, maintaining data integrity.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

This Redis server implementation was inspired by the Redis open-source project (https://redis.io/).

