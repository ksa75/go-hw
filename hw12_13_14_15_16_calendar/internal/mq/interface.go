package mq

type Publisher interface {
	Publish(queueName string, body []byte) error
	Close() error
}

type Consumer interface {
	Consume(queueName string, handler func([]byte)) error
	Close() error
}
