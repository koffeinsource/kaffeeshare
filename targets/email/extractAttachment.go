package email

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/koffeinsource/kaffeeshare/data"
)

const (
	contentTypeImage    = "image/"
	contentTypeSMIMESig = "application/pkcs7-signature"
)

type imageBody struct {
	Body     []byte
	Encoding string
}

// extracts the body of an email
func extractAttachment(con *data.Context, header emailHeader, bodyReader io.Reader) ([]imageBody, error) {
	var images []imageBody
	contentType := header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(mediaType, contentTypeText) {
		con.Log.Debugf("extractAttachment: found text; ignoring; mediaType: %v", mediaType)
		return nil, nil
	}

	if strings.HasPrefix(mediaType, contentTypeSMIMESig) {
		con.Log.Debugf("extractAttachment: SMIME sig; ignoring; mediaType: %v", mediaType)
		return nil, nil
	}

	if strings.HasPrefix(mediaType, contentTypeMulti) {
		con.Log.Debugf("extractAttachment: found multipart; recursion; mediaType: %v", mediaType)
		is, err := extractMimeAttachment(con, params["boundary"], bodyReader)
		images = append(images, is...)
		return images, err
	}

	if strings.HasPrefix(mediaType, contentTypeImage) {
		con.Log.Debugf("extractAttachment: image; mediaType: %v", mediaType)
		i, err := extractImage(con, header, bodyReader)
		if err != nil {
			return nil, err
		}
		images = append(images, *i)
		return images, nil
	}

	return nil, fmt.Errorf("Unsupported content type: %s; media type: %v", contentType, mediaType)
}

// read through the varios multiple parts
func extractMimeAttachment(con *data.Context, boundary string, bodyReader io.Reader) ([]imageBody, error) {
	var images []imageBody // stores an email parse with error

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

		is, err := extractAttachment(con, part.Header, part)

		if err != nil {
			con.Log.Errorf("err: %v", err)
		} else {
			images = append(images, is...)
		}
	}

	return images, nil
}

// Extract image body return it as a string
func extractImage(con *data.Context, header emailHeader, bodyReader io.Reader) (*imageBody, error) {
	s, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	encoding := header.Get("Content-Transfer-Encoding")
	con.Log.Debugf("extractImage encoding: %v", encoding)

	var ret imageBody
	ret.Body = s
	ret.Encoding = encoding
	return &ret, nil
}
