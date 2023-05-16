package check

var (
	_weekplanChecker weekplanChecker
)

func init() {
	go _weekplanChecker.routine()
}

// weekplanChecker 用于检查周计划
type weekplanChecker struct {
}

func (c *weekplanChecker) routine() {
	for {

	}
}
