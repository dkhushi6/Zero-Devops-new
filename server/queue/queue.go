
package queue

import amqp "github.com/rabbitmq/amqp091-go"

func SetUpQueues(queueChannel *amqp.Channel) error{

		err := queueChannel.ExchangeDeclare(		
			"deploy.dlx",
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)

	if err != nil{
		return err
	}	


	// DEAD LETTER QUEUE FOR JOBS
	_,err = queueChannel.QueueDeclare(
		"deploy.jobs.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil{
		return err
	}	

	err = queueChannel.QueueBind("deploy.jobs.dlq", "deploy.jobs.dlq", "deploy.dlx", false, nil)

	if err != nil{
		return err
	}	

	args_jobs := amqp.Table{
		"x-dead-letter-exchange":    "deploy.dlx",
		"x-dead-letter-routing-key": "deploy.jobs.dlq",
	}

	// JOBS QUEUE

	_,err = queueChannel.QueueDeclare(
		"deploy.jobs",
		true,
		false,
		false,
		false,
		args_jobs,
	)

	if err != nil{
		return err
	}

	// DEAD LETTER QUEUE FOR STATUS
	_,err = queueChannel.QueueDeclare(
		"deploy.status.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil{
		return err
	}	

	err = queueChannel.QueueBind("deploy.status.dlq", "deploy.status.dlq", "deploy.dlx", false, nil)

	if err != nil{
		return err
	}	

	args_status := amqp.Table{
		"x-dead-letter-exchange":    "deploy.dlx",
		"x-dead-letter-routing-key": "deploy.status.dlq",
	}

	// STATUS QUEUE
	_,err = queueChannel.QueueDeclare(
		"deploy.status",
		true,
		false,
		false,
		false,
		args_status,
	)

	if err != nil{
		return err
	}

	return nil

}