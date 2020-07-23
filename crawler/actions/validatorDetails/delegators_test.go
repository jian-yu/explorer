package validatorDetails

import (
	"conf"
	"logger"
	"testing"
)

func TestGetDelegations(t *testing.T) {
	config := conf.NewConfig()   // CONFIG
	log := logger.NewLogger()    // LOG
	GetDelegations(config,log)
}
func TestGetDelegations2(t *testing.T) {
	config := conf.NewConfig()   // CONFIG
	log := logger.NewLogger()    // LOG
	GetDelegatorNums(config,log)
}
