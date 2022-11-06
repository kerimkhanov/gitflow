package models

type UserMail struct {
	ID      int      `json:"id,omitempty"`
	Emails  []string `json:"emails,omitempty"`
	Title   string   `json:"title,omitempty"`
	Message string   `json:"message,omitempty"`
}
