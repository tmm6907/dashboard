package mail

import (
	"bytes"
	"html"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	smtp "github.com/emersion/go-smtp"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type Session struct {
	From string
	To   string
	db   *sqlx.DB
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.From = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	s.To = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return err
	}
	msg, err := mail.ReadMessage(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	subject := msg.Header.Get("Subject")
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return err
	}
	bodyStr := strings.ReplaceAll(string(body), "\r\n", "\n")
	bodyStr = strings.ReplaceAll(bodyStr, "\r", "\n")
	bodyStr = html.UnescapeString(bodyStr)
	_, err = s.db.Exec(
		"INSERT OR IGNORE INTO mail (sender, recipient, subject, body) VALUES (?, ?, ?, ?);",
		s.From, s.To, subject, bodyStr,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

type Backend struct{}

func populateDB(db *sqlx.DB) error {

	file, err := os.ReadFile("./build.sql")
	if err != nil {
		return err
	}
	if _, err = db.Exec(string(file)); err != nil {
		return err
	}
	return nil
}

func (b Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	db, err := sqlx.Open("sqlite", "mail.db")
	if err != nil {
		return &Session{}, err
	}
	err = populateDB(db)
	if err != nil {
		return &Session{}, err
	}
	return &Session{db: db}, nil
}

func ConfigureServer(server *smtp.Server,
	port string, domain string, readTimeout time.Duration,
	writeTimeout time.Duration, allowInsecure bool,
) {
	server.Addr = port // Change this to 25 if running as root
	server.Domain = domain
	server.ReadTimeout = readTimeout
	server.WriteTimeout = writeTimeout
	server.AllowInsecureAuth = allowInsecure
}
