package main

import "github.com/skip2/go-qrcode"

// QREncodeToSmallString encodes the given content into a QR code and returns
// a small string representation of the QR code art. It uses the 'qrcode' package's
// New and ToSmallString functions to generate and format the QR code.
//
// Parameters:
//     content (string): The content to be encoded into the QR code.
//     disableBorder (bool): If set to true, the border of the QR code will be disabled.
//     negative (bool): If set to true, the colors of the QR code art will be inverted.
//
// Returns:
//     string: A small string representation of the QR code art.
//     error: An error object indicating any errors that occurred during QR code generation.
//
// Usage:
//     qrArt, err := QREncodeToSmallString("Hello World", false, false)
func QREncodeToSmallString(content string, disableBorder bool, negative bool) (string, error) {
	var q *qrcode.QRCode
	q, err := qrcode.New(content, qrcode.Low)
	if err != nil {
		return "", err
	}

	if disableBorder {
		q.DisableBorder = true
	}

	art := q.ToSmallString(negative)
	return art, nil
}
