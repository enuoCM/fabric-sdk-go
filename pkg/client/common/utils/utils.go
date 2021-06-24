package utils

import (
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
)

var logger = logging.NewLogger("common/utils")

func TimeCost(mark ...string) func() {
	now := time.Now()
	return func() {
		s := time.Since(now)
		m := ""
		if len(mark) == 1 {
			m = mark[0]
		} else {
			for _, v := range mark {
				m += v
			}
		}
		logger.Infof("%s, cost: %s", m, s.String())
	}
}
