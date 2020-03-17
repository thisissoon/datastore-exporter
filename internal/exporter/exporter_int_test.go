package exporter

import (
	"context"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestExporter_Export(t *testing.T) {
	gcpID := os.Getenv("GCP_ID")
	gcsBucketName := os.Getenv("GCS_BUCKETNAME")
	if gcpID == "" && gcsBucketName == "" {
		t.Skipf("Skipping tests environment variables not set for GCP_ID, GCS_BUCKETNAME")
	}
	tests := map[string]struct {
		ctx     context.Context
		wantErr bool
	}{
		"create backup": {
			ctx:     context.Background(),
			wantErr: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log := zerolog.New(os.Stdout)
			e, err := NewExporter(tt.ctx, log, gcpID, gcsBucketName)
			if err != nil {
				t.Fatal(err)
			}
			err = e.Export(tt.ctx)
			assert.Equal(t, err != nil, tt.wantErr, "Exporter.Export() error = %v, wantErr %v", err, tt.wantErr)
		})
	}
}
