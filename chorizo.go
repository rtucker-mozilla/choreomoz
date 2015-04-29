package main

import (
	"fmt"
	"github.com/streadway/amqp"
	config "libchorizo/config"
	"parse_update_script"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Inside Error")
		fmt.Println("%s: %s", msg, err)
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
func main() {
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	exchange := "chorizo"
	exchangeType := "direct"
	HOSTNAME, _ := os.Hostname()
	queue := HOSTNAME
	bindingKey := "chorizo"
	consumerTag := ""
	c, err := NewConsumer("amqp://localhost:5672/", exchange, exchangeType, queue, bindingKey, consumerTag)
	if err != nil {
		fmt.Println("Could not create NewConsumer")
		fmt.Println(err)
	}
	if err := c.Shutdown(); err != nil {
		log.Fatalf("error during shutdown: %s", err)
	}
	//exec_path, _ := os.Getwd()
	// In https://github.com/rtucker-mozilla/chorizo/issues/15
	// going to specify a specific config file path
	exec_path := "/etc/chorizo"
	cfg := config.Config{}
	config := cfg.NewConfig(exec_path)
	fmt.Println(config)

}
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}
	fmt.Println(c)

	var err error

	log.Printf("dialing %q", amqpURI)
	c.conn, err = amqp.Dial(amqpURI)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
	}()

	log.Printf("got Connection, getting Channel")
	c.channel, err = c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring Exchange (%q)", exchange)
	if err = c.channel.ExchangeDeclare(
		exchange,     // name of the exchange
		exchangeType, // type
		false,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return nil, fmt.Errorf("Exchange Declare: %s", err)
	}

	log.Printf("declared Exchange, declaring Queue %q", queueName)
	queue, err := c.channel.QueueDeclare(
		queueName, // name of the queue
		false,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, key)

	if err = c.channel.QueueBind(
		queue.Name, // name of the queue
		key,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	for {
		deliveries, err := c.channel.Consume(
			queue.Name, // name
			c.tag,      // consumerTag,
			false,      // noAck
			false,      // exclusive
			false,      // noLocal
			false,      // noWait
			nil,        // arguments
		)
		if err != nil {
			return nil, fmt.Errorf("Queue Consume: %s", err)
		}
		go handle(deliveries, c.done)
		time.Sleep(1 * time.Second)

	}


	return c, nil
}

func (c *Consumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.channel.Cancel(c.tag, true); err != nil {
		return fmt.Errorf("Consumer cancel failed: %s", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer log.Printf("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.done
}

func handle(deliveries <-chan amqp.Delivery, done chan error) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(true)
	}
	//log.Printf("handle: deliveries channel closed")
	// Here we detect the actionable item
	/* Choices Appear to be:
	   StartUpdate
	   EndUpdate
	   Reboot
	   ExecScript
	*/
	// Need to figure out logic for restarting
	done <- nil
}

type UpdateScriptResponse struct {
	ret_code      int
	system_id     int
	stdout        string
	stderr        string
	is_start      bool
	is_end        bool
	update_script *parse_update_script.UpdateScript
	db            *sql.DB
	api_url       string
}

// SystemReboot executes a shell command to reboot the host
func SystemReboot(exec_reboot bool) bool {
	// Sleep for 2 seconds to give time for the API post
	if exec_reboot == true {
		cmd := exec.Command("/sbin/shutdown", "-r", "now")
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
	time.Sleep(2 * time.Second)
	return exec_reboot
}

func InterpolateConfigOption(exec_path string, config_item string) (retval string) {
	retval = fmt.Sprintf("%s/%s", exec_path, config_item)
	return
}
