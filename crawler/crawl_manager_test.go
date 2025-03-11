package crawler

import (
	"testing"
)

func TestCrawlManager_DeregisterGoroutine(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    int
	}{
		{
			name:   "check goroutine de-allocation",
			fields: fields{AvailableGoroutines: 0},
			want:   1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			if err := cm.DeregisterGoroutine(); (err != nil) != tt.wantErr {
				t.Errorf("DeregisterGoroutine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if cm.AvailableGoroutines != tt.want {
				t.Errorf("DeregisterGoroutine() error got = %v, want %v", cm.AvailableGoroutines, tt.want)
			}
		})
	}
}

func TestCrawlManager_GenerateReport(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	type args struct {
		printToConsole bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]any
	}{
		{
			name: "check generated report",
			fields: fields{
				runningAverage: 10000,
				Failures:       20,
				Success:        80,
			},
			want: map[string]any{
				"averageResponseTime": 10000,
				"totalProcessed":      100,
				"errorRate":           0.2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			if got := cm.GenerateReport(tt.args.printToConsole); len(got) != len(tt.want) {
				t.Errorf("GenerateReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawlManager_IsGoroutineAvailable(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "goroutines available",
			fields: fields{AvailableGoroutines: 10},
			want:   true,
		},
		{
			name:   "goroutines unavailable",
			fields: fields{AvailableGoroutines: 0},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			if got := cm.IsGoroutineAvailable(); got != tt.want {
				t.Errorf("IsGoroutineAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrawlManager_RecordFailure(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "check record failure",
			fields: fields{Failures: 10},
			want:   11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			cm.RecordFailure()
			if cm.Failures != tt.want {
				t.Errorf("RecordFailure() = %v, want %v", cm.Failures, tt.want)
			}
		})
	}
}

func TestCrawlManager_RecordSuccess(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name:   "check record success",
			fields: fields{Success: 10},
			want:   11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			cm.RecordSuccess()
			if cm.Success != tt.want {
				t.Errorf("RecordFailure() = %v, want %v", cm.Success, tt.want)
			}
		})
	}
}

func TestCrawlManager_RegisterGoroutine(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name:    "goroutines exhausted",
			fields:  fields{AvailableGoroutines: 0},
			wantErr: true,
		},
		{
			name:   "goroutine registered",
			fields: fields{AvailableGoroutines: 10},
			want:   9,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			if err := cm.RegisterGoroutine(); (err != nil) != tt.wantErr {
				t.Errorf("RegisterGoroutine() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && cm.AvailableGoroutines != tt.want {
				t.Errorf("RegisterGoroutine() error = %v, wantErr %v", cm.AvailableGoroutines, tt.want)
			}
		})
	}
}

func TestCrawlManager_ShutdownCrawls(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "check shutdown signal",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			cm.ShutdownCrawls()
			if cm.ReceivedShutdownSignal != tt.want {
				t.Errorf("RegisterGoroutine() error = %v, wantErr %v", cm.ReceivedShutdownSignal, tt.want)
			}
		})
	}
}

func TestCrawlManager_UpdateRunningAverage(t *testing.T) {
	type fields struct {
		AvailableGoroutines    int
		runningAverage         float64
		ReceivedShutdownSignal bool
		Failures               int
		Success                int
	}
	type args struct {
		responseTime int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name:   "check updated average value",
			fields: fields{runningAverage: 2500},
			args:   args{responseTime: 7500},
			want:   5000,
		},
		{
			name: "check average value for the first call",
			args: args{responseTime: 7500},
			want: 7500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &CrawlManager{
				AvailableGoroutines:    tt.fields.AvailableGoroutines,
				runningAverage:         tt.fields.runningAverage,
				ReceivedShutdownSignal: tt.fields.ReceivedShutdownSignal,
				Failures:               tt.fields.Failures,
				Success:                tt.fields.Success,
			}
			cm.UpdateRunningAverage(tt.args.responseTime)
			if cm.runningAverage != tt.want {
				t.Errorf("UpdateRunningAverage() error = %v, wantErr %v", cm.runningAverage, tt.want)
			}
		})
	}
}
