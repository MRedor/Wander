package feedback

import (
	"config"
	"net/smtp"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) serverName() string {
	return s.host + ":" + s.port
}

var server = smtpServer{"smtp.gmail.com", "587"}

func Send(msg string) error {
	message := []byte(msg)

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
