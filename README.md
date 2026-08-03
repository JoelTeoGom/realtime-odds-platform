# realtime-odds-platform

Real-time odds platform: ingests event updates from external providers and pushes them
to subscribed clients over WebSocket.

Built from a system design problem I was given in a sportsbook interview. The interesting
constraint is ordering: providers send updates for an event in order, and that order has
to survive all the way to the client, while updates for *different* events run in parallel.

**Status:** the gateway is working. The ingestion pipeline is in progress.

## Architecture

![Architecture](docs/architecture.png)

### Domain

`Event` → `Market` → `Odd`. A provider does not resend the whole event, it sends
`OddUpdate`, a fragment describing what changed. The whole system moves fragments, not
snapshots.

### Ingestion path

External providers push event updates into **N RabbitMQ queues, routed by
`hash(eventId) % N`**. This is the core decision: everything for one event lands in one
queue, so it stays ordered, while different events spread across queues and process in
parallel. One worker per queue, so a worker never has two updates for the same event in
flight.

Each worker does two things with an update:

- Writes the new state to **Redis**, which holds the current state of every event so a
  client connecting mid-match gets a snapshot instead of waiting for the next change.
- Publishes the fragment to **Redis Pub/Sub**, so every gateway instance subscribed to
  that event receives it.

### Gateway and fan-out

A client opens `/ws` and sends `subscribe(eventId)`. One client can subscribe to many
events at once and receives updates for all of them over the same connection.

The hub is **sharded by `hash(eventId)`**. Inside each shard there is a map keyed by
event id, and each entry holds the set of clients subscribed to that event. When a
fragment arrives from Pub/Sub, the hub resolves its shard, looks up the event and pushes
to exactly those clients.

Each connection runs two goroutines: `readPump`, which reads client commands, and
`writePump`, the only goroutine that ever writes to the socket.

## Design notes

**Why hash-partitioned queues.** Parallelising naively breaks ordering: two workers could
process `score: 1` and `score: 2` for the same event out of order and leave the cache
wrong. Partitioning by event id gives parallelism between events and serialisation within
one. Same idea as Kafka partitions, applied to RabbitMQ.

**Why shard the hub.** A single subscription map guarded by one mutex serialises every
broadcast. Sharding by event id means updates for unrelated events never contend for the
same lock.

**Why one write pump per connection.** Gorilla WebSocket connections do not support
concurrent writers. Funnelling every write through a single goroutine per connection
makes fan-out safe without locking the socket.

**Why a non-blocking send.** If the send blocked, one slow consumer would stall the
broadcast for every other client subscribed to that event. Dropping the slow consumer is
the deliberate tradeoff: availability of the fan-out over delivery to one client.

**Why Redis holds state as well as transporting it.** Pub/Sub has no replay. Without the
stored state, a client connecting between two updates would see nothing until the next
change arrives.

## Structure

A monorepo with **two independent services**, each with its own Go module and its own
hexagon (ports + adapters). They share no code: they are deployed, versioned and scaled
separately, and communicate over the network (gRPC / pub-sub).
