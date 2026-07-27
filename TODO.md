# TODO — go-sharded-hub

## internal/ws/pump.go — `Client`
- [ ] Fields: `conn *websocket.Conn`, `send chan []byte`, `id string`
- [ ] `readPump`: loop `conn.ReadMessage()` → `hub.Broadcast(...)`; on exit `hub.Unregister(c)` and close `conn`
- [ ] `writePump`: loop `range c.send` → `conn.WriteMessage(...)`

## internal/hub/shard.go — `shard`
- [ ] Field: `clients map[string]*Client`
- [ ] `newShard() *shard` — initialize the map
- [ ] `shardIndex(key string) int` — hash the key → shard index

## internal/hub/hub.go — `Hub`
- [ ] Field: `shards []*shard`
- [ ] `New()` — initialize the shards (and internal goroutines if needed)
- [ ] `Register(client *Client)` — pick a shard and store the client
- [ ] `Unregister(client *Client)` — find the shard and remove the client
- [ ] `Broadcast(msg []byte)` — walk the shards and enqueue the message

## internal/ws/handler.go — `Handler`
- [ ] Create a `Client` wrapping `conn`
- [ ] `h.Register(client)`
- [ ] Launch `go client.readPump()` and `go client.writePump()`
- [ ] In prod: validate `Origin` in `upgrader.CheckOrigin`

## Suggested order
`Client` (pump.go) → `shard` (shard.go) → `Hub` (hub.go) → wire it all up in `handler.go`
