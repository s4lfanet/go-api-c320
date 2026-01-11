package repository

import (
	"context"
	"fmt"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/utils"
)

// OnuRepository handles ONU data retrieval operations
type OnuRepository struct {
	snmp *gosnmp.GoSNMP
	cfg  *config.Config
}

// NewOnuRepository creates a new OnuRepository instance
func NewOnuRepository(snmp *gosnmp.GoSNMP, cfg *config.Config) *OnuRepository {
	return &OnuRepository{
		snmp: snmp,
		cfg:  cfg,
	}
}

// CalculatePonIndex calculates PON index based on ZTE C320 V2.1 formula
func CalculatePonIndex(boardID, ponID int) int {
	// Formula: 268500992 + (board * 8192) + (pon * 256)
	return 268500992 + (boardID * 8192) + (ponID * 256)
}

// GetByBoardIDAndPonID retrieves all ONUs for a specific board and PON
func (r *OnuRepository) GetByBoardIDAndPonID(ctx context.Context, boardID, ponID int) ([]model.OnuSerialNumber, error) {
	log.Info().Int("board", boardID).Int("pon", ponID).Msg("Getting ONUs for board and PON")

	ponIndex := CalculatePonIndex(boardID, ponID)
	var onus []model.OnuSerialNumber

	// Walk ONU table to get all ONUs (OID suffix 5 = Device SN)
	baseOID := fmt.Sprintf("1.3.6.1.4.1.3902.1012.3.13.3.1.5.%d", ponIndex)
	err := r.snmp.Walk(baseOID, func(variable gosnmp.SnmpPDU) error {
		// Extract ONU ID from OID (last part after PON index)
		oidParts := utils.ParseOID(variable.Name)
		if len(oidParts) < 2 {
			return nil
		}

		onuID := oidParts[len(oidParts)-1]
		serialNumber := utils.ExtractSerialNumber(variable.Value)

		onu := model.OnuSerialNumber{
			Board:        boardID,
			PON:          ponID,
			ID:           onuID,
			SerialNumber: serialNumber,
		}
		onus = append(onus, onu)
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to walk ONU table")
		return nil, err
	}

	log.Info().Int("count", len(onus)).Msg("Successfully retrieved ONUs")
	return onus, nil
}
