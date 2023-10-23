package main

import (
	"reflect"
	"testing"
)

func TestCleanup(t *testing.T) {
	type args struct {
		retain bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Cleanup(tt.args.retain)
		})
	}
}

func TestFormatMarkdownFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FormatMarkdownFile(tt.args.filePath)
		})
	}
}

func TestFunction(t *testing.T) {
	type args struct {
		wdtkID string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TestFunction(tt.args.wdtkID)
		})
	}
}

func TestGenerateHeader(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateHeader(); got != tt.want {
				t.Errorf("GenerateHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateProblemReports(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GenerateProblemReports()
		})
	}
}

func TestGenerateReportHeader(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateReportHeader(tt.args.title); got != tt.want {
				t.Errorf("GenerateReportHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCSVDatasetFromMySociety(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetCSVDatasetFromMySociety()
		})
	}
}

func TestMakeDataset(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeDataset()
		})
	}
}

func TestMakeTableFromGeneratedDataset(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MakeTableFromGeneratedDataset()
		})
	}
}

func TestNewPoliceOrganisation(t *testing.T) {
	type args struct {
		wdtkID string
		emails map[string]string
	}
	tests := []struct {
		name string
		args args
		want *Authority
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthority(tt.args.wdtkID, tt.args.emails); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthority() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessMySocietyDataset(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ProcessMySocietyDataset()
		})
	}
}

func TestReadCSVFile(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    []map[string]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCSVFile(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadCSVFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRebuildDataset(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RebuildDataset()
		})
	}
}
