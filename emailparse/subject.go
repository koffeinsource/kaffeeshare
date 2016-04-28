package emailparse

import (
	"mime"
	"net/mail"
)

func getSubject(msg *mail.Message) (string, error) {
	var dec mime.WordDecoder
	ret, err := dec.DecodeHeader(msg.Header.Get("Subject"))
	if err != nil {
		return "", err
	}
	return ret, nil
}
