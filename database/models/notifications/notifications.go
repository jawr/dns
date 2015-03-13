package notifications

import (
	"encoding/json"
	"github.com/jawr/dns/database/connection"
	"github.com/jawr/dns/database/models/domains"
	"github.com/jawr/dns/database/models/users"
	"time"
)

type Message struct {
	Added   time.Time      `json:"added"`
	Message string         `json:"message"`
	Domain  domains.Domain `json:"domain"`
	ReadAt  time.Time      `json:"read_at,omitempty"`
}

type Notification struct {
	User     users.User `json:"user"`
	Messages []Message  `json:"messages"`
	Updated  time.Time  `json:"updated"`
	Alerts   int        `json:"alerts"`
	Archived []Message  `json:"archived"`
}

func SetupUserNotification(user users.User) (Notification, error) {
	note := Notification{
		User: user,
	}

	conn, err := connection.Get()
	if err != nil {
		return note, err
	}
	_, err = conn.Exec("INSERT INTO notification (user_id) VALUES ($1)", user.ID)
	if err != nil {
		return note, err
	}
	return note, nil
}

func (n *Notification) Save() error {
	conn, err := connection.Get()
	if err != nil {
		return err
	}
	messages, err := json.Marshal(n.Messages)
	if err != nil {
		return err
	}
	archived, err := json.Marshal(n.Archived)
	if err != nil {
		return err
	}
	_, err = conn.Exec(
		`UPDATE notification SET 
			messages = $1,
			alerts = $2,
			archived = $3
		 WHERE user_id = $4`,
		messages,
		n.Alerts,
		archived,
		n.User.ID,
	)
	return err
}

func (n *Notification) AddMessage(message Message) {
	np, err := GetByUser(n.User).One()
	if err != nil {
		return
	}
	n = &np
	if len(n.Messages) == 0 {
		n.Messages = append(n.Messages, message)
	} else {
		n.Messages = append([]Message{message}, n.Messages...)
	}
	n.Alerts++
	n.Save()
}

func (n *Notification) ArchiveMessages(messages []Message) {
	for i, _ := range messages {
		m := messages[i]
		for i, m2 := range n.Messages {
			if m.Added == m2.Added &&
				m.Message == m2.Message &&
				m.Domain.UUID.String() == m2.Domain.UUID.String() {
				m.ReadAt = time.Now()
				n.Archived = append(n.Archived, m)
				n.Messages = append(n.Messages[:i], n.Messages[i+1:]...)
				n.Alerts--
				break

			}
		}
	}
	n.Save()
}
