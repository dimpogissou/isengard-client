package connectors

import "github.com/hpcloud/tail"

type Publisher struct {
	listeners []chan *tail.Line
}

type Subscriber struct {
	Channel   chan *tail.Line
	Connector ConnectorInterface
}

func (p *Publisher) Subscribe(c chan *tail.Line) {
	p.listeners = append(p.listeners, c)
}

func (p *Publisher) Publish(m *tail.Line) {
	for _, c := range p.listeners {
		c <- m
	}
}

func (s *Subscriber) ListenToChannel() {
	for data := range s.Channel {
		s.Connector.Send(data)
	}
}
