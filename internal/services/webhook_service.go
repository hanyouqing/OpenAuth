package services

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WebhookService struct {
	db     *gorm.DB
	logger *logrus.Logger
	client *http.Client
}

func NewWebhookService(db *gorm.DB, logger *logrus.Logger) *WebhookService {
	return &WebhookService{
		db:     db,
		logger: logger,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *WebhookService) List() ([]models.Webhook, error) {
	var webhooks []models.Webhook
	if err := s.db.Find(&webhooks).Error; err != nil {
		return nil, err
	}
	return webhooks, nil
}

func (s *WebhookService) Get(id uint64) (*models.Webhook, error) {
	var webhook models.Webhook
	if err := s.db.First(&webhook, id).Error; err != nil {
		return nil, err
	}
	return &webhook, nil
}

func (s *WebhookService) Create(webhook *models.Webhook) error {
	// Generate secret if not provided
	if webhook.Secret == "" {
		webhook.Secret = uuid.New().String()
	}
	return s.db.Create(webhook).Error
}

func (s *WebhookService) Update(id uint64, data map[string]interface{}) error {
	return s.db.Model(&models.Webhook{}).Where("id = ?", id).Updates(data).Error
}

func (s *WebhookService) Delete(id uint64) error {
	return s.db.Delete(&models.Webhook{}, id).Error
}

func (s *WebhookService) Trigger(event string, payload map[string]interface{}) error {
	var webhooks []models.Webhook
	if err := s.db.Where("enabled = ?", true).Find(&webhooks).Error; err != nil {
		return err
	}

	for _, webhook := range webhooks {
		// Check if webhook subscribes to this event
		subscribed := false
		for _, e := range webhook.Events {
			if e == event || e == "*" {
				subscribed = true
				break
			}
		}
		if !subscribed {
			continue
		}

		// Create webhook event
		webhookEvent := models.WebhookEvent{
			WebhookID: webhook.ID,
			Event:     event,
			Payload:   payload,
			Status:    "pending",
		}
		s.db.Create(&webhookEvent)

		// Send webhook asynchronously
		go s.sendWebhook(webhook, webhookEvent, payload)
	}

	return nil
}

func (s *WebhookService) sendWebhook(webhook models.Webhook, webhookEvent models.WebhookEvent, payload map[string]interface{}) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		s.logger.WithError(err).Error("Failed to marshal webhook payload")
		return
	}

	// Sign payload
	signature := s.signPayload(webhook.Secret, payloadBytes)

	// Create request
	req, err := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		s.updateWebhookEvent(webhookEvent.ID, "failed", err.Error(), 0)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-Event", webhookEvent.Event)
	req.Header.Set("X-Webhook-ID", fmt.Sprintf("%d", webhook.ID))

	// Send request
	resp, err := s.client.Do(req)
	if err != nil {
		s.updateWebhookEvent(webhookEvent.ID, "failed", err.Error(), webhookEvent.Attempts+1)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		s.updateWebhookEvent(webhookEvent.ID, "success", "", webhookEvent.Attempts+1)
	} else {
		s.updateWebhookEvent(webhookEvent.ID, "failed", fmt.Sprintf("HTTP %d", resp.StatusCode), webhookEvent.Attempts+1)
	}
}

func (s *WebhookService) signPayload(secret string, payload []byte) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(payload)
	return hex.EncodeToString(h.Sum(nil))
}

func (s *WebhookService) updateWebhookEvent(id uint64, status, response string, attempts int) {
	s.db.Model(&models.WebhookEvent{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":   status,
		"response": response,
		"attempts": attempts,
		"updated_at": time.Now(),
	})
}

func (s *WebhookService) GetEvents(webhookID uint64, limit int) ([]models.WebhookEvent, error) {
	var events []models.WebhookEvent
	query := s.db.Where("webhook_id = ?", webhookID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}
