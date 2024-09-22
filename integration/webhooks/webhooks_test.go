package webhooks

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type WebhooksTestSuite struct {
	suite.Suite
}

func TestWebhooks(t *testing.T) {
	suite.Run(t, new(WebhooksTestSuite))
}

func (s *WebhooksTestSuite) TestAddingExtensionToContext() {
	s.assertWebhookCall("extensionAddedToContext")
}

func (s *WebhooksTestSuite) TestRemovingExtensionFromContext() {
	extensionInstanceID := uuid.NewString()

	s.assertWebhookCall(
		"extensionAddedToContext",
		WithExtensionInstanceID(extensionInstanceID),
	)

	s.assertWebhookCall(
		"instanceRemovedFromContext",
		WithExtensionInstanceID(extensionInstanceID),
	)
}

func (s *WebhooksTestSuite) TestRemovingExtensionFromContextTwice() {
	extensionInstanceID := uuid.NewString()

	s.assertWebhookCall(
		"extensionAddedToContext",
		WithExtensionInstanceID(extensionInstanceID),
	)

	s.assertWebhookCall(
		"instanceRemovedFromContext",
		WithExtensionInstanceID(extensionInstanceID),
	)

	s.assertWebhookCall(
		"instanceRemovedFromContext",
		WithExtensionInstanceID(extensionInstanceID),
	)
}

func (s *WebhooksTestSuite) TestRemovingUnknownExtensionFromContext() {
	s.assertWebhookCall("instanceRemovedFromContext")
}
