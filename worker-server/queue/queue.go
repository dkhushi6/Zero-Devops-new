package queue

import (
	"encoding/json"
	"log"
	"Zero_Devops/worker_server/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error , msg string){
	if err != nil{
		log.Panicf("%s : %s",msg,err)
	}
}

type queueUsecase struct {
	queueClient	*domain.RabbitMQ
}

func NewQueueUsecase(conn *amqp.Connection , channel *amqp.Channel) domain.QueueUsecase{
	return &queueUsecase{
		queueClient: &domain.RabbitMQ{
			Conn: conn,
			Channel: channel,
		},
	}
}

func (r *queueUsecase) Close(){
	queueClient := r.queueClient
	queueClient.Conn.Close()
	queueClient.Channel.Close()
}

func (r *queueUsecase) Channel() *amqp.Channel {
	return r.queueClient.Channel
}

func (r *queueUsecase) SetUpQueues() error{

		err := r.queueClient.Channel.ExchangeDeclare(		
			"deploy.dlx",
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)

	if err != nil{
		failOnError(err,"Failed to create the Exchange")
		return err
	}	

	_,err = r.queueClient.Channel.QueueDeclare(
		"deploy.jobs.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil{
		failOnError(err,"Failed to  Create DLQ")
		return err
	}	

	err = r.queueClient.Channel.QueueBind("deploy.jobs.dlq", "deploy.jobs.dlq", "deploy.dlx", false, nil)

	if err != nil{
		failOnError(err,"Failed to Bind")
		return err
	}	

	args := amqp.Table{
		"x-dead-letter-exchange":    "deploy.dlx",
		"x-dead-letter-routing-key": "deploy.jobs.dlq",
	}

	_,err = r.queueClient.Channel.QueueDeclare(
		"deploy.jobs",
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil{
		failOnError(err,"Failed to declare a queue")
		return err
	}

	return nil

}

func (r* queueUsecase) PublishJob(job domain.DeployJob) error{
	body, err := json.Marshal(job)

	if err != nil{
		failOnError(err,"Failed to Receive Jobs")
		return err
	}

	return r.queueClient.Channel.Publish(
		"",
		"deploy.jobs",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp.Persistent,
			Body: body,
		},
	)
}

