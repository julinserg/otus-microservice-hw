package delivery_amqp

import (
	"context"
	"encoding/json"
	"fmt"

	amqp_pub "github.com/julinserg/julinserg/OtusMicroserviceHomeWork/hw08_saga/internal/amqp/pub"
	amqp_settings "github.com/julinserg/julinserg/OtusMicroserviceHomeWork/hw08_saga/internal/amqp/settings"
	amqp_sub "github.com/julinserg/julinserg/OtusMicroserviceHomeWork/hw08_saga/internal/amqp/sub"
	order_app "github.com/julinserg/julinserg/OtusMicroserviceHomeWork/hw08_saga/internal/order/app"
	"github.com/streadway/amqp"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type SrvDelivery interface {
	CreateDeliveryOperation(order order_app.Order) error
	RevertDeliveryOperation(idOrder string, statusOrder string) error
}

type SrvDeliveryAMQP struct {
	logger      Logger
	srvDelivery SrvDelivery
	pub         amqp_pub.AmqpPub
	uri         string
}

func New(logger Logger, uri string) *SrvDeliveryAMQP {
	return &SrvDeliveryAMQP{
		logger: logger,
		pub:    *amqp_pub.New(logger),
		uri:    uri,
	}
}

func (a *SrvDeliveryAMQP) SetService(srvDelivery SrvDelivery) {
	a.srvDelivery = srvDelivery
}

func (a *SrvDeliveryAMQP) StartReceiveOrder(ctx context.Context) error {
	conn, err := amqp.Dial(a.uri)
	if err != nil {
		return err
	}
	c := amqp_sub.New("SrvDeliveryAMQPOrder", conn, a.logger)
	msgs, err := c.Consume(ctx, amqp_settings.QueueOrderDeliveryService, amqp_settings.ExchangeOrder,
		"direct", amqp_settings.RoutingKeyDeliveryService)
	if err != nil {
		return err
	}

	a.logger.Info("start consuming order...")

	for m := range msgs {
		order := order_app.Order{}
		json.Unmarshal(m.Data, &order)
		if err != nil {
			return err
		}
		a.logger.Info(fmt.Sprintf("receive new message:%+v\n", order))

		err := a.srvDelivery.CreateDeliveryOperation(order)
		if err != nil {
			a.logger.Warn("Error CreateDeliveryOperation: " + err.Error())
		}
	}
	return nil
}

func (a *SrvDeliveryAMQP) StartReceiveStatus(ctx context.Context) error {
	conn, err := amqp.Dial(a.uri)
	if err != nil {
		return err
	}
	c := amqp_sub.New("SrvDeliveryAMQPStatus", conn, a.logger)
	msgs, err := c.Consume(ctx, amqp_settings.QueueStatusDeliveryService, amqp_settings.ExchangeStatus, "direct", "")
	if err != nil {
		return err
	}

	a.logger.Info("start consuming status...")

	for m := range msgs {
		notifyEvent := order_app.OrderEvent{}
		json.Unmarshal(m.Data, &notifyEvent)
		if err != nil {
			return err
		}
		a.logger.Info(fmt.Sprintf("receive new order status update event:%+v\n", notifyEvent))
		err := a.srvDelivery.RevertDeliveryOperation(notifyEvent.Id, notifyEvent.Status)
		if err != nil {
			a.logger.Warn("Error RevertDeliveryOperation: " + err.Error())
		}
	}
	return nil
}

func (a *SrvDeliveryAMQP) PublishStatus(idOrder string, statusOrder string) error {
	orderStatusEvent := order_app.OrderEvent{Id: idOrder, Status: statusOrder}
	orderStatusStr, err := json.Marshal(orderStatusEvent)
	if err != nil {
		return err
	}
	if err := a.pub.Publish(a.uri, amqp_settings.ExchangeStatus, "direct",
		"", string(orderStatusStr), true); err != nil {
		return err
	}
	a.logger.Info("publish order status for queue is OK ( OrderId: " + orderStatusEvent.Id + ")")
	return nil
}
