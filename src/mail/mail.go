package mail

// send message to email

type Message struct {
	To      string
	Subject string
	Body    string
}

// SendMessage sends a message to the specified email address
func SendMessage(msg Message) {
	// send email
}

// make new message
func NewMessage(to, subject, body string) Message {
	return Message{
		To:      to,
		Subject: subject,
		Body:    body,
	}
}

// send message to email
func SendEmail(to, subject, body string) {
	msg := NewMessage(to, subject, body)
	SendMessage(msg)
}
