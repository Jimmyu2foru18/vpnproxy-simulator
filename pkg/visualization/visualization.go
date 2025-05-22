package visualization

import (
	"fmt"
	"strings"
	"time"
)
type DataFlowEvent struct {
	Direction   string   
	Data        []byte  
	Timestamp   time.Time 
	Source      string  
	Destination string  
	Size        int    
}

type ColorScheme struct {
	Client string
	Proxy  string
	Target string
	Reset  string
}

func DefaultColorScheme() ColorScheme {
	return ColorScheme{
		Client: "\033[32m", // Green
		Proxy:  "\033[36m", // Cyan
		Target: "\033[35m", // Magenta
		Reset:  "\033[0m",  // Reset
	}
}

type Visualizer struct {
	Colors      ColorScheme
	ShowHex     bool
	ShowText    bool
	MaxHexBytes int
	MaxTextLen  int
	Enabled     bool
}

func NewVisualizer() *Visualizer {
	return &Visualizer{
		Colors:      DefaultColorScheme(),
		ShowHex:     true,
		ShowText:    true,
		MaxHexBytes: 16,
		MaxTextLen:  32,
		Enabled:     true,
	}
}

func (v *Visualizer) VisualizeEvent(event *DataFlowEvent) {
	if !v.Enabled {
		return
	}

	var sourceColor, destColor, arrow string
	switch {
	case strings.HasPrefix(event.Direction, "client"):
		sourceColor = v.Colors.Client
		destColor = v.Colors.Proxy
		arrow = "===>"
	case strings.HasPrefix(event.Direction, "proxy") && strings.Contains(event.Direction, "target"):
		sourceColor = v.Colors.Proxy
		destColor = v.Colors.Target
		arrow = "===>"
	case strings.HasPrefix(event.Direction, "target"):
		sourceColor = v.Colors.Target
		destColor = v.Colors.Proxy
		arrow = "<==="
	default: 
		sourceColor = v.Colors.Proxy
		destColor = v.Colors.Client
		arrow = "<==="
	}

	fmt.Printf("%s%s%s %s %s%s%s [%d bytes]\n",
		sourceColor, event.Source, v.Colors.Reset,
		arrow,
		destColor, event.Destination, v.Colors.Reset,
		event.Size)

	if v.ShowHex {
		hexView := formatHexView(event.Data, v.MaxHexBytes)
		fmt.Printf("HEX: %s\n", hexView)
	}

	if v.ShowText {
		textView := formatTextView(event.Data, v.MaxTextLen)
		fmt.Printf("TXT: %s\n\n", textView)
	}
}

func formatHexView(data []byte, maxBytes int) string {
	var builder strings.Builder

	for i, b := range data {
		if i >= maxBytes {
			builder.WriteString("...")
			break
		}
		fmt.Fprintf(&builder, "%02x ", b)
	}

	return builder.String()
}

func formatTextView(data []byte, maxLen int) string {
	var builder strings.Builder

	for i, b := range data {
		if i >= maxLen {
			builder.WriteString("...")
			break
		}

		if b >= 32 && b <= 126 { 
			builder.WriteByte(b)
		} else {
			builder.WriteString(".")
		}
	}

	return builder.String()
}