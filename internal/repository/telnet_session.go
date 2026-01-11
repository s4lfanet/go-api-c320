package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/config"
	"github.com/Cepat-Kilat-Teknologi/go-snmp-olt-zte-c320/internal/model"
	"github.com/rs/zerolog/log"
)

// TelnetSessionPool manages a pool of telnet connections
type TelnetSessionPool struct {
	config    *config.TelnetConfig
	session   TelnetRepository
	mu        sync.Mutex
	lastUsed  time.Time
	inUse     bool
	closeChan chan struct{}
	closed    bool
}

// NewTelnetSessionPool creates a new telnet session pool
func NewTelnetSessionPool(cfg *config.TelnetConfig) *TelnetSessionPool {
	return &TelnetSessionPool{
		config:    cfg,
		session:   NewTelnetRepository(cfg),
		closeChan: make(chan struct{}),
		closed:    false,
	}
}

// GetSession acquires a session from the pool
func (p *TelnetSessionPool) GetSession(ctx context.Context) (TelnetRepository, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil, model.NewTelnetError(model.ErrCodeSessionBusy,
			"session pool is closed", false)
	}

	// Check if session is in use
	deadline := time.Now().Add(30 * time.Second)
	for p.inUse {
		// Release lock and wait
		p.mu.Unlock()

		select {
		case <-ctx.Done():
			p.mu.Lock()
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			// Check if deadline exceeded
			if time.Now().After(deadline) {
				p.mu.Lock()
				return nil, model.NewTelnetError(model.ErrCodeSessionBusy,
					"timeout waiting for available session", true)
			}
		}

		p.mu.Lock()
	}

	// Mark as in use
	p.inUse = true
	p.lastUsed = time.Now()

	// Ensure connection is established
	if !p.session.IsConnected() {
		log.Info().Msg("Session not connected, establishing connection")
		if err := p.session.Connect(); err != nil {
			p.inUse = false
			return nil, err
		}
	}

	// Check if connection is stale
	if time.Since(p.lastUsed) > p.config.MaxIdleTime {
		log.Info().Msg("Session idle too long, reconnecting")
		if err := p.session.Reconnect(); err != nil {
			p.inUse = false
			return nil, err
		}
	}

	return p.session, nil
}

// ReleaseSession releases the session back to the pool
func (p *TelnetSessionPool) ReleaseSession() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.inUse = false
	p.lastUsed = time.Now()

	log.Debug().Msg("Session released back to pool")
}

// Close closes all sessions in the pool
func (p *TelnetSessionPool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.closed {
		return nil
	}

	p.closed = true
	close(p.closeChan)

	if p.session != nil {
		return p.session.Close()
	}

	return nil
}

// StartIdleCleanup starts a goroutine to clean up idle connections
func (p *TelnetSessionPool) StartIdleCleanup() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.cleanupIdleSessions()
			case <-p.closeChan:
				return
			}
		}
	}()
}

// cleanupIdleSessions closes idle sessions
func (p *TelnetSessionPool) cleanupIdleSessions() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.inUse {
		return
	}

	if p.session == nil || !p.session.IsConnected() {
		return
	}

	// Close if idle for too long
	idleDuration := time.Since(p.lastUsed)
	if idleDuration > p.config.MaxIdleTime {
		log.Info().
			Dur("idle_duration", idleDuration).
			Msg("Closing idle telnet session")

		if err := p.session.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close idle session")
		}
	}
}

// GetStatus returns the current status of the pool
func (p *TelnetSessionPool) GetStatus() map[string]interface{} {
	p.mu.Lock()
	defer p.mu.Unlock()

	status := map[string]interface{}{
		"closed":    p.closed,
		"in_use":    p.inUse,
		"last_used": p.lastUsed.Format(time.RFC3339),
		"idle_time": time.Since(p.lastUsed).String(),
	}

	if p.session != nil {
		connInfo := p.session.GetConnectionInfo()
		status["connected"] = connInfo.Connected
		status["mode"] = connInfo.Mode
		status["uptime"] = connInfo.Uptime
	}

	return status
}

// TelnetSessionManager manages the global telnet session pool
type TelnetSessionManager struct {
	pool *TelnetSessionPool
	mu   sync.RWMutex
}

var (
	globalSessionManager *TelnetSessionManager
	sessionManagerOnce   sync.Once
)

// GetGlobalSessionManager returns the singleton session manager
func GetGlobalSessionManager(cfg *config.TelnetConfig) *TelnetSessionManager {
	sessionManagerOnce.Do(func() {
		pool := NewTelnetSessionPool(cfg)
		pool.StartIdleCleanup()

		globalSessionManager = &TelnetSessionManager{
			pool: pool,
		}

		log.Info().Msg("Global telnet session manager initialized")
	})

	return globalSessionManager
}

// ExecuteCommand executes a command using the session pool
func (m *TelnetSessionManager) ExecuteCommand(ctx context.Context, command string) (*model.TelnetResponse, error) {
	session, err := m.pool.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	defer m.pool.ReleaseSession()

	return session.Execute(ctx, command)
}

// ExecuteCommands executes multiple commands using the session pool
func (m *TelnetSessionManager) ExecuteCommands(ctx context.Context, commands []string) (*model.TelnetBatchResponse, error) {
	session, err := m.pool.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	defer m.pool.ReleaseSession()

	return session.ExecuteMulti(ctx, commands)
}

// ExecuteInConfigMode executes commands in configuration mode
func (m *TelnetSessionManager) ExecuteInConfigMode(ctx context.Context, commands []string) (*model.TelnetBatchResponse, error) {
	session, err := m.pool.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	defer m.pool.ReleaseSession()

	// Enter config mode
	if err := session.EnterEnableMode(); err != nil {
		return nil, err
	}

	if err := session.EnterConfigMode(); err != nil {
		return nil, err
	}

	// Execute commands
	result, execErr := session.ExecuteMulti(ctx, commands)

	// Always try to exit config mode
	if exitErr := session.ExitConfigMode(); exitErr != nil {
		log.Error().Err(exitErr).Msg("Failed to exit config mode")
	}

	return result, execErr
}

// SaveConfiguration saves the OLT configuration
func (m *TelnetSessionManager) SaveConfiguration(ctx context.Context) error {
	session, err := m.pool.GetSession(ctx)
	if err != nil {
		return err
	}
	defer m.pool.ReleaseSession()

	return session.SaveConfig()
}

// GetConnectionStatus returns the connection status
func (m *TelnetSessionManager) GetConnectionStatus() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.pool.GetStatus()
}

// Close closes the session manager
func (m *TelnetSessionManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pool != nil {
		return m.pool.Close()
	}

	return nil
}

// WithRetry executes a function with retry logic
func (m *TelnetSessionManager) WithRetry(ctx context.Context, fn func() error) error {
	cfg := m.pool.config
	var lastErr error

	for attempt := 0; attempt <= cfg.RetryCount; attempt++ {
		if attempt > 0 {
			log.Info().
				Int("attempt", attempt).
				Int("max_attempts", cfg.RetryCount).
				Msg("Retrying operation")

			// Wait before retry
			select {
			case <-time.After(cfg.RetryDelay):
			case <-ctx.Done():
				return ctx.Err()
			}

			// Try to reconnect
			session, err := m.pool.GetSession(ctx)
			if err == nil {
				session.Reconnect()
				m.pool.ReleaseSession()
			}
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Check if error is recoverable
		if telnetErr, ok := lastErr.(*model.TelnetError); ok {
			if !telnetErr.Recoverable {
				return lastErr
			}
		}

		log.Warn().
			Err(lastErr).
			Int("attempt", attempt+1).
			Msg("Operation failed, will retry")
	}

	return fmt.Errorf("operation failed after %d attempts: %w", cfg.RetryCount+1, lastErr)
}
