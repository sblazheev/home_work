package amqp

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"                                                 //nolint:depguard
	amqp "github.com/rabbitmq/amqp091-go"                                   //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common" //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
)

type ClientInterface interface {
	Publish(routingKey string, body []byte) error
	Consume(queueName string) (<-chan []byte, error)
	Close() error
}

type Client struct {
	m               *sync.Mutex
	queueName       string
	infolog         *logger.Logger
	errlog          *logger.Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
)

var (
	errNotConnected  = errors.New("not connected to a server")
	errAlreadyClosed = errors.New("already closed: not connected to the server")
	errShutdown      = errors.New("client is shutting down")
)

func New(queueName, addr string, infolog, errlog *logger.Logger) *Client {
	client := Client{
		m:         &sync.Mutex{},
		infolog:   infolog,
		errlog:    errlog,
		queueName: queueName,
		done:      make(chan bool),
	}
	go client.handleReconnect(addr)
	return &client
}

func (client *Client) handleReconnect(addr string) {
	for {
		client.m.Lock()
		client.isReady = false
		client.m.Unlock()

		client.infolog.Info("attempting to connect", "addr", common.ExtractAddr(addr))

		conn, err := client.connect(addr)
		if err != nil {
			client.errlog.Error("failed to connect. Retrying...", "addr", common.ExtractAddr(addr))

			select {
			case <-client.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		if done := client.handleReInit(conn); done {
			break
		}
	}
}

func (client *Client) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	client.changeConnection(conn)
	client.infolog.Info("connected", "addr", client.connection.RemoteAddr())
	return conn, nil
}

func (client *Client) handleReInit(conn *amqp.Connection) bool {
	for {
		client.m.Lock()
		client.isReady = false
		client.m.Unlock()

		err := client.init(conn)
		if err != nil {
			client.errlog.Error("failed to initialize channel, retrying...", "addr", conn.RemoteAddr())
			select {
			case <-client.done:
				return true
			case <-client.notifyConnClose:
				client.infolog.Info("connection closed, reconnecting...", "addr", conn.RemoteAddr())
				return false
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-client.done:
			return true
		case <-client.notifyConnClose:
			client.infolog.Info("connection closed, reconnecting...", "addr", conn.RemoteAddr())
			return false
		case <-client.notifyChanClose:
			client.infolog.Info("hannel closed, re-running init..", "addr", conn.RemoteAddr())
		}
	}
}

func (client *Client) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	_, err = ch.QueueDeclare(
		client.queueName,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	client.changeChannel(ch)
	client.m.Lock()
	client.isReady = true
	client.m.Unlock()
	client.infolog.Info("client init done", "addr", conn.RemoteAddr())

	return nil
}

func (client *Client) changeConnection(connection *amqp.Connection) {
	client.connection = connection
	client.notifyConnClose = make(chan *amqp.Error, 1)
	client.connection.NotifyClose(client.notifyConnClose)
}

func (client *Client) changeChannel(channel *amqp.Channel) {
	client.channel = channel
	client.notifyChanClose = make(chan *amqp.Error, 1)
	client.notifyConfirm = make(chan amqp.Confirmation, 1)
	client.channel.NotifyClose(client.notifyChanClose)
	client.channel.NotifyPublish(client.notifyConfirm)
}

func (client *Client) Push(data []byte) error {
	client.m.Lock()
	if !client.isReady {
		client.m.Unlock()
		return errNotConnected
	}
	client.m.Unlock()
	for {
		err := client.UnsafePush(data)
		if err != nil {
			client.errlog.Error("push failed. Retrying...", "addr", client.connection.RemoteAddr())
			select {
			case <-client.done:
				return errShutdown
			case <-time.After(resendDelay):
			}
			continue
		}
		confirm := <-client.notifyConfirm
		if confirm.Ack {
			return nil
		}
	}
}

func (client *Client) UnsafePush(data []byte) error {
	client.m.Lock()
	if !client.isReady {
		client.m.Unlock()
		return errNotConnected
	}
	client.m.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.channel.PublishWithContext(
		ctx,
		"",               // Exchange
		client.queueName, // Routing key
		false,            // Mandatory
		false,            // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (client *Client) Consume() (<-chan amqp.Delivery, error) {
	client.m.Lock()
	if !client.isReady {
		client.m.Unlock()
		return nil, errNotConnected
	}
	client.m.Unlock()

	if err := client.channel.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	); err != nil {
		return nil, err
	}

	return client.channel.Consume(
		client.queueName,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

// Close will cleanly shut down the channel and connection.
func (client *Client) Close() error {
	client.m.Lock()
	defer client.m.Unlock()

	if !client.isReady {
		return errAlreadyClosed
	}
	close(client.done)
	err := client.channel.Close()
	if err != nil {
		return err
	}
	err = client.connection.Close()
	if err != nil {
		return err
	}

	client.isReady = false
	return nil
}
