package chat

import (
	"log"

	"github.com/escalopa/grpc-ws/internal/api"
	desc "github.com/escalopa/grpc-ws/pkg"
)

func (i *Implementation) JoinChat(stream desc.ChatService_JoinChatServer) error {
	username := api.GetUsername()

	messageCh := make(chan *desc.ChatMessage)
	defer func() { i.eventCh <- event{t: eventTypeRemoveClient, username: username} }()
	i.eventCh <- event{t: eventTypeAddClient, username: username, channel: messageCh}

	var (
		errCh  = make(chan error)
		recvCh = make(chan *desc.ChatMessage)
	)
	defer close(errCh)
	defer close(recvCh)

	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("receive message from %s: %v", username, err)
				errCh <- err
				return
			}
			recvCh <- msg
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case msg := <-messageCh:
			if err := stream.Send(msg); err != nil {
				log.Printf("send message to %s: %v", username, err)
				return err
			}
		case msg := <-recvCh:
			msg.User = username // force the username provided by the token
			i.messageCh <- msg
		case err := <-errCh:
			return err
		}
	}
}
