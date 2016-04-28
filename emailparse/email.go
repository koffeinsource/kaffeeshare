package emailparse

import (
	"fmt"
	"io"
	"net/mail"

	"github.com/koffeinsource/kaffeeshare/data"
)

// Email represents everything we extracted from an email
type Email struct {
	Images     []ImageBody
	Texts      []TextBody
	Namespaces []string
	Subject    string
}

// TextBody is the text extracted from an EMail
type TextBody struct {
	Body        string
	ContentType string
}

// ImageBody is an image extracted from an EMail
type ImageBody struct {
	Body     []byte
	Encoding string
}

// Get parses the email and returns all relevant data
func Get(con *data.Context, m io.Reader) (*Email, error) {
	msg, err := mail.ReadMessage(m)
	if err != nil {
		return nil, fmt.Errorf("Error at mail.ReadMessage in DispatchEmail. Error: %v", err)
	}
	con.Log.Debugf("header: %v", msg.Header)

	// get body
	email, err := extract(con, msg.Header, msg.Body)
	if err != nil {
		return nil, fmt.Errorf("Error at extracting the body. Error: %v", err)
	}
	//con.Log.Debugf("Received mail: %v", email)

	namespaces, err := getNamespaces(msg)
	if err != nil {
		return nil, fmt.Errorf("Error at parsing the receiver fields. Error: %v", err)
	}
	con.Log.Debugf("Detected namespaces: %v", namespaces)
	email.Namespaces = namespaces

	email.Subject, err = getSubject(msg)
	if err != nil {
		con.Log.Errorf("Could not decode subject. Error: %v", err)
		email.Subject = ""
	}

	return email, nil
}
