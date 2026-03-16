package routes

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// Errors
var (
	errInvalidRequest = errors.New("invalid request")
	errMissingSrcDst  = errors.New("src and dst are required")
	errInvalidFormat  = errors.New("expected format lng,lat")
	errInvalidLong    = errors.New("invalid longitude")
	errInvalidLang    = errors.New("invalid latitude")
)

type Request struct {
	Src string   `json:"src"`
	Dst []string `json:"dst"`
}

// parseRequest reads and normalizes query params from the incoming HTTP request.
func parseRequest(r *http.Request) Request {
	dst := r.URL.Query()["dst"]
	for i := range dst {
		dst[i] = strings.TrimSpace(dst[i])
	}

	return Request{
		Src: strings.TrimSpace(r.URL.Query().Get("src")),
		Dst: dst,
	}
}

// Validate ensures src and all dst values are present and follow lng,lat format.
func (r Request) Validate() error {
	if r.Src == "" || len(r.Dst) == 0 {
		return errors.Join(errInvalidRequest, errMissingSrcDst)
	}

	if err := validate(r.Src); err != nil {
		return errors.Join(errInvalidRequest, err)
	}

	for _, part := range r.Dst {
		if err := validate(part); err != nil {
			return errors.Join(errInvalidRequest, err)
		}
	}

	return nil
}

// validate checks that a coordinate is provided as two float values in lng,lat order.
func validate(raw string) error {
	parts := strings.Split(raw, ",")
	if len(parts) != 2 {
		return errInvalidFormat
	}

	_, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return errInvalidLong
	}

	_, err = strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return errInvalidLang
	}

	return nil
}
