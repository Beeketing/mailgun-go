package mailgun

import (
	"testing"

	"encoding/json"

	"github.com/facebookgo/ensure"
)

func TestEmailValidation(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
	validator, err := NewEmailValidatorFromEnv()
	ensure.Nil(t, err)

	ev, err := validator.ValidateEmail("foo@mailgun.com", false)
	ensure.Nil(t, err)

	ensure.True(t, ev.IsValid)
	ensure.DeepEqual(t, ev.MailboxVerification, "")
	ensure.False(t, ev.IsDisposableAddress)
	ensure.False(t, ev.IsRoleAddress)
	ensure.True(t, ev.Parts.DisplayName == "")
	ensure.DeepEqual(t, ev.Parts.LocalPart, "foo")
	ensure.DeepEqual(t, ev.Parts.Domain, "mailgun.com")
	ensure.True(t, ev.Reason == "")
}

func TestParseAddresses(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
	validator, err := NewEmailValidatorFromEnv()
	ensure.Nil(t, err)

	addressesThatParsed, unparsableAddresses, err := validator.ParseAddresses(
		"Alice <alice@example.com>",
		"bob@example.com",
		"example.com")
	ensure.Nil(t, err)
	hittest := map[string]bool{
		"Alice <alice@example.com>": true,
		"bob@example.com":           true,
	}
	for _, a := range addressesThatParsed {
		ensure.True(t, hittest[a])
	}
	ensure.True(t, len(unparsableAddresses) == 1)
}

func TestUnmarshallResponse(t *testing.T) {
	payload := []byte(`{
		"address": "some_email@aol.com",
		"did_you_mean": null,
		"is_disposable_address": false,
		"is_role_address": false,
		"is_valid": true,
		"mailbox_verification": "unknown",
		"parts":
		{
			"display_name": null,
			"domain": "aol.com",
			"local_part": "some_email"
		},
		"reason": null
	}`)
	var ev EmailVerification
	err := json.Unmarshal(payload, &ev)
	ensure.Nil(t, err)

	ensure.True(t, ev.IsValid)
	ensure.DeepEqual(t, ev.MailboxVerification, "unknown")
	ensure.False(t, ev.IsDisposableAddress)
	ensure.False(t, ev.IsRoleAddress)
	ensure.True(t, ev.Parts.DisplayName == "")
	ensure.DeepEqual(t, ev.Parts.LocalPart, "some_email")
	ensure.DeepEqual(t, ev.Parts.Domain, "aol.com")
	ensure.True(t, ev.Reason == "")
}
