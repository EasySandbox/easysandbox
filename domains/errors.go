package domains

type DomainIsRunningError struct {
	Domain string
	Msg    string
}

func (e *DomainIsRunningError) Error() string {
	return e.Msg
}
