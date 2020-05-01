package parcs

type Runner struct {
	*Engine
}

func NewRunner(engine *Engine) *Runner {
	return &Runner{engine}
}

func (r *Runner) Init() {}

func (r *Runner) Shutdown() {}
