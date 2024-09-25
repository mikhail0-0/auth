package mail

import (
	"auth/src/config"
	"fmt"
	"log"
	"net/smtp"
)

func SendMessage(email, message string) error {
	log.Printf(
		"Trying to send message '%v' from %v to %v\n",
		message,
		config.SmtpEmail,
		email,
	)

	c, err := smtp.Dial(fmt.Sprintf("%v:%v", config.SmtpHost, config.SmtpPort))
	if err != nil {
		return err
	}

	if err := c.Mail(config.SmtpEmail); err != nil {
		return err
	}
	if err := c.Rcpt(email); err != nil {
		return err
	}

	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(wc, message)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}
