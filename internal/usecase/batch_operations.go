package usecase

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
)

// BatchOperationsUsecaseInterface defines business logic for batch ONU operations
type BatchOperationsUsecaseInterface interface {
	BatchRebootONUs(ctx context.Context, req *model.BatchONURebootRequest) (*model.BatchONURebootResponse, error)
	BatchBlockONUs(ctx context.Context, req *model.BatchONUBlockRequest) (*model.BatchONUBlockResponse, error)
	BatchUnblockONUs(ctx context.Context, req *model.BatchONUBlockRequest) (*model.BatchONUBlockResponse, error)
	BatchDeleteONUs(ctx context.Context, req *model.BatchONUDeleteRequest) (*model.BatchONUDeleteResponse, error)
	BatchUpdateDescriptions(ctx context.Context, req *model.BatchONUDescriptionRequest) (*model.BatchONUDescriptionResponse, error)
}

// BatchOperationsUsecase implements batch ONU operations business logic
type BatchOperationsUsecase struct {
	telnetSessionManager *repository.TelnetSessionManager
	onuMgmtUsecase       ONUManagementUsecaseInterface
	cfg                  *config.Config
}

// NewBatchOperationsUsecase creates a new batch operations usecase
func NewBatchOperationsUsecase(
	telnetSessionManager *repository.TelnetSessionManager,
	onuMgmtUsecase ONUManagementUsecaseInterface,
	cfg *config.Config,
) BatchOperationsUsecaseInterface {
	return &BatchOperationsUsecase{
		telnetSessionManager: telnetSessionManager,
		onuMgmtUsecase:       onuMgmtUsecase,
		cfg:                  cfg,
	}
}

// BatchRebootONUs reboots multiple ONUs in parallel
func (u *BatchOperationsUsecase) BatchRebootONUs(ctx context.Context, req *model.BatchONURebootRequest) (*model.BatchONURebootResponse, error) {
	startTime := time.Now()

	// Validate targets
	if err := u.validateBatchTargets(req.Targets); err != nil {
		return nil, err
	}

	log.Info().Int("target_count", len(req.Targets)).Msg("Starting batch ONU reboot operation")

	// Execute operations sequentially (Telnet session is single-threaded)
	results := make([]model.BatchOperationResult, 0, len(req.Targets))
	successCount := 0
	failureCount := 0

	for _, target := range req.Targets {
		rebootReq := &model.ONURebootRequest{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		resp, err := u.onuMgmtUsecase.RebootONU(ctx, rebootReq)

		result := model.BatchOperationResult{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		if err != nil {
			result.Success = false
			result.Message = "Reboot failed"
			result.Error = err.Error()
			failureCount++
			log.Warn().
				Str("pon_port", target.PONPort).
				Int("onu_id", target.ONUID).
				Err(err).
				Msg("Failed to reboot ONU in batch")
		} else {
			result.Success = resp.Success
			result.Message = resp.Message
			if resp.Success {
				successCount++
			} else {
				failureCount++
			}
		}

		results = append(results, result)
	}

	executionTime := time.Since(startTime).Milliseconds()

	response := &model.BatchONURebootResponse{
		TotalTargets:    len(req.Targets),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		ExecutionTimeMs: executionTime,
	}

	log.Info().
		Int("total", len(req.Targets)).
		Int("success", successCount).
		Int("failure", failureCount).
		Int64("execution_time_ms", executionTime).
		Msg("Batch ONU reboot operation completed")

	return response, nil
}

// BatchBlockONUs blocks multiple ONUs
func (u *BatchOperationsUsecase) BatchBlockONUs(ctx context.Context, req *model.BatchONUBlockRequest) (*model.BatchONUBlockResponse, error) {
	startTime := time.Now()

	// Validate targets
	if err := u.validateBatchTargets(req.Targets); err != nil {
		return nil, err
	}

	action := "block"
	if !req.Block {
		action = "unblock"
	}

	log.Info().
		Int("target_count", len(req.Targets)).
		Str("action", action).
		Msg("Starting batch ONU block/unblock operation")

	// Execute operations sequentially
	results := make([]model.BatchOperationResult, 0, len(req.Targets))
	successCount := 0
	failureCount := 0

	for _, target := range req.Targets {
		blockReq := &model.ONUBlockRequest{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
			Block:   req.Block,
		}

		var resp *model.ONUBlockResponse
		var err error

		if req.Block {
			resp, err = u.onuMgmtUsecase.BlockONU(ctx, blockReq)
		} else {
			resp, err = u.onuMgmtUsecase.UnblockONU(ctx, blockReq)
		}

		result := model.BatchOperationResult{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("%s failed", action)
			result.Error = err.Error()
			failureCount++
			log.Warn().
				Str("pon_port", target.PONPort).
				Int("onu_id", target.ONUID).
				Str("action", action).
				Err(err).
				Msg("Failed to block/unblock ONU in batch")
		} else {
			result.Success = resp.Success
			result.Message = resp.Message
			if resp.Success {
				successCount++
			} else {
				failureCount++
			}
		}

		results = append(results, result)
	}

	executionTime := time.Since(startTime).Milliseconds()

	response := &model.BatchONUBlockResponse{
		Blocked:         req.Block,
		TotalTargets:    len(req.Targets),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		ExecutionTimeMs: executionTime,
	}

	log.Info().
		Int("total", len(req.Targets)).
		Int("success", successCount).
		Int("failure", failureCount).
		Str("action", action).
		Int64("execution_time_ms", executionTime).
		Msg("Batch ONU block/unblock operation completed")

	return response, nil
}

// BatchUnblockONUs is a convenience method for unblocking ONUs
func (u *BatchOperationsUsecase) BatchUnblockONUs(ctx context.Context, req *model.BatchONUBlockRequest) (*model.BatchONUBlockResponse, error) {
	req.Block = false
	return u.BatchBlockONUs(ctx, req)
}

// BatchDeleteONUs deletes multiple ONUs
func (u *BatchOperationsUsecase) BatchDeleteONUs(ctx context.Context, req *model.BatchONUDeleteRequest) (*model.BatchONUDeleteResponse, error) {
	startTime := time.Now()

	// Validate targets
	if err := u.validateBatchTargets(req.Targets); err != nil {
		return nil, err
	}

	log.Info().Int("target_count", len(req.Targets)).Msg("Starting batch ONU delete operation")

	// Execute operations sequentially
	results := make([]model.BatchOperationResult, 0, len(req.Targets))
	successCount := 0
	failureCount := 0

	for _, target := range req.Targets {
		deleteReq := &model.ONUDeleteRequest{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		resp, err := u.onuMgmtUsecase.DeleteONU(ctx, deleteReq)

		result := model.BatchOperationResult{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		if err != nil {
			result.Success = false
			result.Message = "Delete failed"
			result.Error = err.Error()
			failureCount++
			log.Warn().
				Str("pon_port", target.PONPort).
				Int("onu_id", target.ONUID).
				Err(err).
				Msg("Failed to delete ONU in batch")
		} else {
			result.Success = resp.Success
			result.Message = resp.Message
			if resp.Success {
				successCount++
			} else {
				failureCount++
			}
		}

		results = append(results, result)
	}

	executionTime := time.Since(startTime).Milliseconds()

	response := &model.BatchONUDeleteResponse{
		TotalTargets:    len(req.Targets),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		ExecutionTimeMs: executionTime,
	}

	log.Info().
		Int("total", len(req.Targets)).
		Int("success", successCount).
		Int("failure", failureCount).
		Int64("execution_time_ms", executionTime).
		Msg("Batch ONU delete operation completed")

	return response, nil
}

// BatchUpdateDescriptions updates descriptions for multiple ONUs
func (u *BatchOperationsUsecase) BatchUpdateDescriptions(ctx context.Context, req *model.BatchONUDescriptionRequest) (*model.BatchONUDescriptionResponse, error) {
	startTime := time.Now()

	// Validate targets
	if err := u.validateBatchDescriptionTargets(req.Targets); err != nil {
		return nil, err
	}

	log.Info().Int("target_count", len(req.Targets)).Msg("Starting batch ONU description update operation")

	// Execute operations sequentially
	results := make([]model.BatchOperationResult, 0, len(req.Targets))
	successCount := 0
	failureCount := 0

	for _, target := range req.Targets {
		descReq := &model.ONUDescriptionRequest{
			PONPort:     target.PONPort,
			ONUID:       target.ONUID,
			Description: target.Description,
		}

		resp, err := u.onuMgmtUsecase.UpdateDescription(ctx, descReq)

		result := model.BatchOperationResult{
			PONPort: target.PONPort,
			ONUID:   target.ONUID,
		}

		if err != nil {
			result.Success = false
			result.Message = "Description update failed"
			result.Error = err.Error()
			failureCount++
			log.Warn().
				Str("pon_port", target.PONPort).
				Int("onu_id", target.ONUID).
				Err(err).
				Msg("Failed to update ONU description in batch")
		} else {
			result.Success = resp.Success
			result.Message = resp.Message
			if resp.Success {
				successCount++
			} else {
				failureCount++
			}
		}

		results = append(results, result)
	}

	executionTime := time.Since(startTime).Milliseconds()

	response := &model.BatchONUDescriptionResponse{
		TotalTargets:    len(req.Targets),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		Results:         results,
		ExecutionTimeMs: executionTime,
	}

	log.Info().
		Int("total", len(req.Targets)).
		Int("success", successCount).
		Int("failure", failureCount).
		Int64("execution_time_ms", executionTime).
		Msg("Batch ONU description update operation completed")

	return response, nil
}

// ============================================
// Validation Functions
// ============================================

// validateBatchTargets validates an array of ONU targets
func (u *BatchOperationsUsecase) validateBatchTargets(targets []model.ONUTarget) error {
	if len(targets) == 0 {
		return fmt.Errorf("at least one target is required")
	}

	if len(targets) > 50 {
		return fmt.Errorf("maximum 50 targets allowed per batch operation, got %d", len(targets))
	}

	// Validate each target
	ponPortRegex := regexp.MustCompile(`^(\d+)/(\d+)/(\d+)$`)
	seenTargets := make(map[string]bool)

	for i, target := range targets {
		// Validate PON port format
		if !ponPortRegex.MatchString(target.PONPort) {
			return fmt.Errorf("invalid PON port format at index %d: %s (expected format: rack/shelf/port, e.g., 1/1/1)", i, target.PONPort)
		}

		// Validate ONU ID range
		if target.ONUID < 1 || target.ONUID > 128 {
			return fmt.Errorf("invalid ONU ID at index %d: %d (must be between 1 and 128)", i, target.ONUID)
		}

		// Check for duplicates
		targetKey := fmt.Sprintf("%s:%d", target.PONPort, target.ONUID)
		if seenTargets[targetKey] {
			return fmt.Errorf("duplicate target at index %d: PON %s, ONU ID %d", i, target.PONPort, target.ONUID)
		}
		seenTargets[targetKey] = true
	}

	return nil
}

// validateBatchDescriptionTargets validates an array of ONU description targets
func (u *BatchOperationsUsecase) validateBatchDescriptionTargets(targets []model.ONUDescriptionTarget) error {
	if len(targets) == 0 {
		return fmt.Errorf("at least one target is required")
	}

	if len(targets) > 50 {
		return fmt.Errorf("maximum 50 targets allowed per batch operation, got %d", len(targets))
	}

	// Validate each target
	ponPortRegex := regexp.MustCompile(`^(\d+)/(\d+)/(\d+)$`)
	descriptionRegex := regexp.MustCompile(`^[a-zA-Z0-9\s\-_.]+$`)
	seenTargets := make(map[string]bool)

	for i, target := range targets {
		// Validate PON port format
		if !ponPortRegex.MatchString(target.PONPort) {
			return fmt.Errorf("invalid PON port format at index %d: %s (expected format: rack/shelf/port, e.g., 1/1/1)", i, target.PONPort)
		}

		// Validate ONU ID range
		if target.ONUID < 1 || target.ONUID > 128 {
			return fmt.Errorf("invalid ONU ID at index %d: %d (must be between 1 and 128)", i, target.ONUID)
		}

		// Validate description
		if target.Description == "" {
			return fmt.Errorf("description cannot be empty at index %d", i)
		}

		if len(target.Description) > 64 {
			return fmt.Errorf("description too long at index %d: %d characters (max 64)", i, len(target.Description))
		}

		if !descriptionRegex.MatchString(target.Description) {
			return fmt.Errorf("invalid description format at index %d: %s (allowed: alphanumeric, spaces, dashes, underscores, dots)", i, target.Description)
		}

		// Check for duplicates
		targetKey := fmt.Sprintf("%s:%d", target.PONPort, target.ONUID)
		if seenTargets[targetKey] {
			return fmt.Errorf("duplicate target at index %d: PON %s, ONU ID %d", i, target.PONPort, target.ONUID)
		}
		seenTargets[targetKey] = true
	}

	return nil
}
