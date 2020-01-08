package feedback

import (
	"config"
	"fmt"
	"net/smtp"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

var server = smtpServer{config.Config.Email.Host, config.Config.Email.Port}

func createMessage(email, text string, relatedId int, relatedType string) string {
	return fmt.Sprintf("New feedback from %s related to %s:%v!\n\n %s",
		email, relatedType, relatedId, text)
}

func Send(email, text string, relatedId int, relatedType string) error {
	message := []byte(createMessage(email, text, relatedId, relatedType))

	err := smtp.SendMail(
		server.serverName(),
		smtp.PlainAuth("", config.Config.Email.Login, config.Config.Email.Password, server.host),
		config.Config.Email.Login, //from
		[]string{config.Config.Email.Login}, //to
		message,
	)

	if err != nil {
		return err
	}
	return nil
}
