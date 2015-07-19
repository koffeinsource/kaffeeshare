package email

import (
	"net/mail"
	"strings"
)

func getNamespaces(msg *mail.Message) ([]string, error) {
	// use a 'set' to remove duplicates
	set := make(map[string]bool)

	fields := [...]string{"To", "CC"}
	for _, field := range fields {
		if msg.Header.Get(field) == "" {
			// the field is not present in the email
			continue
		}
		addresses, err := msg.Header.AddressList(field)
		if err != nil {
			return nil, err
		}

		for _, addr := range addresses {
			// TODO we should actually check for the receiver domain
			temp := addr.Address[:strings.Index(addr.Address, "@")]
			set[temp] = true
		}
	}

	// copy into return slice
	namespaces := make([]string, len(set))
	i := 0
	for k := range set {
		namespaces[i] = k
		i++
	}
	return namespaces, nil
}
