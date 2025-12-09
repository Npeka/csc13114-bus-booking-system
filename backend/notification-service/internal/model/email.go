package model

type Email struct {
	From    EmailAddress `json:"from"`
	To      []string     `json:"to"`
	Subject string       `json:"subject"`
	Body    string       `json:"body"`
}
type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
