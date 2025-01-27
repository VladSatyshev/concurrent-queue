package queues

const (
	RepositoryErr = iota
	RepositoryNotFoundErr
	UseCaseErr
	UseCaseNotFoundErr
)

type QueueErr struct {
	ErrType int
	Msg     string
}

func (q *QueueErr) Error() string {
	return q.Msg
}

func NewQueueErr(errType int, msg string) error {
	return &QueueErr{
		ErrType: errType,
		Msg:     msg,
	}
}
