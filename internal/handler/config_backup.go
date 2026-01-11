package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
	"github.com/s4lfanet/go-api-c320/internal/utils"
)

// ConfigBackupHandler handles configuration backup and restore requests
type ConfigBackupHandler struct {
	configBackupUsecase usecase.ConfigBackupUsecase
}

// NewConfigBackupHandler creates a new config backup handler
func NewConfigBackupHandler(configBackupUsecase usecase.ConfigBackupUsecase) *ConfigBackupHandler {
	return &ConfigBackupHandler{
		configBackupUsecase: configBackupUsecase,
	}
}

// BackupONU godoc
// @Summary Backup single ONU configuration
// @Description Creates a complete backup of ONU configuration including VLAN, T-CONT, GEM ports, and service ports
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param pon path string true "PON port (format: rack/shelf/port, e.g., 1/1/1)"
// @Param onuId path int true "ONU ID (1-128)"
// @Param request body model.BackupCreateRequest false "Backup options"
// @Success 200 {object} utils.WebResponse{data=model.ConfigBackup}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/onu/{pon}/{onuId} [post]
func (h *ConfigBackupHandler) BackupONU(w http.ResponseWriter, r *http.Request) {
	ponPort := chi.URLParam(r, "pon")
	onuIDStr := chi.URLParam(r, "onuId")

	onuID, err := strconv.Atoi(onuIDStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid ONU ID format")
		utils.HandleError(w, err)
		return
	}

	// Parse request body (optional)
	var req model.BackupCreateRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// If no body or invalid, use defaults
			req.Type = "onu"
			req.PONPort = ponPort
			req.ONUID = onuID
		}
	}

	// Create backup
	backup, err := h.configBackupUsecase.BackupONU(ponPort, onuID, req.Description, req.Tags)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create ONU backup")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   backup,
	})
}

// BackupOLT godoc
// @Summary Backup entire OLT configuration
// @Description Creates a complete backup of all ONU configurations on the OLT
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param request body model.BackupCreateRequest false "Backup options"
// @Success 200 {object} utils.WebResponse{data=model.ConfigBackup}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/olt [post]
func (h *ConfigBackupHandler) BackupOLT(w http.ResponseWriter, r *http.Request) {
	var req model.BackupCreateRequest
	if r.Body != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// Use defaults if no body
			req.Type = "olt"
		}
	}

	backup, err := h.configBackupUsecase.BackupOLT(req.Description, req.Tags)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create OLT backup")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   backup,
	})
}

// ListBackups godoc
// @Summary List all configuration backups
// @Description Retrieves a list of all available configuration backups with metadata
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param type query string false "Filter by backup type: onu or olt"
// @Param limit query int false "Maximum number of backups to return (default: all)"
// @Success 200 {object} utils.WebResponse{data=[]model.BackupListItem}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backups [get]
func (h *ConfigBackupHandler) ListBackups(w http.ResponseWriter, r *http.Request) {
	backupType := r.URL.Query().Get("type")
	limitStr := r.URL.Query().Get("limit")

	limit := 0
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	backups, err := h.configBackupUsecase.ListBackups(backupType, limit)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list backups")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   backups,
	})
}

// GetBackup godoc
// @Summary Get specific backup details
// @Description Retrieves complete details of a specific configuration backup
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param backupId path string true "Backup ID (UUID)"
// @Success 200 {object} utils.WebResponse{data=model.ConfigBackup}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/{backupId} [get]
func (h *ConfigBackupHandler) GetBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backupId")

	backup, err := h.configBackupUsecase.GetBackup(backupID)
	if err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to get backup")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   backup,
	})
}

// DeleteBackup godoc
// @Summary Delete a configuration backup
// @Description Permanently deletes a configuration backup file
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param backupId path string true "Backup ID (UUID)"
// @Success 200 {object} utils.WebResponse{data=string}
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/{backupId} [delete]
func (h *ConfigBackupHandler) DeleteBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backupId")

	if err := h.configBackupUsecase.DeleteBackup(backupID); err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to delete backup")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Backup deleted successfully",
	})
}

// RestoreFromBackup godoc
// @Summary Restore configuration from backup
// @Description Restores ONU or OLT configuration from a previously created backup
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param backupId path string true "Backup ID (UUID)"
// @Param request body model.RestoreRequest true "Restore options"
// @Success 200 {object} utils.WebResponse{data=model.RestoreResult}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/restore/{backupId} [post]
func (h *ConfigBackupHandler) RestoreFromBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backupId")

	var req model.RestoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request body")
		utils.HandleError(w, err)
		return
	}

	// Set backup ID from path parameter
	req.BackupID = backupID

	result, err := h.configBackupUsecase.RestoreFromBackup(&req)
	if err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to restore from backup")
		utils.HandleError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	})
}

// ExportBackup godoc
// @Summary Export backup to downloadable file
// @Description Exports a configuration backup as a downloadable JSON file
// @Tags Config Backup
// @Accept json
// @Produce json
// @Param backupId path string true "Backup ID (UUID)"
// @Success 200 {file} application/json
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/{backupId}/export [get]
func (h *ConfigBackupHandler) ExportBackup(w http.ResponseWriter, r *http.Request) {
	backupID := chi.URLParam(r, "backupId")

	backup, err := h.configBackupUsecase.GetBackup(backupID)
	if err != nil {
		log.Error().Err(err).Str("backup_id", backupID).Msg("Failed to export backup")
		utils.HandleError(w, err)
		return
	}

	// Set headers for file download
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename="+backupID+".json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(backup)
}

// ImportBackup godoc
// @Summary Import backup from file
// @Description Imports a configuration backup from an uploaded JSON file
// @Tags Config Backup
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Backup JSON file"
// @Success 200 {object} utils.WebResponse{data=model.ConfigBackup}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/config/backup/import [post]
func (h *ConfigBackupHandler) ImportBackup(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Error().Err(err).Msg("Failed to parse multipart form")
		utils.HandleError(w, err)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Error().Err(err).Msg("No file uploaded")
		utils.HandleError(w, err)
		return
	}
	defer file.Close()

	// Read backup from file
	var backup model.ConfigBackup
	if err := json.NewDecoder(file).Decode(&backup); err != nil {
		log.Error().Err(err).Msg("Failed to decode backup file")
		utils.HandleError(w, err)
		return
	}

	// Import backup (save to backup directory with new ID)
	// For now, just return the uploaded backup
	// TODO: Implement proper import with file saving

	log.Info().Str("filename", header.Filename).Msg("Backup file uploaded")

	utils.SendJSONResponse(w, http.StatusOK, utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   backup,
	})
}
