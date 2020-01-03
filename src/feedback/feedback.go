package feedback

import (
	"config"
	"controllers"
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

var server = smtpServer{"smtp.gmail.com", "587"}

func CreateMessage(req controllers.FeedbackRequest) string {
	return fmt.Sprintf("New feedback from %s related to %s:%v!\n\n %s",
		req.Email, req.RelatedTo.Type, req.RelatedTo.Id, req.Text)
}

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
