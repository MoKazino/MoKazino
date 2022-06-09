package web

type Broker struct {
	stopCh    chan struct{}
	publishCh chan string
	subCh     chan chan string
	unsubCh   chan chan string
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan string, 1),
		subCh:     make(chan chan string, 1),
		unsubCh:   make(chan chan string, 1),
	}
}

func (b *Broker) Start() {
	subs := map[chan string]struct{}{}
	for {
		select {
		case <-b.stopCh:
			return
		case msgCh := <-b.subCh:
			subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

func (b *Broker) Subscribe() chan string {
	msgCh := make(chan string)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan string) {
	b.unsubCh <- msgCh
}

func (b *Broker) Publish(msg string) {
	b.publishCh <- msg
}
