package usecase

import (
	"bytes"
	"fmt"
	"gitflow/internal/email/models"
	"html/template"
	"log"
	"net/smtp"
	"strings"
	"time"

	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
)

type EmailUseCase struct {
}

func (e *EmailUseCase) StartWorker() {
	//rep
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
	//rep
	cli, err := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		10,
	)
	if err != nil {
		log.Fatal(err)
	}

	sendMail := func(from string, password string, host string, port string, email string, body string) string {
		auth := smtp.PlainAuth("", from, password, host)
		if err != nil {
			fmt.Println(err)
			return "error"
		}
		fmt.Println(email)
		data := strings.Split(email, ",")
		for i := 0; i < len(data); i++ {
			fmt.Println(data[i], "!")
		}
		err = smtp.SendMail(host+":"+port, auth, from, data, []byte(body))

		if err != nil {
			fmt.Println(err)
			return "error"
		}
		return "SUCCESS"
	}

	cli.Register("worker.mailing", sendMail)

	cli.StartWorker()
	cli.WaitForStopWorker()
}

func (e *EmailUseCase) StartClient(input models.UserMail, tmpl *template.Template) (string, error) {
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	cli, err := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		10,
	)

	if err != nil {
		return "", err
	}
	taskName := "worker.mailing"

	template, err := tmpl.ParseFiles("client/template.html")
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", input.Title, mimeHeaders)))

	_ = template.Execute(&body, input)
	data := strings.Join(input.Emails, ",")

	asyncResult, err := cli.Delay(taskName, "english2two@gmail.com", "ztbmxouvonzpicvw", "smtp.gmail.com", "587", data, body.String())
	if err != nil {
		return "", err
	}
	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		return "", err
	}

	return res.(string), nil
}
