// Package micro is a pluggable framework for microservices
package micro2




type Publisher struct {
	topic string
	//Publish(topic string) error
}

func (p *Publisher) Publish(data interface{}){
	Send(p.topic, data)

}


// NewService creates and returns a new Service based on the packages within.
func NewService() GRPCService {
	return NewGRPCService()
}



// NewPublisher returns a new Publisher
func NewPublisher(topic string) *Publisher {
	return &Publisher{ topic}
}


// RegisterSubscriber is syntactic sugar for registering a subscriber
func RegisterSubscriber(topic string, handler interface{}) error {
	Receive(topic, handler)
	return nil
}
