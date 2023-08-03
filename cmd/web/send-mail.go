package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AdiAkhileshSingh15/bookmyroom/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
		}
		mailTemplate := string(data)

		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)

		//1//https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcSByKlEQXwjstWR7Fhuq1ku5zo88dp7yvcEI64rm-0ZQqgvSHHb
		//2//https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSE76YgwYNFaGj1D014zzldnoEq4nhFn5lSwTeJftKPshBFuZ7T

		if m.RoomID == 1 {
			msgToSend = strings.Replace(msgToSend, "[%image%]", "https://encrypted-tbn1.gstatic.com/images?q=tbn:ANd9GcSByKlEQXwjstWR7Fhuq1ku5zo88dp7yvcEI64rm-0ZQqgvSHHb", 1)
		} else {
			msgToSend = strings.Replace(msgToSend, "[%image%]", "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSE76YgwYNFaGj1D014zzldnoEq4nhFn5lSwTeJftKPshBFuZ7T", 1)
		}

		email.SetBody(mail.TextHTML, msgToSend)
	}

	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
	} else {
		app.InfoLog.Println("Email sent")
	}
}
