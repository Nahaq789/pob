package phase

type PostMovePhaseHandler struct{}

func NewPostMovePhaseHandler() *PostMovePhaseHandler {
	return &PostMovePhaseHandler{}
}

func (post *PostMovePhaseHandler) Handle(ctx PostMoveContext) Result {
	return Result{}
}
