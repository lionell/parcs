package parcs

type Runner struct {
	*Engine
}

func DefaultRunner() *Runner {
	return NewRunner(NewEnvEngine())
}

func NewRunner(engine *Engine) *Runner {
	return &Runner{engine}
}

func (r *Runner) Init() {}

func (r *Runner) Shutdown() {}
