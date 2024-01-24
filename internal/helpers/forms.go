package helpers

import (
	"net/url"

	"github.com/monoculum/formam/v3"
)

var decoder *formam.Decoder = nil

func Decode(values url.Values, dst interface{}) (err error) {
	decoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
	err = decoder.Decode(values, dst)
	return
}
