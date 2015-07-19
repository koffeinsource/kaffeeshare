package email

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"

	"appengine"
)

// Required to be able to pass different kind of headers in the following functions
type emailHeader interface {
	Get(key string) string
}

// extracts the body of an email
func extractBody(c appengine.Context, header emailHeader, bodyReader io.Reader) (*email, error) {
	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if mediaType[:4] == "text" {
		c.Infof("extractBody: found text")
		return extractTextBody(c, header, bodyReader)
	}

	if mediaType[:9] == "multipart" {
		c.Infof("extractBody: multipart")
		return extractMimeBody(c, params["boundary"], bodyReader)
	}

	return nil, fmt.Errorf("Unsupported content type: %s", contentType)
}

// read through the varios multiple parts
func extractMimeBody(c appengine.Context, boundary string, bodyReader io.Reader) (*email, error) {
	var withError *email // stores an email parse with error

	mimeReader := multipart.NewReader(bodyReader, boundary)

	for {
		part, err := mimeReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		defer part.Close()

		result, err := extractBody(c, part.Header, part)

		// this means we tried to decode it, but are not sure
		// lets save this result and try the other parts before return this result
		if result != nil && err != nil {
			withError = result
			c.Infof("extractMimeBody: email guess with error %v", err)
			continue
		}

		if result != nil && result.ContentType[:4] == "text" {
			return result, nil
		}
	}

	if withError != nil {
		return withError, nil
	}

	return nil, fmt.Errorf("Could not parse any of the multiple parts:")
}

// Decode body text and store it in a string
func extractTextBody(c appengine.Context, header emailHeader, bodyReader io.Reader) (*email, error) {
	var returnee email
	encoding := header.Get("Content-Transfer-Encoding")
	c.Infof("extractTextBody encoding: %v", encoding)

	s, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	if encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(string(s))
		if err != nil {
			return nil, err
		}
		returnee.Body = string(b)
		returnee.ContentType = header.Get("Content-Type")
		return &returnee, nil
	}
	if encoding == "quoted-printable" {
		// https://stackoverflow.com/questions/24883742/how-to-decode-mail-body-in-go
		// looks like it will be in go 1.5
		// maybe wait until then?
		// TODO
	}

	// ok, let's guess this is just plain text and put it into a string
	returnee.Body = string(s)
	returnee.ContentType = header.Get("Content-Type")

	return &returnee, fmt.Errorf("Unsupported Content-Transfer-Encoding: %v", encoding)
}
