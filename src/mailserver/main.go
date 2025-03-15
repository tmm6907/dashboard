package main

import (
	"log"
	"net"
	"time"

	smtp "github.com/emersion/go-smtp"
	"github.com/gofiber/fiber/v2"
	"github.com/tmm6907/mailserver/httpServer"
	"github.com/tmm6907/mailserver/mail"
)

func main() {
	backend := mail.Backend{}
	mailServer := smtp.NewServer(backend)
	if mailServer == nil {
		log.Fatalln("Failed to make mail server")
	}
	smtpPort := ":2525"
	mail.ConfigureServer(mailServer, smtpPort, "example.com", 10*time.Second, 10*time.Second, true)

	log.Println("Starting SMTP server on", smtpPort)
	go func(m *smtp.Server) {
		listener, err := net.Listen("tcp", m.Addr)
		if err != nil {
			log.Fatal(err)
		}

		if err = m.Serve(listener); err != nil {
			log.Printf("Mail Server ERROR: %s", err.Error())
		}
	}(mailServer)

	server := fiber.New()

	server.Get("/", httpServer.FeedHandler)

	go func(s *fiber.App) {
		if err := server.Listen(":8888"); err != nil {
			log.Printf("HTTP Server ERROR: %s", err.Error())
		}
	}(server)

	select {}
}
