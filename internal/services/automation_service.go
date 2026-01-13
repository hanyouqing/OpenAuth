package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hanyouqing/openauth/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AutomationService struct {
	db       *gorm.DB
	logger   *logrus.Logger
	Services *Services
}

func NewAutomationService(db *gorm.DB, logger *logrus.Logger) *AutomationService {
	return &AutomationService{
		db:     db,
		logger: logger,
	}
}

func (s *AutomationService) SetServices(services *Services) {
	s.Services = services
}

// CreateWorkflow creates a new automation workflow
func (s *AutomationService) CreateWorkflow(name, description string, trigger models.AutomationTrigger, actions []models.AutomationAction, priority int) (*models.AutomationWorkflow, error) {
	triggerJSON, err := json.Marshal(trigger)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal trigger: %w", err)
	}

	actionsJSON, err := json.Marshal(actions)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal actions: %w", err)
	}

	workflow := &models.AutomationWorkflow{
		Name:        name,
		Description: description,
		Enabled:     true,
		Trigger:     string(triggerJSON),
		Actions:     string(actionsJSON),
		Priority:    priority,
	}

	if err := s.db.Create(workflow).Error; err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return workflow, nil
}

// GetWorkflow gets a workflow by ID
func (s *AutomationService) GetWorkflow(id uint64) (*models.AutomationWorkflow, error) {
	var workflow models.AutomationWorkflow
	if err := s.db.First(&workflow, id).Error; err != nil {
		return nil, err
	}
	return &workflow, nil
}

// ListWorkflows lists all workflows
func (s *AutomationService) ListWorkflows(enabled *bool) ([]models.AutomationWorkflow, error) {
	var workflows []models.AutomationWorkflow
	query := s.db
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}
	if err := query.Order("priority DESC, created_at DESC").Find(&workflows).Error; err != nil {
		return nil, err
	}
	return workflows, nil
}

// UpdateWorkflow updates a workflow
func (s *AutomationService) UpdateWorkflow(id uint64, name, description string, trigger *models.AutomationTrigger, actions []models.AutomationAction, enabled *bool, priority *int) error {
	updates := make(map[string]interface{})
	if name != "" {
		updates["name"] = name
	}
	if description != "" {
		updates["description"] = description
	}
	if trigger != nil {
		triggerJSON, err := json.Marshal(trigger)
		if err != nil {
			return fmt.Errorf("failed to marshal trigger: %w", err)
		}
		updates["trigger"] = string(triggerJSON)
	}
	if actions != nil {
		actionsJSON, err := json.Marshal(actions)
		if err != nil {
			return fmt.Errorf("failed to marshal actions: %w", err)
		}
		updates["actions"] = string(actionsJSON)
	}
	if enabled != nil {
		updates["enabled"] = *enabled
	}
	if priority != nil {
		updates["priority"] = *priority
	}

	return s.db.Model(&models.AutomationWorkflow{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteWorkflow deletes a workflow
func (s *AutomationService) DeleteWorkflow(id uint64) error {
	return s.db.Delete(&models.AutomationWorkflow{}, id).Error
}

// TriggerWorkflow triggers a workflow execution
func (s *AutomationService) TriggerWorkflow(workflowID uint64, input map[string]interface{}) (*models.AutomationExecution, error) {
	workflow, err := s.GetWorkflow(workflowID)
	if err != nil {
		return nil, err
	}

	if !workflow.Enabled {
		return nil, fmt.Errorf("workflow is disabled")
	}

	// Create execution
	inputJSON, _ := json.Marshal(input)
	execution := &models.AutomationExecution{
		WorkflowID: workflowID,
		Status:     "pending",
		Input:      string(inputJSON),
	}
	if err := s.db.Create(execution).Error; err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}

	// Execute workflow asynchronously
	go s.executeWorkflow(execution, workflow)

	return execution, nil
}

// executeWorkflow executes a workflow
func (s *AutomationService) executeWorkflow(execution *models.AutomationExecution, workflow *models.AutomationWorkflow) {
	now := time.Now()
	execution.Status = "running"
	execution.StartedAt = &now
	s.db.Save(execution)

	// Parse actions
	var actions []models.AutomationAction
	if err := json.Unmarshal([]byte(workflow.Actions), &actions); err != nil {
		s.logger.WithError(err).Error("Failed to parse workflow actions")
		execution.Status = "failed"
		execution.Error = fmt.Sprintf("Failed to parse actions: %v", err)
		s.db.Save(execution)
		return
	}

	// Parse input
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(execution.Input), &input); err != nil {
		s.logger.WithError(err).Error("Failed to parse execution input")
		execution.Status = "failed"
		execution.Error = fmt.Sprintf("Failed to parse input: %v", err)
		s.db.Save(execution)
		return
	}

	// Execute each action
	output := make(map[string]interface{})
	for i, action := range actions {
		result, err := s.executeAction(action, input)
		if err != nil {
			s.logger.WithError(err).Errorf("Action %d failed: %s", i, action.Type)
			if action.OnError == "stop" {
				execution.Status = "failed"
				execution.Error = fmt.Sprintf("Action %d failed: %v", i, err)
				s.db.Save(execution)
				return
			} else if action.OnError == "retry" {
				// Simple retry logic (in production, use a proper retry mechanism)
				result, err = s.executeAction(action, input)
				if err != nil {
					if action.OnError == "stop" {
						execution.Status = "failed"
						execution.Error = fmt.Sprintf("Action %d failed after retry: %v", i, err)
						s.db.Save(execution)
						return
					}
				}
			}
			// Continue on error if on_error is "continue"
		}
		output[fmt.Sprintf("action_%d", i)] = result
	}

	// Save output
	outputJSON, _ := json.Marshal(output)
	execution.Output = string(outputJSON)
	execution.Status = "completed"
	completedAt := time.Now()
	execution.CompletedAt = &completedAt
	s.db.Save(execution)
}

// executeAction executes a single action
func (s *AutomationService) executeAction(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	switch action.Type {
	case "send_email":
		return s.executeSendEmail(action, input)
	case "create_user":
		return s.executeCreateUser(action, input)
	case "assign_role":
		return s.executeAssignRole(action, input)
	case "remove_role":
		return s.executeRemoveRole(action, input)
	case "update_user":
		return s.executeUpdateUser(action, input)
	case "webhook":
		return s.executeWebhook(action, input)
	default:
		return nil, fmt.Errorf("unknown action type: %s", action.Type)
	}
}

// executeSendEmail executes a send_email action
func (s *AutomationService) executeSendEmail(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	email, ok := action.Config["email"].(string)
	if !ok {
		email, _ = input["email"].(string)
	}
	if email == "" {
		return nil, fmt.Errorf("email address is required")
	}

	subject, _ := action.Config["subject"].(string)
	body, _ := action.Config["body"].(string)

	if s.Services != nil && s.Services.Notification != nil {
		if err := s.Services.Notification.SendEmail(email, subject, body); err != nil {
			return nil, fmt.Errorf("failed to send email: %w", err)
		}
	}

	return map[string]interface{}{"status": "sent", "email": email}, nil
}

// executeCreateUser executes a create_user action
func (s *AutomationService) executeCreateUser(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	if s.Services == nil || s.Services.User == nil {
		return nil, fmt.Errorf("user service not available")
	}

	username, _ := input["username"].(string)
	email, _ := input["email"].(string)
	password, _ := input["password"].(string)

	if username == "" || email == "" {
		return nil, fmt.Errorf("username and email are required")
	}

	user, err := s.Services.User.Create(username, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return map[string]interface{}{"user_id": user.ID, "username": user.Username}, nil
}

// executeAssignRole executes an assign_role action
func (s *AutomationService) executeAssignRole(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	if s.Services == nil || s.Services.Role == nil {
		return nil, fmt.Errorf("role service not available")
	}

	userID, ok := input["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id is required")
	}

	roleName, _ := action.Config["role"].(string)
	if roleName == "" {
		return nil, fmt.Errorf("role name is required")
	}

	if err := s.Services.Role.AssignRoleToUser(uint64(userID), roleName); err != nil {
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return map[string]interface{}{"user_id": uint64(userID), "role": roleName}, nil
}

// executeRemoveRole executes a remove_role action
func (s *AutomationService) executeRemoveRole(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	if s.Services == nil || s.Services.Role == nil {
		return nil, fmt.Errorf("role service not available")
	}

	userID, ok := input["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id is required")
	}

	roleName, _ := action.Config["role"].(string)
	if roleName == "" {
		return nil, fmt.Errorf("role name is required")
	}

	if err := s.Services.Role.RemoveRoleFromUser(uint64(userID), roleName); err != nil {
		return nil, fmt.Errorf("failed to remove role: %w", err)
	}

	return map[string]interface{}{"user_id": uint64(userID), "role": roleName}, nil
}

// executeUpdateUser executes an update_user action
func (s *AutomationService) executeUpdateUser(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	if s.Services == nil || s.Services.User == nil {
		return nil, fmt.Errorf("user service not available")
	}

	userID, ok := input["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id is required")
	}

	updates := make(map[string]interface{})
	if status, ok := action.Config["status"].(string); ok {
		updates["status"] = status
	}
	if email, ok := action.Config["email"].(string); ok {
		updates["email"] = email
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no updates specified")
	}

	user, err := s.Services.User.Update(uint64(userID), updates)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return map[string]interface{}{"user_id": user.ID}, nil
}

// executeWebhook executes a webhook action
func (s *AutomationService) executeWebhook(action models.AutomationAction, input map[string]interface{}) (interface{}, error) {
	if s.Services == nil || s.Services.Webhook == nil {
		return nil, fmt.Errorf("webhook service not available")
	}

	url, _ := action.Config["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("webhook URL is required")
	}

	event := "automation.triggered"
	if eventName, ok := action.Config["event"].(string); ok {
		event = eventName
	}

	if err := s.Services.Webhook.Trigger(event, input); err != nil {
		return nil, fmt.Errorf("failed to trigger webhook: %w", err)
	}

	return map[string]interface{}{"status": "triggered", "url": url}, nil
}

// HandleEvent handles an event and triggers matching workflows
func (s *AutomationService) HandleEvent(eventType string, payload map[string]interface{}) error {
	// Get all enabled workflows
	enabled := true
	workflows, err := s.ListWorkflows(&enabled)
	if err != nil {
		return fmt.Errorf("failed to list workflows: %w", err)
	}

	// Check each workflow for matching triggers
	for _, workflow := range workflows {
		var trigger models.AutomationTrigger
		if err := json.Unmarshal([]byte(workflow.Trigger), &trigger); err != nil {
			s.logger.WithError(err).Errorf("Failed to parse trigger for workflow %d", workflow.ID)
			continue
		}

		// Check if trigger matches event
		if trigger.Type == "event" && trigger.Event == eventType {
			// Check conditions if any
			if trigger.Conditions != nil {
				if !s.matchesConditions(trigger.Conditions, payload) {
					continue
				}
			}

			// Trigger workflow
			if _, err := s.TriggerWorkflow(workflow.ID, payload); err != nil {
				s.logger.WithError(err).Errorf("Failed to trigger workflow %d", workflow.ID)
			} else {
				s.logger.Infof("Triggered workflow %d for event %s", workflow.ID, eventType)
			}
		}
	}

	return nil
}

// matchesConditions checks if payload matches trigger conditions
func (s *AutomationService) matchesConditions(conditions map[string]interface{}, payload map[string]interface{}) bool {
	for key, expectedValue := range conditions {
		actualValue, ok := payload[key]
		if !ok {
			return false
		}
		if actualValue != expectedValue {
			return false
		}
	}
	return true
}

// GetExecution gets an execution by ID
func (s *AutomationService) GetExecution(id uint64) (*models.AutomationExecution, error) {
	var execution models.AutomationExecution
	if err := s.db.Preload("Workflow").First(&execution, id).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

// ListExecutions lists executions for a workflow
func (s *AutomationService) ListExecutions(workflowID uint64, limit int) ([]models.AutomationExecution, error) {
	var executions []models.AutomationExecution
	query := s.db.Where("workflow_id = ?", workflowID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}
