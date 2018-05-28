package checkmail

import (
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

type SmtpError struct {
	Err error
}

func (e SmtpError) Error() string {
	return e.Err.Error()
}

func (e SmtpError) Code() string {
	return e.Err.Error()[0:3]
}

func NewSmtpError(err error) SmtpError {
	return SmtpError{
		Err: err,
	}
}

var (
	ErrBadFormat        = errors.New("invalid format")
	ErrUnresolvableHost = errors.New("unresolvable host")

	emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type CheckerOptions struct {
	fromHost, fromEmail string
	timeout             time.Duration
}

type CheckerOption func(*CheckerOptions)

var (
	FromHost = func(fromHost string) CheckerOption {
		return func(o *CheckerOptions) {
			o.fromHost = fromHost
		}
	}

	FromEmail = func(fromEmail string) CheckerOption {
		return func(o *CheckerOptions) {
			o.fromEmail = fromEmail
		}
	}

	Timeout = func(timeout time.Duration) CheckerOption {
		return func(o *CheckerOptions) {
			o.timeout = timeout
		}
	}
)

type Checker struct {
	opts CheckerOptions
}

func NewChecker(opts ...CheckerOption) Checker {
	c := Checker{CheckerOptions{}}

	for _, o := range opts {
		o(&c.opts)
	}
	if c.opts.fromEmail == "" || c.opts.fromHost == "" {
		panic("FromEmail and FromHost options are required")
	}
	if c.opts.timeout == 0 {
		c.opts.timeout = 10 * time.Second
	}

	return c
}

func (c Checker) FromHost() string {
	return c.opts.fromHost
}

func (c Checker) FromEmail() string {
	return c.opts.fromEmail
}

func (c Checker) Timeout() time.Duration {
	return c.opts.timeout
}

func (c Checker) ValidateFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return ErrBadFormat
	}
	return nil
}

func (c Checker) ValidateHost(email string) error {
	_, host := split(email)
	mx, err := net.LookupMX(host)
	if err != nil {
		return ErrUnresolvableHost
	}

	client, err := smtp.Dial(fmt.Sprintf("%s:%d", mx[0].Host, 25))
	if err != nil {
		return NewSmtpError(err)
	}
	defer client.Close()

	t := time.AfterFunc(c.Timeout(), func() { client.Close() })
	defer t.Stop()

	err = client.Hello(c.FromHost())
	if err != nil {
		return NewSmtpError(err)
	}

	err = client.Mail(c.FromEmail())
	if err != nil {
		return NewSmtpError(err)
	}

	err = client.Rcpt(email)
	if err != nil {
		return NewSmtpError(err)
	}

	return nil
}

func split(email string) (account, host string) {
	i := strings.LastIndexByte(email, '@')
	account = email[:i]
	host = email[i+1:]
	return
}
