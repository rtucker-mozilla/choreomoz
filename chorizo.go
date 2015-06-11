package main

import (
	"fmt"
	"github.com/streadway/amqp"
	config "libchorizo/config"
	util "libchorizo/util"
	"parse_update_script"
	"database/sql"
	"log"
	"os"
	"os/exec"
	"time"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"encoding/json"
)
type HelloResp struct {
	Hash 						string
	Hostname    				string
	Command 					string
	ReturnString 				string
	GroupId 					int
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println("Inside Error")
		fmt.Println("%s: %s", msg, err)
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func SendHello(masterRoutingKey string, c *Consumer){
	hello_resp := HelloResp{}
	hello_resp.Hash = "xxx"
	hello_resp.Hostname, _ = os.Hostname()
	hello_resp.Command = "hello_resp"
	hello_resp.ReturnString = "hello_resp"
	body_string, _ := json.Marshal(hello_resp)
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
	fmt.Println(string(body_string))
	if puberr != nil {
		fmt.Println(fmt.Errorf("Exchange Publish: %s", puberr))
	}

}

func main() {
	//conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	exec_path := "/etc/chorizo"
	cfg := config.Config{}
	config := cfg.NewConfig(exec_path)
	useTls := config.Main.UseTls
	exchange := "chorizo"
	exchangeType := "topic"
	HOSTNAME, _ := os.Hostname()
	fmt.Println(HOSTNAME)
	queue := util.HostnameToQueueName(HOSTNAME)
	bindingKey := util.HostnameToBindingKey(HOSTNAME)
	consumerTag := ""
	conn_prefix := ""
	if useTls == true {
		conn_prefix = "amqps"
	} else {
		conn_prefix = "amqp"
	}
	conn_string := fmt.Sprintf("%s://%s:%s@%s:%s/", 
		conn_prefix,
		config.Main.RabbitmqUser,
		config.Main.RabbitmqPass,
		config.Main.RabbitmqHost,
		config.Main.RabbitmqPort)
	c, err := NewConsumer(conn_string, exchange, exchangeType, queue, bindingKey, consumerTag, useTls)
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
	fmt.Println(config)

}
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	done    chan error
}


/*func publish(amqpURI, exchange, exchangeType, routingKey, body string, reliable bool) error {

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
}*/

func GetConnection(amqpURI string, useTls bool) (*amqp.Connection, error){
	cfg := new(tls.Config)
	cfg.RootCAs = x509.NewCertPool()
	if ca, err := ioutil.ReadFile("/etc/chorizo/ssl/mozilla-cacert.pem"); err == nil {
    	cfg.RootCAs.AppendCertsFromPEM(ca)
	}


	if useTls == true {
		fmt.Println("dialing tls")
		conn, err := amqp.DialTLS(amqpURI, cfg)
		return conn, err
	} else {
		conn, err := amqp.Dial(amqpURI)
		return conn, err
	}

}

func NewConsumer(amqpURI, exchange, exchangeType, queueName, key, ctag string, useTls bool) (*Consumer, error) {
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
	c.conn, err = GetConnection(amqpURI, useTls)
	if err != nil {
		return nil, fmt.Errorf("Dial: %s", err)
	}

	go func() {
		// Let supervisord restart the process
		fmt.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		os.Exit(2)
	}()

	log.Printf("got Connection, getting Channel")
	if c.channel == nil {
		c.channel, err = c.conn.Channel()
	}
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

	masterRoutingKey := fmt.Sprintf("master.%s.host", queue.Name)
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
	// Announce to master queue that we have connected
	// forever chan to block
	SendHello(masterRoutingKey, c)
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
		body_string := ""
		respObj, err := res.ExecuteCommand()
		if err != nil {
			body_string = "Error occurred"
		} else {
			body_string, _ = respObj.Response()
		}
		masterRoutingKey := fmt.Sprintf("master.%s", d.RoutingKey)
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
		if puberr != nil {
			fmt.Println(fmt.Errorf("Exchange Publish: %s", puberr))
		}
		if res.Command == "start_reboot" {
			fmt.Println("Start Reboot Now")

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
