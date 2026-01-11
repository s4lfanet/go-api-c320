package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/s4lfanet/go-api-c320/internal/handler"
	"github.com/s4lfanet/go-api-c320/internal/model"
	"github.com/s4lfanet/go-api-c320/internal/usecase"
)

// mockOnuUsecase for testing routes
type mockOnuUsecase struct{}

// mockPonUsecase for testing routes
type mockPonUsecase struct{}

// mockProfileUsecase for testing routes
type mockProfileUsecase struct{}

// mockCardUsecase for testing routes
type mockCardUsecase struct{}

func (m *mockOnuUsecase) GetByBoardIDAndPonID(ctx context.Context, boardID, ponID int) ([]model.ONUInfoPerBoard, error) {
	return nil, nil
}

func (m *mockOnuUsecase) GetByBoardIDPonIDAndOnuID(boardID, ponID, onuID int) (model.ONUCustomerInfo, error) {
	return model.ONUCustomerInfo{}, nil
}

func (m *mockOnuUsecase) GetEmptyOnuID(ctx context.Context, boardID, ponID int) ([]model.OnuID, error) {
	return nil, nil
}

func (m *mockOnuUsecase) GetOnuIDAndSerialNumber(boardID, ponID int) ([]model.OnuSerialNumber, error) {
	return nil, nil
}

func (m *mockOnuUsecase) UpdateEmptyOnuID(ctx context.Context, boardID, ponID int) error {
	return nil
}

func (m *mockOnuUsecase) GetByBoardIDAndPonIDWithPagination(boardID, ponID, page, pageSize int) ([]model.ONUInfoPerBoard, int) {
	return nil, 0
}

func (m *mockOnuUsecase) DeleteCache(ctx context.Context, boardID, ponID int) error {
	return nil
}

// Mock methods for PonUsecase
func (m *mockPonUsecase) GetPonPortInfo(ctx context.Context, boardID, ponID int) (*model.PonPortInfo, error) {
	return &model.PonPortInfo{}, nil
}

// Mock methods for ProfileUsecase
func (m *mockProfileUsecase) GetAllTrafficProfiles(ctx context.Context) ([]*model.TrafficProfile, error) {
	return nil, nil
}

func (m *mockProfileUsecase) GetTrafficProfile(ctx context.Context, profileID int) (*model.TrafficProfile, error) {
	return &model.TrafficProfile{}, nil
}

func (m *mockProfileUsecase) GetAllVlanProfiles(ctx context.Context) ([]*model.VlanProfile, error) {
	return nil, nil
}

// Mock methods for CardUsecase
func (m *mockCardUsecase) GetAllCards(ctx context.Context) ([]*model.CardInfo, error) {
	return nil, nil
}

func (m *mockCardUsecase) GetCard(ctx context.Context, rack, shelf, slot int) (*model.CardInfo, error) {
	return &model.CardInfo{}, nil
}

func TestRootHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	rootHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	expectedBody := "Hello, this is the root endpoint!"
	if rr.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, rr.Body.String())
	}
}

func TestLoadRoutes(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)

	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	if router == nil {
		t.Error("Expected non-nil router")
	}
}

func TestLoadRoutes_RootEndpoint(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)
	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK for root endpoint, got %d", rr.Code)
	}
}

func TestLoadRoutes_MiddlewareApplied(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)
	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// Check if RequestID middleware added the header
	requestID := rr.Header().Get("X-Request-ID")
	if requestID == "" {
		t.Error("Expected X-Request-ID header to be set by middleware")
	}

	// Check if security headers are set
	if rr.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Error("Expected X-Content-Type-Options header to be set")
	}
}

func TestLoadRoutes_CORSHeaders(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)
	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	req := httptest.NewRequest("OPTIONS", "/", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	// CORS middleware should handle OPTIONS requests
	if rr.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("Expected CORS headers to be set")
	}
}

func TestLoadRoutes_APIv1ProfileRoutes(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)

	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"GET all traffic profiles", "/api/v1/profiles/traffic/", http.StatusOK},
		{"GET specific traffic profile", "/api/v1/profiles/traffic/1", http.StatusOK},
		{"GET all VLAN profiles", "/api/v1/profiles/vlan/", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, rr.Code)
			}
		})
	}
}

func TestLoadRoutes_APIv1SystemRoutes(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)

	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	tests := []struct {
		name   string
		path   string
		status int
	}{
		{"GET all cards", "/api/v1/system/cards/", http.StatusOK},
		{"GET specific card", "/api/v1/system/cards/1/1/1", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			if rr.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, rr.Code)
			}
		})
	}
}

func TestLoadRoutes_APIv1RouteNotFound(t *testing.T) {
	onuUsecase := &mockOnuUsecase{}
	ponUsecase := &mockPonUsecase{}
	profileUsecase := &mockProfileUsecase{}
	cardUsecase := &mockCardUsecase{}

	onuHandler := handler.NewOnuHandler(onuUsecase)
	ponHandler := handler.NewPonHandler(ponUsecase)
	profileHandler := handler.NewProfileHandler(profileUsecase)
	cardHandler := handler.NewCardHandler(cardUsecase)
	provisionUsecase := usecase.NewProvisionUsecase(nil, nil)
	provisionHandler := handler.NewProvisionHandler(provisionUsecase)
	vlanUsecase := usecase.NewVLANUsecase(nil, nil)
	vlanHandler := handler.NewVLANHandler(vlanUsecase)
	trafficUsecase := usecase.NewTrafficUsecase(nil, nil)
	trafficHandler := handler.NewTrafficHandler(trafficUsecase)
	onuMgmtUsecase := usecase.NewONUManagementUsecase(nil, nil)
	onuMgmtHandler := handler.NewONUManagementHandler(onuMgmtUsecase)
	batchUsecase := usecase.NewBatchOperationsUsecase(nil, onuMgmtUsecase, nil)
	batchHandler := handler.NewBatchOperationsHandler(batchUsecase)
	configBackupUsecase := usecase.NewConfigBackupUsecase(nil, onuMgmtUsecase, vlanUsecase, trafficUsecase, provisionUsecase)
	configBackupHandler := handler.NewConfigBackupHandler(configBackupUsecase)
	router := loadRoutes(onuHandler, ponHandler, profileHandler, cardHandler, provisionHandler, vlanHandler, trafficHandler, onuMgmtHandler, batchHandler, configBackupHandler, nil)

	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent route, got %d", http.StatusNotFound, rr.Code)
	}
}
