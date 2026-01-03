package providers

import (
	"context"
	"log"

	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioProvider struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioProvider(accountSID, authToken, fromNumber string) *TwilioProvider {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &TwilioProvider{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (p *TwilioProvider) SendSMS(ctx context.Context, to string, message string) (*service.DeliveryResult, error) {
	params := &openapi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(p.fromNumber)
	params.SetBody(message)

	resp, err := p.client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Twilio Error: %v", err)
		return &service.DeliveryResult{Success: false, ErrorMessage: err.Error()}, err
	}

	return &service.DeliveryResult{
		Success:    true,
		ProviderID: *resp.Sid,
	}, nil
}
