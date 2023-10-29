package md

import (
	"testing"

	"github.com/woweh/forge-api-go-client"
)

func TestNewMdApi(t *testing.T) {
	tests := []struct {
		name   string
		region forge.Region
		want   string
	}{
		{
			name:   "Test US Region",
			region: forge.US,
			want:   usPath,
		},
		{
			name:   "Test EMEA Region",
			region: forge.EMEA,
			want:   euPath,
		},
		{
			name:   "Test EU Region",
			region: forge.EU,
			want:   euPath,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mdApi := NewMdApi(nil, tt.region)
				if mdApi.RelativePath() != tt.want {
					t.Errorf("got RelativePath = %s, want %s", mdApi.RelativePath(), tt.want)
				}
				if mdApi.Region() != tt.region {
					t.Errorf("got Region = %s, want %s", mdApi.Region(), tt.region)
				}
			},
		)
	}
}

func TestModelDerivativeAPI_SetRegion(t *testing.T) {
	tests := []struct {
		name          string
		initialRegion forge.Region
		initialPath   string
		newRegion     forge.Region
		newPath       string
	}{
		{
			name:          "Test US Region",
			initialRegion: forge.US,
			initialPath:   usPath,
			newRegion:     forge.EMEA,
			newPath:       euPath,
		},
		{
			name:          "Test EMEA Region",
			initialRegion: forge.EMEA,
			initialPath:   euPath,
			newRegion:     forge.US,
			newPath:       usPath,
		},
		{
			name:          "Test EU Region",
			initialRegion: forge.EU,
			initialPath:   euPath,
			newRegion:     forge.US,
			newPath:       usPath,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				mdApi := NewMdApi(nil, tt.initialRegion)
				if mdApi.RelativePath() != tt.initialPath {
					t.Errorf("got RelativePath = %s, initialPath %s", mdApi.RelativePath(), tt.initialPath)
				}
				if mdApi.Region() != tt.initialRegion {
					t.Errorf("got Region = %s, initialPath %s", mdApi.Region(), tt.initialRegion)
				}
				t.Log("Set new region: ", tt.newRegion)
				mdApi.SetRegion(tt.newRegion)
				if mdApi.RelativePath() != tt.newPath {
					t.Errorf("got RelativePath = %s, newPath %s", mdApi.RelativePath(), tt.newPath)
				}
				if mdApi.Region() != tt.newRegion {
					t.Errorf("got Region = %s, newPath %s", mdApi.Region(), tt.newRegion)
				}
			},
		)
	}
}
