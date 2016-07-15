package emailparse

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"strings"

	"github.com/koffeinsource/kaffeeshare/data"
)

const (
	contentTypeText  = "text"
	contentTypeHTML  = "html"
	contentTypeMulti = "multipart"
	contentTypeImage = "image/"
)

// Required to be able to pass different kind of headers in the following functions
type emailHeader interface {
	Get(key string) string
}

// extracts the body of an email
func extract(con *data.Context, header emailHeader, bodyReader io.Reader) (*Email, error) {
	var em Email
	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediaType, contentTypeText) {
		con.Log.Debugf("extract: found text")
		t, err := extractText(con, header, bodyReader)
		if err != nil {
			con.Log.Errorf("Error while extracting text body: %v", err)
		}
		if t != nil {
			em.Texts = append(em.Texts, *t)
		}
		return &em, nil
	}

	if strings.HasPrefix(mediaType, contentTypeMulti) {
		con.Log.Debugf("extract: found multipart; recursion; mediaType: %v", mediaType)
		e, err := extractMime(con, params["boundary"], bodyReader)
		if err != nil {
			con.Log.Errorf("Error while extracting multipart: %v", err)
		}
		em.Images = append(em.Images, e.Images...)
		em.Texts = append(em.Texts, e.Texts...)
		return &em, err
	}

	if strings.HasPrefix(mediaType, contentTypeImage) {
		con.Log.Debugf("extract: image; mediaType: %v", mediaType)
		i, err := extractImage(con, header, bodyReader)
		if err != nil {
			con.Log.Errorf("Error while extracting image: %v", err)
		}
		if i != nil {
			em.Images = append(em.Images, *i)
		}
		return &em, err
	}

	con.Log.Debugf("extract: %v; ignoring", mediaType)
	return nil, nil
}

// read through the varios multiple parts
func extractMime(con *data.Context, boundary string, bodyReader io.Reader) (e Email, err error) {
	mimeReader := multipart.NewReader(bodyReader, boundary)

	for {
		part, errT := mimeReader.NextPart()
		if errT == io.EOF {
			break
		}
		if errT != nil {
			return e, err
		}
		defer part.Close()

		result, errT := extract(con, part.Header, part)
		if errT != nil {
			con.Log.Errorf("Error in while calling extract in extractMime: %v", err)
			err = errT
		}

		if result != nil {
			e.Texts = append(e.Texts, result.Texts...)
			e.Images = append(e.Images, result.Images...)
		}
	}

	return e, err
}

// Extract image body return it as a string
func extractImage(con *data.Context, header emailHeader, bodyReader io.Reader) (*ImageBody, error) {
	s, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	encoding := header.Get("Content-Transfer-Encoding")
	con.Log.Debugf("extractImage encoding: %v", encoding)

	var ret ImageBody
	ret.Body = s
	ret.Encoding = encoding
	return &ret, nil
}

// Decode body text and store it in a string
func extractText(con *data.Context, header emailHeader, bodyReader io.Reader) (*TextBody, error) {
	var ret TextBody

	s, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}
	con.Log.Debugf("header: %v", header)

	ret.ContentType = header.Get("Content-Type")

	encoding := header.Get("Content-Transfer-Encoding")
	con.Log.Debugf("extractTextBody encoding: %v", encoding)

	if encoding == "base64" {
		b, err := base64.StdEncoding.DecodeString(string(s))
		if err != nil {
			return nil, err
		}
		ret.Body = string(b)
		return &ret, nil
	}

	if encoding == "7bit" {
		// that is just US ASCII (7bit)
		// https://stackoverflow.com/questions/25710599/content-transfer-encoding-7bit-or-8-bit
		ret.Body = string(s)
		return &ret, nil
	}

	if encoding == "quoted-printable" {
		// https://stackoverflow.com/questions/24883742/how-to-decode-mail-body-in-go
		r := quotedprintable.NewReader(bytes.NewReader(s))
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}
		ret.Body = string(b)
		return &ret, nil
	}

	// ok, let's guess this is just plain text and put it into a string
	ret.Body = string(s)

	return &ret, nil
}
