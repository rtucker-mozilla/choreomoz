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
	"encoding/json"
)

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Inside Error")
		fmt.Println("%s: %s", msg, err)
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
/*func parseQueueName(input_hostname string) (err error){
		return input_hostnane, err
}

func parseBindingKey(input_hostname string) (err error){
		return input_hostnane, err
}*/


func main() {
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	exchange := "chorizo"
	exchangeType := "topic"
	HOSTNAME, _ := os.Hostname()
	fmt.Println(HOSTNAME)
	queue := "robs-macbook-pro"
	bindingKey := "robs-macbook-pro.host"
	//masterBindingKey := "master.robs-macbook-pro.host"
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

func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {

	// This function dials, connects, declares, publishes, and tears down,
	// all in one go. In a real service, you probably want to maintain a
	// long-lived connection as state, and publish against that.

	log.Printf("dialing %q", amqpURI)
	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	log.Printf("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	log.Printf("got Channel, declaring %q Exchange (%q)", exchangeType, exchange)
	if err := channel.ExchangeDeclare(
		exchange,     // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.

	log.Printf("declared Exchange, publishing %dB body (%q)", len(body), body)

	return nil
}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string) (*Consumer, error) {
	c := &Consumer{
		conn:    nil,
		channel: nil,
		tag:     ctag,
		done:    make(chan error),
	}
	fmt.Println(exchange, exchangeType)
	masterQueueName := fmt.Sprintf("master.%s", queueName)
	fmt.Println(masterQueueName)

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
	_, err = c.channel.QueueDeclare(
		masterQueueName, // name of the queue
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

	masterRoutingKey := "master.robs-macbook-pro.host"
	if err = c.channel.QueueBind(
		masterQueueName, // name of the queue
		masterRoutingKey,        // bindingKey
		exchange,   // sourceExchange
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	log.Printf("Queue bound to Exchange, starting Consume (consumer tag %q)", c.tag)
	// forever chan to block
	forever := make(chan bool)

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
	go handle(deliveries, c.done, c)

	<-forever



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

func handle(deliveries <-chan amqp.Delivery, done chan error, c *Consumer) {
	for d := range deliveries {
		log.Printf(
			"got %dB delivery: [%v] %q",
			len(d.Body),
			d.DeliveryTag,
			d.Body,
		)
		d.Ack(true)
		res := ParseCommand{}
		json.Unmarshal(d.Body, &res)
		body_string := "test"
		respObj, err := res.ExecuteCommand()
		if err != nil {
			body_string = "Error occurred"
		} else {
			body_string, _ = respObj.Response()
		}
		masterRoutingKey := "master.robs-macbook-pro.host"
		puberr := c.channel.Publish(
			"chorizo",   // publish to an exchange
			masterRoutingKey, // routing to 0 or more queues
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Headers:         amqp.Table{},
				ContentType:     "text/plain",
				ContentEncoding: "",
				Body:            []byte(body_string),
				DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
				Priority:        9,              // 0-9
				// a bunch of application/implementation-specific fields
			},
		) 
		fmt.Println("Here")
		if puberr != nil {
			fmt.Println(fmt.Errorf("Exchange Publish: %s", puberr))
		}
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
