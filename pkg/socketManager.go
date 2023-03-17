package pkg

type SocketManager interface {
}

type Socket interface {
	Send()
	Close()
}
