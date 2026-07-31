package outbound

type Client interface {
	ID() string
	Send(payload []byte)
	Close() error
}
