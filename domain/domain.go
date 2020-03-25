/*Package domain contains entity models and business rules

Author: Karpov Artem, mailto: art.frela@gmail.com
Date: 2020-03-25
*/
package domain

// User is chat user entity
type User struct {
	ID       string `json:"id"`
	Nick     string `json:"nick"`
	EMail    string `json:"email"`
	Password string `json:"-"`
}

type Users []User

// Message is chat text message entity,
// which has id, author - who wrote message and body
type Message struct {
	ID     string
	Author string
	Body   string
}

type Messages []Message

// there are will be repositories of the messages and users
