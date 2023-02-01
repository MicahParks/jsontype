package jsontype_test

import (
	"bytes"
	"encoding/json"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/MicahParks/jsontype"
)

func TestEmailAddress(t *testing.T) {
	const rawAddress = "example@example.com"
	addr, err := mail.ParseAddress(rawAddress)
	if err != nil {
		t.Fatalf("failed to parse address: %s", err)
	}

	// Success marshal.
	{
		b, err := json.Marshal(jsontype.New(addr))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		var s string
		err = json.Unmarshal(b, &s)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if s != addr.String() {
			t.Fatalf("invalid JSON produced for address")
		}

		opts := jsontype.Options{
			MailAddressAddressOnlyMarshal: true,
		}
		b, err = json.Marshal(jsontype.NewWithOptions(addr, opts))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if string(b) != encloseDoubleQuotes(addr.Address) {
			t.Fatalf("invalid JSON produced for address only")
		}

		const mixCaseAddr = "mixCASE@example.com"
		mixCase, err := mail.ParseAddress(mixCaseAddr)
		if err != nil {
			t.Fatalf("failed to parse uppercase address: %s", err)
		}
		opts.MailAddressLowerMarshal = true
		b, err = json.Marshal(jsontype.NewWithOptions(mixCase, opts))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if string(b) != encloseDoubleQuotes(strings.ToLower(mixCaseAddr)) {
			t.Fatalf("invalid JSON produced for lowercase address")
		}

		opts.MailAddressLowerMarshal = false
		opts.MailAddressUpperMarshal = true
		b, err = json.Marshal(jsontype.NewWithOptions(mixCase, opts))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if string(b) != encloseDoubleQuotes(strings.ToUpper(mixCaseAddr)) {
			t.Fatalf("invalid JSON produced for uppercase address")
		}
	}

	// Failure marshal.
	{

	}

	// Success unmarshal.
	unmarshal := jsontype.New(&mail.Address{})
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes(rawAddress)), &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if unmarshal.Get().String() != addr.String() {
			t.Fatalf("invalid address unmarshalled")
		}
	}

	// Failure unmarshal.
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes("no at symbol")), &unmarshal)
		if err == nil {
			t.Fatalf("expected error for invalid address")
		}
	}
}

func TestRegexp(t *testing.T) {
	const rawRegexp = ".*"
	validJSON := json.RawMessage(encloseDoubleQuotes(rawRegexp))
	r, err := regexp.Compile(rawRegexp)
	if err != nil {
		t.Fatalf("failed to compile regexp: %s", err)
	}

	// Success marshal.
	{
		b, err := json.Marshal(jsontype.New(r))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if !bytes.Equal(b, validJSON) {
			t.Fatalf("invalid JSON produced for regexp")
		}
	}

	// Failure marshal.
	{

	}

	// Success unmarshal.
	unmarshal := jsontype.New(&regexp.Regexp{})
	{
		err = json.Unmarshal(validJSON, &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if unmarshal.Get().String() != r.String() {
			t.Fatalf("invalid URL unmarshalled")
		}
	}

	// Failure unmarshal.
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes("*")), &unmarshal)
		if err == nil {
			t.Fatalf("expected error for invalid regexp")
		}
	}
}

func TestTimeDuration(t *testing.T) {
	const rawDuration = "1h0m0s"
	validJSON := json.RawMessage(encloseDoubleQuotes(rawDuration))
	d, err := time.ParseDuration(rawDuration)
	if err != nil {
		t.Fatalf("failed to parse duration: %s", err)
	}

	// Success marshal.
	{
		b, err := json.Marshal(jsontype.New(d))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if !bytes.Equal(b, validJSON) {
			t.Fatalf("invalid JSON produced for duration")
		}
	}

	// Failure marshal.
	{

	}

	// Success unmarshal.
	unmarshal := jsontype.New(time.Duration(0))
	{
		err = json.Unmarshal(validJSON, &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if unmarshal.Get().String() != d.String() {
			t.Fatalf("invalid duration unmarshalled")
		}
	}

	// Failure unmarshal.
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes("")), &unmarshal)
		if err == nil {
			t.Fatalf("expected error for invalid duration")
		}
	}
}

func TestTimeTime(t *testing.T) {
	const rawTime = "2022-10-04T00:00:00Z"
	validJSON := json.RawMessage(encloseDoubleQuotes(rawTime))
	tt, err := time.Parse(time.RFC3339, rawTime)
	if err != nil {
		t.Fatalf("failed to parse time: %s", err)
	}
	tn := tt.Add(time.Nanosecond)
	validJSONNS := encloseDoubleQuotes("2022-10-04T00:00:00.000000001Z")

	// Success marshal.
	{
		b, err := json.Marshal(jsontype.New(tt))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if !bytes.Equal(b, validJSON) {
			t.Fatalf("invalid JSON produced for time")
		}

		opts := jsontype.Options{
			TimeFormatMarshal: time.RFC3339Nano,
		}
		b, err = json.Marshal(jsontype.NewWithOptions(tn, opts))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if string(b) != validJSONNS {
			t.Fatalf("invalid JSON produced for time")
		}
	}

	// Failure marshal.
	{

	}

	// Success unmarshal.
	unmarshal := jsontype.New(time.Time{})
	{
		err = json.Unmarshal(validJSON, &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if unmarshal.Get().String() != tt.String() {
			t.Fatalf("invalid time unmarshalled")
		}

		opts := jsontype.Options{
			TimeFormatUnmarshal: time.RFC3339Nano,
		}
		unmarshal = jsontype.NewWithOptions(time.Time{}, opts)
		err = json.Unmarshal(json.RawMessage(validJSONNS), &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal JSON: %s", err)
		}
		if unmarshal.Get().String() != tn.String() {
			t.Fatalf("invalid time unmarshalled")
		}
	}

	// Failure unmarshal.
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes("")), &unmarshal)
		if err == nil {
			t.Fatalf("expected error for invalid time")
		}
	}
}

func TestURL(t *testing.T) {
	const rawURL = "https://github.com"
	validJSON := json.RawMessage(encloseDoubleQuotes(rawURL))
	u, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("failed to parse URL: %s", err)
	}

	// Success marshal.
	{
		b, err := json.Marshal(jsontype.New(u))
		if err != nil {
			t.Fatalf("failed to marshal struct: %s", err)
		}
		if !bytes.Equal(b, validJSON) {
			t.Fatalf("invalid JSON produced for URL")
		}
	}

	// Failure marshal.
	{

	}

	// Success unmarshal.
	unmarshal := jsontype.New(&url.URL{})
	{
		err = json.Unmarshal(validJSON, &unmarshal)
		if err != nil {
			t.Fatalf("failed to unmarshal struct: %s", err)
		}
		if unmarshal.Get().String() != u.String() {
			t.Fatalf("invalid URL unmarshalled")
		}
	}

	// Failure unmarshal.
	{
		err = json.Unmarshal(json.RawMessage(encloseDoubleQuotes(string(rune(0x7f)))), &unmarshal)
		if err == nil {
			t.Fatal("expected error for URL with invalid character")
		}
	}
}

func encloseDoubleQuotes(s string) string {
	return `"` + s + `"`
}
