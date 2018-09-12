package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/gngeorgiev/checkmail"
)

var (
	fromEmailFlag string
	fromHostFlag  string
	timeoutFlag   string

	toFlag string
)

func init() {
	flag.StringVar(&fromEmailFlag, "from", "", "")
	flag.StringVar(&fromHostFlag, "host", "", "")
	flag.StringVar(&timeoutFlag, "timeout", "", "")

	flag.StringVar(&toFlag, "to", "", "")
}

func main() {
	flag.Parse()

	opts := make([]checkmail.CheckerOption, 0)
	if timeoutFlag != "" {
		d, err := time.ParseDuration(timeoutFlag)
		if err != nil {
			panic(err)
		}

		opts = append(opts, checkmail.Timeout(d))
	}
	opts = append(opts, checkmail.FromEmail(fromEmailFlag))
	opts = append(opts, checkmail.FromHost(fromHostFlag))

	checker := checkmail.NewChecker(opts...)

	if err := checker.ValidateFormat(toFlag); err != nil {
		fmt.Println(fmt.Sprintf("error validating email format - %s", err))
		return
	}

	var mx []*net.MX
	var err error
	if mx, err = checker.ValidateDNS(toFlag); err != nil {
		fmt.Println(fmt.Sprintf("error validating dns %s", err))
		return
	}

	if err := checker.ValidateSMTP(toFlag, mx); err != nil {
		fmt.Println(fmt.Sprintf("error validating host %s", err))
		return
	}

	fmt.Println("OK!")
}
