package metrics

import (
	"sync"
	"time"
)

type ConnectionStats struct {
	ID            string
	ClientAddr    string
	TargetAddr    string
	BytesIn       int64
	BytesOut      int64
	StartTime     time.Time
	EndTime       time.Time
	Duration      time.Duration
	Active        bool
	ErrorOccurred bool
	ErrorMessage  string
}

type ProxyMetrics struct {
	StartTime          time.Time
	TotalConnections   int64
	ActiveConnections  int64
	TotalBytesIn       int64
	TotalBytesOut      int64
	FailedConnections  int64
	ConnectionsHistory []*ConnectionStats
	mutex              sync.RWMutex
	maxHistory         int
}

func NewProxyMetrics(maxHistory int) *ProxyMetrics {
	return &ProxyMetrics{
		StartTime:          time.Now(),
		ConnectionsHistory: make([]*ConnectionStats, 0, maxHistory),
		maxHistory:         maxHistory,
	}
}

func (pm *ProxyMetrics) NewConnection(clientAddr, targetAddr string) *ConnectionStats {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.TotalConnections++
	pm.ActiveConnections++

	stats := &ConnectionStats{
		ID:         generateID(),
		ClientAddr: clientAddr,
		TargetAddr: targetAddr,
		StartTime:  time.Now(),
		Active:     true,
	}

	pm.ConnectionsHistory = append(pm.ConnectionsHistory, stats)
	if len(pm.ConnectionsHistory) > pm.maxHistory {
		pm.ConnectionsHistory = pm.ConnectionsHistory[1:]
	}

	return stats
}

func (pm *ProxyMetrics) CloseConnection(stats *ConnectionStats, bytesIn, bytesOut int64, err error) {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	stats.Active = false
	stats.EndTime = time.Now()
	stats.Duration = stats.EndTime.Sub(stats.StartTime)
	stats.BytesIn = bytesIn
	stats.BytesOut = bytesOut

	if err != nil {
		stats.ErrorOccurred = true
		stats.ErrorMessage = err.Error()
		pm.FailedConnections++
	}

	pm.ActiveConnections--
	pm.TotalBytesIn += bytesIn
	pm.TotalBytesOut += bytesOut
}

func (pm *ProxyMetrics) GetSummary() map[string]interface{} {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	return map[string]interface{}{
		"uptime":             time.Since(pm.StartTime).String(),
		"total_connections":  pm.TotalConnections,
		"active_connections": pm.ActiveConnections,
		"failed_connections": pm.FailedConnections,
		"total_bytes_in":     pm.TotalBytesIn,
		"total_bytes_out":    pm.TotalBytesOut,
	}
}

func (pm *ProxyMetrics) GetRecentConnections(limit int) []*ConnectionStats {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	if limit <= 0 || limit > len(pm.ConnectionsHistory) {
		limit = len(pm.ConnectionsHistory)
	}

	result := make([]*ConnectionStats, limit)
	start := len(pm.ConnectionsHistory) - limit
	copy(result, pm.ConnectionsHistory[start:])

	return result
}

func generateID() string {
	return time.Now().Format("20060102-150405.000")
}