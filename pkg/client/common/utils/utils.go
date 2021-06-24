package utils

import (
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

var logger = logging.NewLogger("common/utils")

func TimeCost(mark string) func() {
	now := time.Now()
	return func() {
		s := time.Since(now)
		logger.Infof("%s, cost: %s", mark, s.String())
	}
}
