package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/gosnmp/gosnmp"
	"github.com/rs/zerolog/log"
	"github.com/s4lfanet/go-api-c320/config"
	apperrors "github.com/s4lfanet/go-api-c320/internal/errors"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/repository"
	"golang.org/x/sync/singleflight"
)

// CardUseCaseInterface defines the interface for card/slot operations
type CardUseCaseInterface interface {
	GetAllCards(ctx context.Context) ([]*model.CardInfo, error)
	GetCard(ctx context.Context, rack, shelf, slot int) (*model.CardInfo, error)
}

// cardUsecase implements card usecase
type cardUsecase struct {
	snmpRepository  repository.SnmpRepositoryInterface
	redisRepository repository.OnuRedisRepositoryInterface
	cfg             *config.Config
	sg              singleflight.Group
}

// NewCardUsecase creates a new card usecase instance
func NewCardUsecase(
	snmpRepository repository.SnmpRepositoryInterface,
	redisRepository repository.OnuRedisRepositoryInterface,
	cfg *config.Config,
) CardUseCaseInterface {
	return &cardUsecase{
		snmpRepository:  snmpRepository,
		redisRepository: redisRepository,
		cfg:             cfg,
		sg:              singleflight.Group{},
	}
}

// GetAllCards retrieves all card/slot information
func (u *cardUsecase) GetAllCards(ctx context.Context) ([]*model.CardInfo, error) {
	log.Info().Msg("Getting all card information")

	// Map to store cards indexed by rack.shelf.slot
	cards := make(map[string]*model.CardInfo)

	// Walk the card info OID
	// OID: 1015.2.1.1.3.1.{col}.{rack}.{shelf}.{slot}
	baseOID := "1.3.6.1.4.1.3902.1015.2.1.1.3.1"

	err := u.snmpRepository.Walk(baseOID, func(pdu gosnmp.SnmpPDU) error {
		// Extract rack, shelf, slot from OID
		oidParts := strings.Split(pdu.Name, ".")

		// OID format: ...3902.1015.2.1.1.3.1.{col}.{rack}.{shelf}.{slot}
		// Base OID parts count
		baseOIDParts := strings.Split(baseOID, ".")
		if len(oidParts) < len(baseOIDParts)+4 {
			return nil
		}

		// Extract column, rack, shelf, slot
		column := 0
		rack := 0
		shelf := 0
		slot := 0

		_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)], "%d", &column)
		_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)+1], "%d", &rack)
		_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)+2], "%d", &shelf)
		_, _ = fmt.Sscanf(oidParts[len(baseOIDParts)+3], "%d", &slot)

		// Create unique key for this card
		cardKey := fmt.Sprintf("%d.%d.%d", rack, shelf, slot)

		// Initialize card if not exists
		if _, exists := cards[cardKey]; !exists {
			cards[cardKey] = &model.CardInfo{
				Rack:  rack,
				Shelf: shelf,
				Slot:  slot,
			}
		}

		card := cards[cardKey]

		// Parse value based on column
		switch column {
		case 2: // Card Type - appears to be a numeric type code
			if cardTypeNum, ok := pdu.Value.(int); ok {
				card.CardType = fmt.Sprintf("type_%d", cardTypeNum)
			} else {
				switch v := pdu.Value.(type) {
				case string:
					card.CardType = v
				case []byte:
					card.CardType = string(v)
				default:
					card.CardType = fmt.Sprintf("%v", v)
				}
			}
		case 3: // Also appears to be type-related
			// Skip or use as secondary type info
		case 4: // Serial Number (STRING)
			switch v := pdu.Value.(type) {
			case string:
				card.SerialNumber = v
			case []byte:
				card.SerialNumber = string(v)
			default:
				card.SerialNumber = fmt.Sprintf("%v", v)
			}
		case 5: // Hardware Version (INTEGER)
			if hwVer, ok := pdu.Value.(int); ok {
				card.HardwareVer = fmt.Sprintf("v%d", hwVer)
			}
		case 6: // Software Version (INTEGER)
			if swVer, ok := pdu.Value.(int); ok {
				card.SoftwareVer = fmt.Sprintf("v%d", swVer)
			}
		case 7: // Status or other info (INTEGER)
			if status, ok := pdu.Value.(int); ok {
				switch status {
				case 0:
					card.Status = "inactive"
				case 3:
					card.Status = "active"
				case 16:
					card.Status = "online"
				default:
					card.Status = fmt.Sprintf("status_%d", status)
				}
			}
		case 8: // Description
			if desc, ok := pdu.Value.(int); ok {
				card.Description = fmt.Sprintf("desc_%d", desc)
			}
		}

		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to walk card information")
		return nil, apperrors.NewSNMPError("failed to get card information", err)
	}

	// Convert map to slice
	result := make([]*model.CardInfo, 0, len(cards))
	for _, card := range cards {
		result = append(result, card)
	}

	log.Info().Int("count", len(result)).Msg("Successfully retrieved card information")
	return result, nil
}

// GetCard retrieves specific card information
func (u *cardUsecase) GetCard(ctx context.Context, rack, shelf, slot int) (*model.CardInfo, error) {
	log.Info().Int("rack", rack).Int("shelf", shelf).Int("slot", slot).Msg("Getting card information")

	card := &model.CardInfo{
		Rack:  rack,
		Shelf: shelf,
		Slot:  slot,
	}

	// Get card type (column 2)
	cardTypeOID := fmt.Sprintf("1.3.6.1.4.1.3902.1015.2.1.1.3.1.2.%d.%d.%d", rack, shelf, slot)
	if result, err := u.snmpRepository.Get([]string{cardTypeOID}); err == nil && len(result.Variables) > 0 {
		if cardTypeNum, ok := result.Variables[0].Value.(int); ok {
			card.CardType = fmt.Sprintf("type_%d", cardTypeNum)
		} else {
			switch v := result.Variables[0].Value.(type) {
			case string:
				card.CardType = v
			case []byte:
				card.CardType = string(v)
			default:
				card.CardType = fmt.Sprintf("%v", v)
			}
		}
	}

	// Get serial number (column 4)
	serialOID := fmt.Sprintf("1.3.6.1.4.1.3902.1015.2.1.1.3.1.4.%d.%d.%d", rack, shelf, slot)
	if result, err := u.snmpRepository.Get([]string{serialOID}); err == nil && len(result.Variables) > 0 {
		switch v := result.Variables[0].Value.(type) {
		case string:
			card.SerialNumber = v
		case []byte:
			card.SerialNumber = string(v)
		default:
			card.SerialNumber = fmt.Sprintf("%v", v)
		}
	}

	// Get hardware version (column 5)
	hwVerOID := fmt.Sprintf("1.3.6.1.4.1.3902.1015.2.1.1.3.1.5.%d.%d.%d", rack, shelf, slot)
	if result, err := u.snmpRepository.Get([]string{hwVerOID}); err == nil && len(result.Variables) > 0 {
		if hwVer, ok := result.Variables[0].Value.(int); ok {
			card.HardwareVer = fmt.Sprintf("v%d", hwVer)
		}
	}

	// Get software version (column 6)
	swVerOID := fmt.Sprintf("1.3.6.1.4.1.3902.1015.2.1.1.3.1.6.%d.%d.%d", rack, shelf, slot)
	if result, err := u.snmpRepository.Get([]string{swVerOID}); err == nil && len(result.Variables) > 0 {
		if swVer, ok := result.Variables[0].Value.(int); ok {
			card.SoftwareVer = fmt.Sprintf("v%d", swVer)
		}
	}

	// Get status/description (column 7)
	statusOID := fmt.Sprintf("1.3.6.1.4.1.3902.1015.2.1.1.3.1.7.%d.%d.%d", rack, shelf, slot)
	if result, err := u.snmpRepository.Get([]string{statusOID}); err == nil && len(result.Variables) > 0 {
		if status, ok := result.Variables[0].Value.(int); ok {
			switch status {
			case 0:
				card.Status = "inactive"
			case 3:
				card.Status = "active"
			case 16:
				card.Status = "online"
			default:
				card.Status = fmt.Sprintf("status_%d", status)
			}
		}
	}

	log.Info().
		Int("rack", rack).
		Int("shelf", shelf).
		Int("slot", slot).
		Str("card_type", card.CardType).
		Str("status", card.Status).
		Msg("Successfully retrieved card information")

	return card, nil
}

// getCardTypeName converts card type code to readable name
func (u *cardUsecase) getCardTypeName(cardType int) string {
	switch cardType {
	case 1:
		return "CTRL"
	case 2:
		return "GPON"
	case 3:
		return "EPON"
	case 4:
		return "GE"
	case 5:
		return "10GE"
	case 6:
		return "XGE"
	default:
		return fmt.Sprintf("type_%d", cardType)
	}
}
