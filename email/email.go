package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/zeromicro/go-zero/core/conf"
	"net/smtp"
)

func init() {
	conf.MustLoad("./etc/email.yaml", &cfg)
	fmt.Printf("Read email config: %+v\n", cfg)
}

var (
	cfg Config
)

type Config struct {
	From     string
	Username string
	Password string
	Host     string
	Addr     string
}

func SendEmail(to, subject, content string) error {
	em := email.NewEmail()
	em.From = cfg.From
	em.To = []string{to}
	em.Subject = subject
	em.HTML = []byte(content)
	return em.Send(cfg.Addr, smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host))
}
