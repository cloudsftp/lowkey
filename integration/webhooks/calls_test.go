package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cloudsftp/lowkey/integration/pointer"
)

func (s *WebhooksTestSuite) assertWebhookCall(
	webhookType string,
	options ...WebhookCallOption,
) {
	s.T().Helper()

	config := defaultWebhookCallConfiguration()

	for _, o := range options {
		o(&config)
	}

	var requestBody bytes.Buffer
	requestBodyEncoder := json.NewEncoder(&requestBody)
	err := requestBodyEncoder.Encode(config.WebhookCallRequest)
	s.NoError(err, "encoutered error when encoding the webhook call request")

	localDevHost := os.Getenv("LOCAL_DEV_URL")

	response, err := http.Post(
		fmt.Sprintf("%s/internal/calls/%s", localDevHost, webhookType),
		"application/json",
		&requestBody,
	)
	s.NoError(err, "encountered error when requesting webhook call")
	defer response.Body.Close()

	s.Equal(config.StatusCode, response.StatusCode)

	var extensionResponse ExtensionResponse
	responseBodyDecoder := json.NewDecoder(response.Body)
	err = responseBodyDecoder.Decode(&extensionResponse)
	s.NoError(err, "encountered error when decoding the extension response")

	s.Equal(config.ExtensionResponse, extensionResponse)
}

const (
	ExtensionID   = "0c3f42aa-536b-4562-abaf-ccbd91f53bda"
	ContributorID = "e77eb8f9-9b25-4e13-b0ab-92f6eaa8d3f7"
)

type WebhookCallRequest struct {
	ExtensionID               *string  `json:"extensionId"`
	ContributorID             *string  `json:"contributorId"`
	Context                   *string  `json:"context"`
	ContextAggregateID        *string  `json:"contextAggregateId"`
	ExtensionInstanceID       *string  `json:"extensionInstanceId"`
	ConsentedScopes           []string `json:"consentedScopes"`
	ExtensionInstanceDisabled bool     `json:"extensionInstanceDisabled"`
}

type ExtensionResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Successful bool   `json:"successful"`
}

type WebhookCallConfiguration struct {
	WebhookCallRequest
	ExtensionResponse
}

func defaultWebhookCallConfiguration() WebhookCallConfiguration {
	return WebhookCallConfiguration{
		WebhookCallRequest: WebhookCallRequest{
			ExtensionID:               pointer.Of(ExtensionID),
			ContributorID:             pointer.Of(ContributorID),
			Context:                   pointer.Of("customer"),
			ExtensionInstanceDisabled: false,
		},
		ExtensionResponse: ExtensionResponse{
			StatusCode: 200,
			Body:       "Ok",
			Successful: true,
		},
	}
}

type WebhookCallOption func(*WebhookCallConfiguration)

func WithExtensionInstanceID(extensionInstanceID string) WebhookCallOption {
	return func(config *WebhookCallConfiguration) {
		config.ExtensionInstanceID = &extensionInstanceID
	}
}
