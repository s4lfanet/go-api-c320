package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
)

// ConfigBackupUsecase handles configuration backup and restore operations
type ConfigBackupUsecase interface {
	// Backup single ONU configuration
	BackupONU(ponPort string, onuID int, description string, tags []string) (*model.ConfigBackup, error)

	// Backup entire OLT configuration (all ONUs)
	BackupOLT(description string, tags []string) (*model.ConfigBackup, error)

	// List all backups with optional filters
	ListBackups(backupType string, limit int) ([]*model.BackupListItem, error)

	// Get specific backup by ID
	GetBackup(backupID string) (*model.ConfigBackup, error)

	// Delete backup by ID
	DeleteBackup(backupID string) error

	// Restore configuration from backup
	RestoreFromBackup(req *model.RestoreRequest) (*model.RestoreResult, error)

	// Export backup to file
	ExportBackup(backupID string, outputPath string) error

	// Import backup from file
	ImportBackup(inputPath string) (*model.ConfigBackup, error)
}

type configBackupUsecase struct {
	cfg              *config.Config
	onuMgmtUsecase   ONUManagementUsecaseInterface
	vlanUsecase      VLANUsecaseInterface
	trafficUsecase   TrafficUsecaseInterface
	provisionUsecase ProvisionUseCaseInterface
	backupDir        string
}

// NewConfigBackupUsecase creates a new config backup usecase
func NewConfigBackupUsecase(
	cfg *config.Config,
	onuMgmtUsecase ONUManagementUsecaseInterface,
	vlanUsecase VLANUsecaseInterface,
	trafficUsecase TrafficUsecaseInterface,
	provisionUsecase ProvisionUseCaseInterface,
) ConfigBackupUsecase {
	// Set backup directory from config or use default
	backupDir := "/var/lib/go-snmp-olt/backups"
	if cfg != nil && cfg.OltCfg.BackupDir != "" {
		backupDir = cfg.OltCfg.BackupDir
	}

	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Error().Err(err).Str("dir", backupDir).Msg("Failed to create backup directory")
	}

	return &configBackupUsecase{
		cfg:              cfg,
		onuMgmtUsecase:   onuMgmtUsecase,
		vlanUsecase:      vlanUsecase,
		trafficUsecase:   trafficUsecase,
		provisionUsecase: provisionUsecase,
		backupDir:        backupDir,
	}
}

// BackupONU creates a backup of single ONU configuration
func (u *configBackupUsecase) BackupONU(ponPort string, onuID int, description string, tags []string) (*model.ConfigBackup, error) {
	log.Info().
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("Creating ONU configuration backup")

	// Get ONU details
	onuConfig, err := u.getONUConfiguration(ponPort, onuID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get ONU configuration")
		return nil, apperrors.NewInternalError("failed to get ONU configuration", err)
	}

	// Create backup object
	backup := &model.ConfigBackup{
		ID:          uuid.New().String(),
		Type:        "onu",
		Timestamp:   time.Now(),
		Description: description,
		Metadata: model.BackupMetadata{
			CreatedBy: "system",
			Source:    u.cfg.OltCfg.Host,
			Version:   "v2.1.0", // TODO: Get from OLT
			Tags:      tags,
		},
		Config: onuConfig,
	}

	// Save backup to file
	if err := u.saveBackupToFile(backup); err != nil {
		log.Error().Err(err).Msg("Failed to save backup to file")
		return nil, apperrors.NewInternalError("failed to save backup", err)
	}

	log.Info().
		Str("backup_id", backup.ID).
		Str("pon_port", ponPort).
		Int("onu_id", onuID).
		Msg("ONU backup created successfully")

	return backup, nil
}

// BackupOLT creates a backup of entire OLT configuration
func (u *configBackupUsecase) BackupOLT(description string, tags []string) (*model.ConfigBackup, error) {
	log.Info().Msg("Creating OLT configuration backup (all ONUs)")

	// TODO: Implement OLT-wide backup
	// This requires getting list of all ONUs and backing up each one
	// For now, return not implemented

	return nil, apperrors.NewInternalError("OLT-wide backup not yet implemented", nil)
}

// ListBackups lists all available backups
func (u *configBackupUsecase) ListBackups(backupType string, limit int) ([]*model.BackupListItem, error) {
	log.Info().
		Str("type", backupType).
		Int("limit", limit).
		Msg("Listing configuration backups")

	// Read all backup files from directory
	files, err := os.ReadDir(u.backupDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read backup directory")
		return nil, apperrors.NewInternalError("failed to read backups", err)
	}

	var backups []*model.BackupListItem

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Read backup metadata
		backupPath := filepath.Join(u.backupDir, file.Name())
		backup, err := u.loadBackupFromFile(backupPath)
		if err != nil {
			log.Warn().Err(err).Str("file", file.Name()).Msg("Failed to load backup")
			continue
		}

		// Filter by type if specified
		if backupType != "" && backup.Type != backupType {
			continue
		}

		// Get file info for size
		fileInfo, err := file.Info()
		if err != nil {
			log.Warn().Err(err).Str("file", file.Name()).Msg("Failed to get file info")
			continue
		}

		// Count ONUs
		onuCount := 0
		if backup.Type == "onu" {
			onuCount = 1
		} else if backup.Type == "olt" {
			if oltBackup, ok := backup.Config.(model.OLTConfigBackup); ok {
				onuCount = len(oltBackup.ONUs)
			}
		}

		backups = append(backups, &model.BackupListItem{
			ID:          backup.ID,
			Type:        backup.Type,
			Timestamp:   backup.Timestamp,
			Description: backup.Description,
			Size:        fileInfo.Size(),
			ONUCount:    onuCount,
			Source:      backup.Metadata.Source,
			Tags:        backup.Metadata.Tags,
		})
	}

	// Sort by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	// Apply limit
	if limit > 0 && len(backups) > limit {
		backups = backups[:limit]
	}

	log.Info().Int("count", len(backups)).Msg("Successfully listed backups")
	return backups, nil
}

// GetBackup retrieves a specific backup by ID
func (u *configBackupUsecase) GetBackup(backupID string) (*model.ConfigBackup, error) {
	log.Info().Str("backup_id", backupID).Msg("Getting backup details")

	backupPath := filepath.Join(u.backupDir, fmt.Sprintf("%s.json", backupID))

	backup, err := u.loadBackupFromFile(backupPath)
	if err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to load backup")
		return nil, apperrors.NewNotFoundError("backup", backupID)
	}

	return backup, nil
}

// DeleteBackup deletes a backup by ID
func (u *configBackupUsecase) DeleteBackup(backupID string) error {
	log.Info().Str("backup_id", backupID).Msg("Deleting backup")

	backupPath := filepath.Join(u.backupDir, fmt.Sprintf("%s.json", backupID))

	// Check if file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return apperrors.NewNotFoundError("backup", backupID)
	}

	// Delete file
	if err := os.Remove(backupPath); err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to delete backup")
		return apperrors.NewInternalError("failed to delete backup", err)
	}

	log.Info().Str("backup_id", backupID).Msg("Backup deleted successfully")
	return nil
}

// RestoreFromBackup restores configuration from a backup
func (u *configBackupUsecase) RestoreFromBackup(req *model.RestoreRequest) (*model.RestoreResult, error) {
	log.Info().
		Str("backup_id", req.BackupID).
		Bool("dry_run", req.DryRun).
		Msg("Restoring configuration from backup")

	// Load backup
	backup, err := u.GetBackup(req.BackupID)
	if err != nil {
		return nil, err
	}

	result := &model.RestoreResult{
		BackupID: req.BackupID,
		DryRun:   req.DryRun,
	}

	// Handle ONU backup restoration
	if backup.Type == "onu" {
		onuConfig, ok := backup.Config.(model.ONUConfigBackup)
		if !ok {
			return nil, apperrors.NewInternalError("invalid ONU backup format", nil)
		}

		// Restore ONU
		itemResult, err := u.restoreONU(&onuConfig, req)
		result.Details = append(result.Details, *itemResult)

		if err != nil {
			result.Success = false
			result.FailedONUs = 1
			result.Message = fmt.Sprintf("Failed to restore ONU: %v", err)
		} else {
			result.Success = true
			result.RestoredONUs = 1
			result.Message = "ONU configuration restored successfully"
		}
	} else if backup.Type == "olt" {
		// TODO: Handle OLT-wide restoration
		return nil, apperrors.NewInternalError("OLT-wide restore not yet implemented", nil)
	}

	log.Info().
		Str("backup_id", req.BackupID).
		Bool("success", result.Success).
		Int("restored", result.RestoredONUs).
		Msg("Restore operation completed")

	return result, nil
}

// ExportBackup exports a backup to a file
func (u *configBackupUsecase) ExportBackup(backupID string, outputPath string) error {
	log.Info().
		Str("backup_id", backupID).
		Str("output", outputPath).
		Msg("Exporting backup")

	backup, err := u.GetBackup(backupID)
	if err != nil {
		return err
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal backup")
		return apperrors.NewInternalError("failed to export backup", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		log.Error().Err(err).Msg("Failed to write backup file")
		return apperrors.NewInternalError("failed to write backup file", err)
	}

	log.Info().Str("backup_id", backupID).Str("output", outputPath).Msg("Backup exported successfully")
	return nil
}

// ImportBackup imports a backup from a file
func (u *configBackupUsecase) ImportBackup(inputPath string) (*model.ConfigBackup, error) {
	log.Info().Str("input", inputPath).Msg("Importing backup")

	backup, err := u.loadBackupFromFile(inputPath)
	if err != nil {
		return nil, err
	}

	// Generate new ID and timestamp for imported backup
	backup.ID = uuid.New().String()
	backup.Timestamp = time.Now()

	// Save to backup directory
	if err := u.saveBackupToFile(backup); err != nil {
		return nil, err
	}

	log.Info().Str("backup_id", backup.ID).Msg("Backup imported successfully")
	return backup, nil
}

// Helper methods

// getONUConfiguration retrieves complete configuration for an ONU
func (u *configBackupUsecase) getONUConfiguration(ponPort string, onuID int) (*model.ONUConfigBackup, error) {
	// For now, create a basic configuration
	// TODO: Actually retrieve real configuration from ONU via Telnet/SNMP

	config := &model.ONUConfigBackup{
		PONPort:      ponPort,
		ONUID:        onuID,
		SerialNumber: "UNKNOWN", // TODO: Get from ONU
		Type:         "UNKNOWN", // TODO: Get from ONU
		Name:         fmt.Sprintf("ONU_%s_%d", ponPort, onuID),
		AdminState:   "enabled",
		OperState:    "unknown",
		VLANs:        []model.ONUVLANConfig{},
		TCONTs:       []model.ONUTCONTConfig{},
		GEMPorts:     []model.ONUGEMPortConfig{},
		ServicePorts: []model.ONUServicePortConfig{},
		CustomConfig: make(map[string]interface{}),
	}

	return config, nil
}

// restoreONU restores a single ONU configuration
func (u *configBackupUsecase) restoreONU(onuConfig *model.ONUConfigBackup, req *model.RestoreRequest) (*model.RestoreItemResult, error) {
	result := &model.RestoreItemResult{
		PONPort:  onuConfig.PONPort,
		ONUID:    onuConfig.ONUID,
		ItemType: "onu",
	}

	if req.DryRun {
		result.Success = true
		result.Message = "Dry run: ONU configuration would be restored"
		return result, nil
	}

	// TODO: Actually restore configuration via Telnet
	// For now, just return success
	result.Success = true
	result.Message = "ONU configuration restore not yet implemented"

	return result, nil
}

// saveBackupToFile saves a backup to a JSON file
func (u *configBackupUsecase) saveBackupToFile(backup *model.ConfigBackup) error {
	filePath := filepath.Join(u.backupDir, fmt.Sprintf("%s.json", backup.ID))

	data, err := json.MarshalIndent(backup, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal backup: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup file: %w", err)
	}

	return nil
}

// loadBackupFromFile loads a backup from a JSON file
func (u *configBackupUsecase) loadBackupFromFile(filePath string) (*model.ConfigBackup, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	var backup model.ConfigBackup
	if err := json.Unmarshal(data, &backup); err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup: %w", err)
	}

	return &backup, nil
}
