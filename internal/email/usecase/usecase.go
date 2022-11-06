package usecase

import (
	"bytes"
	"fmt"
	"gitflow/config"
	"gitflow/internal/email/models"
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"html/template"
	"log"
	"net/smtp"
	"reflect"
	"time"
)

type EmailUseCase struct {
}

func (e *EmailUseCase) InitTask(input []models.UserMail, config config.Config) {
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

	sendMail := func(from string, password string, host string, port string, email []string, body string) error {
		auth := smtp.PlainAuth("", from, password, host)
		if err != nil {
			return fmt.Errorf("send mail: %v\n", err)
		}

		err = smtp.SendMail(host+":"+port, auth, from, email, []byte(body))

		if err != nil {
			return fmt.Errorf("send mail: %v\n", err)
		}
		return nil
	}

	cli.Register("worker.mailing", sendMail)

	cli.StartWorker()
	time.Sleep(100 * time.Second)
	cli.StopWorker()

}
func (e *EmailUseCase) DoTasks(input models.UserMail, tmpl *template.Template) error {
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
	cli, err1 := gocelery.NewCeleryClient(
		gocelery.NewRedisBroker(redisPool),
		&gocelery.RedisCeleryBackend{Pool: redisPool},
		10,
	)
	if err1 != nil {
		log.Fatal(err1)
	}
	taskName := "worker.mailing"
	//  func(from []string, email string, body []byte) error {

	template, err := tmpl.ParseFiles("client/template.html")
	if err != nil {
		fmt.Println(err)
		return err
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", input, mimeHeaders)))

	err = template.Execute(&body, input)

	var asyncResult, err2 = cli.Delay(taskName, "english2two@gmail.com", "yobruuhyjlcfdmua", "smtp.gmail.com", "587", input.Emails, body.String())
	if err2 != nil {
		panic(err2)
	}
	res, err := asyncResult.Get(10 * time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	log.Printf("result: %+v of type %+v", res, reflect.TypeOf(res))
	return nil
}
