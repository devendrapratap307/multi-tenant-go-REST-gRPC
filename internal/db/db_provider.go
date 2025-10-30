package db

import (
	"context"
	"fmt"
	"go-multitenant/internal/model"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// cachedDB holds gorm.DB and lastUsed timestamp
type cachedDB struct {
	DB       *gorm.DB
	lastUsed time.Time
}

// DBProvider holds master DB and cache of client DBs
type DBProvider struct {
	master *gorm.DB
	cache  sync.Map // map[string]cachedDB
	ttl    time.Duration
}

// NewDBProvider constructs provider and starts cleanup goroutine
func NewDBProvider(master *gorm.DB) *DBProvider {
	p := &DBProvider{master: master, ttl: 10 * time.Minute}
	go p.cleanupLoop()
	return p
}

// GetClientDB returns a cached client DB or opens a new one
func (p *DBProvider) GetClientDB(ctx context.Context, clientID string) (*gorm.DB, error) {
	// quick cache read
	if v, ok := p.cache.Load(clientID); ok {
		entry := v.(cachedDB)
		entry.lastUsed = time.Now()
		p.cache.Store(clientID, entry)
		return entry.DB.WithContext(ctx), nil
	}
	// load client info from master DB
	var client model.Client
	if err := p.master.WithContext(ctx).First(&client, "id = ?",
		clientID).Error; err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d	sslmode=disable",
		client.DBHost, client.DBUser, client.DBPassword, client.DBName,
		client.DBPort)
	clientDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed open client db: %w", err)
	}
	// configure pool
	sqlDB, err := clientDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	p.cache.Store(clientID, cachedDB{DB: clientDB, lastUsed: time.Now()})
	return clientDB.WithContext(ctx), nil
}

// cleanupLoop periodically removes stale client DB connections
func (p *DBProvider) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		p.cache.Range(func(key, value interface{}) bool {
			k := key.(string)
			entry := value.(cachedDB)
			if now.Sub(entry.lastUsed) > p.ttl {
				if sqlDB, err := entry.DB.DB(); err == nil {
					_ = sqlDB.Close()
				}
				p.cache.Delete(k)
			}
			return true
		})
	}
}
