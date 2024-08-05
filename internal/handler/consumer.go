// Package handler
package handler

import (
	"context"

	"gitlab.privy.id/privypass/privypass-boilerplate/internal/appctx"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	uContract "gitlab.privy.id/privypass/privypass-boilerplate/internal/ucase/contract"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/awssqs"
)

// SQSConsumerHandler sqs consumer message processor handler
func SQSConsumerHandler(msgHandler uContract.MessageProcessor) awssqs.MessageProcessorFunc {
	return func(decoder *awssqs.MessageDecoder) error {
		return msgHandler.Serve(context.Background(), &appctx.ConsumerData{
			Body:        []byte(*decoder.Body),
			Key:         []byte(*decoder.MessageId),
			ServiceType: consts.ServiceTypeConsumer,
		})
	}
}
