package chat

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/catalystgo/catalystgo"
	desc "github.com/escalopa/grpc-ws/pkg"
	"golang.org/x/sync/errgroup"
)

type Implementation struct {
	desc.UnimplementedChatServiceServer

	clients   map[string]chan *desc.ChatMessage
	messageCh chan *desc.ChatMessage
	eventCh   chan event

	closeCh chan struct{}
	wg      sync.WaitGroup
}

type (
	eventType int

	event struct {
		t        eventType
		username string
		channel  chan *desc.ChatMessage
	}
)

const (
	eventTypeAddClient eventType = iota
	eventTypeRemoveClient
)

func NewChatService() *Implementation {
	return &Implementation{
		clients:   make(map[string]chan *desc.ChatMessage),
		messageCh: make(chan *desc.ChatMessage, 100),
		eventCh:   make(chan event),
		closeCh:   make(chan struct{}),
	}
}

func (i *Implementation) GetDescription() catalystgo.ServiceDesc {
	return desc.NewChatServiceServiceDesc(i)
}

func (i *Implementation) run() {
	i.wg.Add(1)
	defer i.wg.Done()

	for {
		select {
		case <-i.closeCh:
			return
		case msg := <-i.messageCh:
			i.broadcast(msg)
		case e := <-i.eventCh:
			var msg string

			switch e.t {
			case eventTypeAddClient:
				i.clients[e.username] = e.channel
				msg = fmt.Sprintf("%s joined the chat", e.username)
			case eventTypeRemoveClient:
				close(e.channel)
				delete(i.clients, e.username)
				msg = fmt.Sprintf("%s left the chat", e.username)
			}

			i.broadcast(&desc.ChatMessage{User: "system", Message: msg})
		}
	}
}

func (i *Implementation) broadcast(msg *desc.ChatMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errG := &errgroup.Group{}
	for _, ch := range i.clients {
		errG.Go(func() error {
			select {
			case ch <- msg:
			case <-i.closeCh:
			case <-ctx.Done():
			}
			return nil
		})
	}

	_ = errG.Wait()
}

func (i *Implementation) Close() {
	close(i.closeCh)
	i.wg.Wait()
}
