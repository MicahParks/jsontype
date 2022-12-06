package jsontype

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// errStringUnmarshal is an error message for when an expected JSON string could not be unmarshalled.
	errStringUnmarshal = "failed to unmarshal as a string"

	// errUnmarshalPackage is prepended to all unmarshal errors to make troubleshooting easier.
	errUnmarshalPackage = "github.com/MicahParks/jsontype JSON unmarshal error"
)

// Options is a set of options for a JSONType. It modifies the behavior of JSON marshal/unmarshal.
type Options struct {
	MailAddressAddressOnlyMarshal bool
	MailAddressLowerMarshal       bool
	MailAddressUpperMarshal       bool
	TimeFormatMarshal             string
	TimeFormatUnmarshal           string
}

// J is a set of common Go types that can be marshaled and unmarshalled with this package.
type J interface {
	*mail.Address | *regexp.Regexp | time.Duration | time.Time | *url.URL | uuid.UUID
}

// JSONType holds a generic J value. It can be used to marshal and unmarshal its value to and from JSON.
type JSONType[T J] struct {
	mux     sync.RWMutex
	options Options
	v       T
}

// New creates a new JSONType.
func New[T J](v T) *JSONType[T] {
	return &JSONType[T]{
		v: v,
	}
}

// NewWithOptions creates a new JSONType with options.
func NewWithOptions[T J](v T, options Options) *JSONType[T] {
	return &JSONType[T]{
		options: options,
		v:       v,
	}
}

// Get returns the held value.
func (j *JSONType[T]) Get() T {
	if j == nil {
		var t T
		return t
	}
	j.mux.RLock()
	defer j.mux.RUnlock()
	return j.v
}

// MarshalJSON helps implement the json.Marshaler interface.
func (j *JSONType[T]) MarshalJSON() ([]byte, error) {
	var s string
	switch v := any(j.Get()).(type) {
	case *mail.Address:
		if j.options.MailAddressAddressOnlyMarshal {
			s = v.Address
		} else {
			s = v.String()
		}
		if j.options.MailAddressLowerMarshal {
			s = strings.ToLower(s)
		} else if j.options.MailAddressUpperMarshal {
			s = strings.ToUpper(s)
		}
	case *regexp.Regexp:
		s = v.String()
	case time.Duration:
		s = v.String()
	case time.Time:
		format := time.RFC3339
		if j.options.TimeFormatMarshal != "" {
			format = j.options.TimeFormatMarshal
		}
		s = v.Format(format)
	case *url.URL:
		s = v.String()
	case uuid.UUID:
		s = v.String()
	}
	return json.Marshal(s)
}

// Options returns the options for the held value.
func (j *JSONType[T]) Options() Options {
	j.mux.RLock()
	defer j.mux.RUnlock()
	return j.options
}

// Set sets the held value.
func (j *JSONType[T]) Set(v T) {
	j.mux.Lock()
	j.v = v
	j.mux.Unlock()
}

// SetOptions sets the options for the held value.
func (j *JSONType[T]) SetOptions(options Options) {
	j.mux.Lock()
	j.options = options
	j.mux.Unlock()
}

// UnmarshalJSON helps implement the json.Unmarshaler interface.
func (j *JSONType[T]) UnmarshalJSON(bytes []byte) error {
	switch any(j.v).(type) {
	case *mail.Address:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errUnmarshalPackage, errStringUnmarshal, err)
		}
		addr, err := mail.ParseAddress(s)
		if err != nil {
			return fmt.Errorf("%s: failed to parse email address: %w", errUnmarshalPackage, err)
		}
		j.Set(any(addr).(T))
	case *regexp.Regexp:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errStringUnmarshal, errUnmarshalPackage, err)
		}
		re, err := regexp.Compile(s)
		if err != nil {
			return fmt.Errorf("%s: failed to compile regexp: %w", errUnmarshalPackage, err)
		}
		j.Set(any(re).(T))
	case time.Duration:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errUnmarshalPackage, errStringUnmarshal, err)
		}
		d, err := time.ParseDuration(s)
		if err != nil {
			return fmt.Errorf("%s: failed to parse duration: %w", errUnmarshalPackage, err)
		}
		j.Set(any(d).(T))
	case time.Time:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errUnmarshalPackage, errStringUnmarshal, err)
		}
		format := time.RFC3339
		if j.options.TimeFormatUnmarshal != "" {
			format = j.options.TimeFormatUnmarshal
		}
		t, err := time.Parse(format, s)
		if err != nil {
			return fmt.Errorf("%s: failed to parse time: %w", errUnmarshalPackage, err)
		}
		j.Set(any(t).(T))
	case *url.URL:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errUnmarshalPackage, errStringUnmarshal, err)
		}
		u, err := url.Parse(s)
		if err != nil {
			return fmt.Errorf("%s: failed to parse url: %w", errUnmarshalPackage, err)
		}
		j.Set(any(u).(T))
	case uuid.UUID:
		var s string
		err := json.Unmarshal(bytes, &s)
		if err != nil {
			return fmt.Errorf("%s: %s: %w", errUnmarshalPackage, errStringUnmarshal, err)
		}
		u, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("%s: failed to parse uuid: %w", errUnmarshalPackage, err)
		}
		j.Set(any(u).(T))
	}
	return nil
}
