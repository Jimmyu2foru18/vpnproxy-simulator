package tunnel

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Stats struct {
	BytesIn   int64
	BytesOut  int64
	StartTime time.Time
	EndTime   time.Time
}

type Tunnel struct {
	ClientConn net.Conn
	TargetConn net.Conn
	Logger     *logrus.Logger
	Stats      Stats
	Done       chan struct{}
	mutex      sync.Mutex
}

func NewTunnel(clientConn, targetConn net.Conn, logger *logrus.Logger) *Tunnel {
	return &Tunnel{
		ClientConn: clientConn,
		TargetConn: targetConn,
		Logger:     logger,
		Stats: Stats{
			StartTime: time.Now(),
		},
		Done:  make(chan struct{}),
		mutex: sync.Mutex{},
	}
}

func (t *Tunnel) Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		n, err := t.pipe(t.ClientConn, t.TargetConn, "client -> target")
		t.mutex.Lock()
		t.Stats.BytesIn += n
		t.mutex.Unlock()
		if err != nil && err != io.EOF {
			t.Logger.Errorf("Error in client->target: %v", err)
		}
		t.ClientConn.Close()
	}

	go func() {
		defer wg.Done()
		n, err := t.pipe(t.TargetConn, t.ClientConn, "target -> client")
		t.mutex.Lock()
		t.Stats.BytesOut += n
		t.mutex.Unlock()
		if err != nil && err != io.EOF {
			t.Logger.Errorf("Error in target->client: %v", err)
		}
		t.TargetConn.Close()
	}

	go func() {
		wg.Wait()
		t.mutex.Lock()
		t.Stats.EndTime = time.Now()
		t.mutex.Unlock()
		close(t.Done)
		t.Logger.Infof("Tunnel closed. Bytes in: %d, Bytes out: %d, Duration: %v",
			t.Stats.BytesIn, t.Stats.BytesOut, t.Stats.EndTime.Sub(t.Stats.StartTime))
	}()
}

func (t *Tunnel) pipe(src, dst net.Conn, direction string) (int64, error) {
	buffer := make([]byte, 4096)
	var total int64

	for {
		n, err := src.Read(buffer)
		if n > 0 {
			t.Logger.Debugf("%s: Read %d bytes", direction, n)
			_, werr := dst.Write(buffer[:n])
			if werr != nil {
				return total, werr
			}
			total += int64(n)
		}

		if err != nil {
			return total, err
		}
	}
}

func (t *Tunnel) Close() {
	t.ClientConn.Close()
	t.TargetConn.Close()
}

func (t *Tunnel) GetStats() Stats {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.Stats
}