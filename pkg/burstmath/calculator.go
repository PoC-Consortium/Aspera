package burstmath

var RequestHandler *DeadlineRequestHandler

func Init(workerCount int, timeoutSeconds int64) {
	RequestHandler = NewDeadlineRequestHandler(workerCount, timeoutSeconds)
}
