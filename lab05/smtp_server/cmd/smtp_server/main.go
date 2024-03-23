package main

import (
	"fmt"
	gomail "gopkg.in/gomail.v2"
	"strconv"
)

func makeMailer() (*gomail.Dialer, error) {
	smtp := "127.0.0.1"
	port, err := strconv.Atoi("2500")
	if err != nil {
		fmt.Println("error parsing smtp port.", err)
		return nil, err
	}

	name := ""
	pwd := ""

	return gomail.NewDialer(smtp, port, name, pwd), nil
}

func send() error {
	mailer, err := makeMailer()
	if err != nil {
		return fmt.Errorf("can't create mailer. %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "st072182@student.spbu.ru")
	m.SetHeader("To", "st072182@student.spbu.ru")
	m.SetHeader("Subject", "MAILER_SUBJECT_HTML")
	//m.SetBody("text/plain", "HI MKN!")
	m.SetBody("text", `<head><title>Hi MKN!</title></head>
		<body>
			<h1>hi MKN!</h1>
			<p>hi MKn!</p>
			<p><i>BY!</i><br>Mike</p>
		</body>`)

	err = mailer.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("can't send digest. %v", err)
	}

	return nil
}

func main() {
	err := send()
	if err != nil {
		fmt.Println(err)
	}

}
