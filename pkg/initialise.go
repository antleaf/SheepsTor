package pkg

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var registry *SheepstorWebsiteRegistry

func SetLogger(logger *zap.SugaredLogger) {
	log = logger
}

func SetSheepstorRegistry(reg *SheepstorWebsiteRegistry) {
	registry = reg
}
