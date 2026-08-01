package outbound

type Connection interface {
	ID() string
	Send(payload []byte)
	Close() error
}
