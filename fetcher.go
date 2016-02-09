package email

import (
	"bytes"
	"fmt"
	"net/mail"

	"github.com/bytbox/go-pop3"
)

type EmailRetriever interface {
	Retr(int) (string, error)
}

type EmailFetcher interface {
	FetchEmails() ([]mail.Message, error)
}

type tlsEmailFetcher struct {
	username string
	password string
	popUrl   string
	popPort  int
}

func NewGmailFetcher(username, password string) *tlsEmailFetcher {
	return NewTlsEmailFetcher(username, password, "pop.gmail.com", 995)
}

func NewTlsEmailFetcher(username, password, url string, port int) *tlsEmailFetcher {
	return &tlsEmailFetcher{
		username: username,
		password: password,
		popUrl:   url,
		popPort:  port,
	}
}

func (f *tlsEmailFetcher) FetchEmails() ([]*mail.Message, error) {
	uri := fmt.Sprintf("%s:%d", f.popUrl, f.popPort)
	client, err := pop3.DialTLS(uri)
	if err != nil {
		return nil, fmt.Errorf("could not dial server: %v", err)
	}
	defer client.Quit()

	err = client.Auth(f.username, f.password)
	if err != nil {
		return nil, fmt.Errorf("could not authenticate: %v", err)
	}

	msgIds, _, err := client.ListAll()
	if err != nil {
		return nil, fmt.Errorf("could not list messages: %v", err)
	}

	return f.harvestMessages(client, msgIds)
}

func (f *tlsEmailFetcher) harvestSingle(retriever EmailRetriever, msgId int) (*mail.Message, error) {
	var message *mail.Message
	text, err := retriever.Retr(msgId)
	if err != nil {
		return message, fmt.Errorf("could not retrieve message (id=%d): %v", msgId, err)
	}
	msg, err := mail.ReadMessage(bytes.NewBufferString(text))
	if err != nil {
		return message, fmt.Errorf("could not read message (id=%d): %v", msgId, err)
	}
	message = msg
	return message, nil
}

func (f *tlsEmailFetcher) FetchEmailsChan(msg chan mail.Message) error {
	uri := fmt.Sprintf("%s:%d", f.popUrl, f.popPort)
	client, err := pop3.DialTLS(uri)
	if err != nil {
		return fmt.Errorf("could not dial server: %v", err)
	}
	defer client.Quit()

	err = client.Auth(f.username, f.password)
	if err != nil {
		return fmt.Errorf("could not authenticate: %v", err)
	}

	msgIds, _, err := client.ListAll()
	if err != nil {
		return fmt.Errorf("could not list messages: %v", err)
	}
	fmt.Printf("Found %d messages\n", len(msgIds))
	for _, msgId := range msgIds {

		m_, err := f.harvestSingle(client, msgId)
		fmt.Printf("Writing message\n")
		if err != nil {
			return err
		}
		msg <- *m_
	}
	return nil
}

func (f *tlsEmailFetcher) harvestMessages(retriever EmailRetriever, msgIds []int) ([]*mail.Message, error) {
	messages := make([]*mail.Message, len(msgIds))
	for i, id := range msgIds {
		text, err := retriever.Retr(id)
		if err != nil {
			return messages, fmt.Errorf("could not retrieve message (id=%d): %v", id, err)
		}
		msg, err := mail.ReadMessage(bytes.NewBufferString(text))
		if err != nil {
			return messages, fmt.Errorf("could not read message (id=%d): %v", id, err)
		}
		messages[i] = msg
	}
	return messages, nil
}
