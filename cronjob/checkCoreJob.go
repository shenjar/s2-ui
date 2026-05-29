package cronjob

import (
	"github.com/shen-sky6/s2-ui/service"
)

type CheckCoreJob struct {
	service.ConfigService
}

func NewCheckCoreJob() *CheckCoreJob {
	return &CheckCoreJob{}
}

func (s *CheckCoreJob) Run() {
	s.ConfigService.StartCore()
}
