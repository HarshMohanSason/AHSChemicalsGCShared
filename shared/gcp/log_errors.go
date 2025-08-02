package gcp

import (
	"context"
	"log"
	
	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
	"github.com/HarshMohanSason/AHSChemicalsGCShared/shared"
)

var (
	LogClient *logging.Client
)

func InitLogger(ctx context.Context) {
	shared.InitGCPOnce.Do(func() {
		projectID, err := metadata.ProjectIDWithContext(ctx)
		if err != nil {
			log.Fatalf("Failed to get project ID: %v", err)
		}

		LogClient, err = logging.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create logger: %v", err)
		}
	})
}

func CloseLogger() {
	if LogClient != nil {
		_ = LogClient.Close()
	}
}

func LogDebug(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Debug})
}

func LogError(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Error})
}

func LogInfo(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Info})
}

func LogWarning(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Warning})
}

func LogCritical(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Critical})
}

func LogNotice(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Notice})
}

func LogEmergency(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Emergency})
}

func LogAlert(cloudFnName string, message string) {
	if LogClient == nil {
		return
	}
	LogClient.Logger(cloudFnName).Log(logging.Entry{Payload: message, Severity: logging.Alert})
}