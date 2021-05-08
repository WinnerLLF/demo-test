package push

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

var p *producer

type producer struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	Address  string
	exchange string
	key      string
}
type ReChargeMsg struct {
	RechargeAddress string `json:"recharge_address"`
	FromAddress     string `json:"from_address"`
	Cid             string `json:"cid"`
	Height          uint64 `json:"height"`
	Amount          string `json:"amount"`
	Status          int    `json:"status"`
}
type PConfig struct {
	Address  string
	Exchange string
	Key      string
}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func consumerRun() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover:%v\n", err)
		}
		time.Sleep(5 * time.Second)
		go consumerRun()
	}()

	conn, err := amqp.Dial("amqp://ipfsmain:ipfsmain@39.105.56.95:5672/")
	failOnError(err, "Failed to connect to RabbitMQ!")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel!")

	go func() {
		fmt.Printf("Listening to exit!\n")

		connNotify := conn.NotifyClose(make(chan *amqp.Error))
		channelNotify := ch.NotifyClose(make(chan *amqp.Error))

		select {
		case err := <-connNotify:
			if err != nil {
				for {
					fmt.Printf("Cleaning messages in progress!\n")
					select {
					case err := <-connNotify:
						fmt.Printf("rabbitmq consumer - connection NotifyClose: %v\n", err)
						if err == nil {
							time.Sleep(5 * time.Second)
							fmt.Printf("Is to restart!--------\n")
							go consumerRun()
							fmt.Printf("The restart command is complete!----------\n")
							return
						} else {
							fmt.Printf("err = %v\n", err)
						}
					case <-time.After(3 * time.Second):
						fmt.Printf("Message wait has timed out and is restarting!\n")
						_ = conn.Close()
						_ = ch.Close()
						time.Sleep(5 * time.Second)
						go consumerRun()
						return
					}
				}
			}
		case err := <-channelNotify:
			if err != nil {
				for {
					fmt.Printf("Cleaning messages in progress!\n")
					select {
					case err := <-channelNotify:
						fmt.Printf("rabbitmq consumer - channel NotifyClose: %v\n", err)
						if err == nil {
							time.Sleep(5 * time.Second)
							go consumerRun()
							return
						}
					case <-time.After(3 * time.Second):
						_ = conn.Close()
						_ = ch.Close()
						time.Sleep(5 * time.Second)
						go consumerRun()
						return
					}
				}
			}
		}
	}()

	q, err := ch.QueueDeclare(
		"recharge_notice", // name
		true,              // durable
		false,             // delete when usused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	err = ch.QueueBind("recharge_notice", "", "recharge_notice", false, nil)
	if err != nil {
		fmt.Printf("QueueBind err:%v\n", err)
		return
	}
	for d := range msgs {
		fmt.Printf("get a message:%v\n", string(d.Body))

		var req string
		_ = d.Ack(true)
		if err := json.Unmarshal(d.Body, &req); err != nil {
			fmt.Printf("解析错误：%v\n", err)
			return
		}
		fmt.Printf("数据展示：%v\n", req)
	}

}

func sendMsg() {
	conn, err := amqp.Dial("amqp://ipfsmain:ipfsmain@39.105.56.95:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"recharge_notice", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	failOnError(err, "Failed to declare a queue")
	body := "{\"recharge_address\":\"t152oeuctqo7seqxr3dmgp4exbag6tzat2jjxph6i\",\"from_address\":\"f3taase67itpl4z5qs3nrdkregerotwcgzyep6n6l5x2g2fd6vw22efx773vneitpa3lkgmxlq5yrbc76am4oq\",\"cid\":\"bafy2bzacecqlt5ocp3woly7fufg5ztwrzl4w4jksalmahis3sqjxqphz3zv3s\",\"height\":100043,\"amount\":\"50000000000000000000\",\"status\":0}"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}

func connMq() error {
	var err error

	fmt.Printf("dialing:%v\n", p.Address)
	p.conn, err = amqp.Dial(p.Address)
	if err != nil {
		fmt.Printf("Dial err:%v\n", err)
		return err
	}

	fmt.Printf("got Connection, getting Channel\n")
	p.channel, err = p.conn.Channel()
	if err != nil {
		fmt.Printf("Channel: %v\n", err)
		return err
	}
	return nil
}

func (p *producer) publish(body []byte) error {
	err := p.channel.Publish(
		"",
		p.exchange,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		fmt.Printf("Message sending exception:%v\n", err)
		return err
	}

	return nil
}

func Init_Producer(config *PConfig) {
	newProducer(config.Address, config.Exchange, config.Key)
	err := connMq()
	if err != nil {
		panic("Initialize the MQ producer exception:" + err.Error())
	}
}

func newProducer(Address, exchange, key string) {
	p = &producer{
		conn:     nil,
		channel:  nil,
		Address:  Address,
		exchange: exchange,
		key:      key,
	}
}

func PushMsg(data *ReChargeMsg) error {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Marshal err:%v\n", err)
		return err
	}
	fmt.Printf("数据展示信息：%v\n", string(body))
	if err := p.publish(body); err != nil {
		fmt.Printf("Send message exception err:%v\n", err)
		return err
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}
