# go-sharded-ws-hub

Hub de WebSockets shardeado para difundir cambios de cuotas en tiempo real.

## Arquitectura

![Arquitectura del hub](docs/architecture.png)

**Camino de entrada (arriba).** Los clientes abren `/ws` contra el gateway. El
handler hace el upgrade y crea un `Client` con dos goroutines: `readPump`, que
lee los comandos del navegador (`subscribe(eventId)`), y `writePump`, la única
que escribe en el socket. El `readPump` llama al hub, que elige shard con
`hash(eventId) % nShards` y guarda al cliente en el mapa de ese shard.

**Camino de actualizaciones (abajo).** Los `eventUPDATES` externos se reparten
por colas con el mismo criterio, `hash(eventId) % N`, para que los cambios de un
mismo evento se procesen en orden. Los workers de fan-out publican en pub/sub, y
cada instancia del hub recibe lo que corresponde a los eventos que tiene
suscritos y se lo empuja a los clientes de ese shard.

## Estructura

```
cmd/gateway/          binario que sirve /ws
internal/
  domain/             Event, Market, Odd, OddUpdate
  hub/                Hub + Shard: registro de suscripciones y fan-out
  ports/
    inbound/          Hub — lo que el adapter ws invoca
    outbound/         Client — lo que el hub invoca
  adapters/
    inbound/ws/       handler + client (gorilla/websocket)
```

La dependencia va siempre hacia dentro: `hub` no importa nunca `adapters`.
