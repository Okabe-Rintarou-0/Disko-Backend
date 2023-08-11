package email

import (
	"disko/utils"
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
)

func init() {
	utils.ReadConfig("./etc/email.yaml", &cfg)
	fmt.Printf("Read email config: %+v\n", cfg)
}

var (
	cfg Config
)

type Config struct {
	From     string `yaml:"from"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Addr     string `yaml:"addr"`
}

func SendEmail(to, subject, content string) error {
	em := email.NewEmail()
	em.From = cfg.From
	em.To = []string{to}
	em.Subject = subject
	em.HTML = []byte(content)
	return em.Send(cfg.Addr, smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host))
}
