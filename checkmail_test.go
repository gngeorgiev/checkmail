package checkmail

import (
	"testing"
	"time"
)

const (
	_FromHost  = "checkmail.me"
	_FromEmail = "lansome-cowboy@gmail.com"
)

var (
	samples = []struct {
		mail    string
		format  bool
		account bool //host+user
	}{
		{mail: "florian@carrere.cc", format: true, account: true},
		{mail: " florian@carrere.cc", format: false, account: false},
		{mail: "florian@carrere.cc ", format: false, account: false},
		{mail: "test@912-wrong-domain902.com", format: true, account: false},
		{mail: "0932910-qsdcqozuioqkdmqpeidj8793@gmail.com", format: true, account: false},
		{mail: "@gmail.com", format: false, account: false},
		{mail: "test@gmail@gmail.com", format: false, account: false},
		{mail: "test test@gmail.com", format: false, account: false},
		{mail: " test@gmail.com", format: false, account: false},
		{mail: "test@wrong domain.com", format: false, account: false},
		{mail: "é&ààà@gmail.com", format: false, account: false},
		{mail: "admin@jalopyjournal.com", format: true, account: false},
		{mail: "admin@busyboo.com", format: true, account: false},
	}

	checker = NewChecker(
		FromHost(_FromHost),
		FromEmail(_FromEmail),
		Timeout(5*time.Second),
	)
)

func TestInitialize(t *testing.T) {
	if checker.FromEmail() != _FromEmail {
		t.Fatal(checker.FromEmail())
	}

	if checker.FromHost() != _FromHost {
		t.Fatal(checker.FromHost())
	}

	if checker.Timeout() != 5*time.Second {
		t.Fatal(checker.Timeout())
	}
}

func TestValidateHost(t *testing.T) {
	for _, s := range samples {
		if !s.format {
			continue
		}

		mx, err := checker.ValidateDNS(s.mail)
		if err != nil && s.account == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
			return
		}

		err = checker.ValidateSMTP(s.mail, mx)
		if err != nil && s.account == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.account == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}

func TestValidateFormat(t *testing.T) {
	for _, s := range samples {
		err := checker.ValidateFormat(s.mail)
		if err != nil && s.format == true {
			t.Errorf(`"%s" => unexpected error: "%v"`, s.mail, err)
		}
		if err == nil && s.format == false {
			t.Errorf(`"%s" => expected error`, s.mail)
		}
	}
}
