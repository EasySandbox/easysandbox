package domains

type DomainIsRunningError struct {
	Domain string
	Msg    string
}

func (e *DomainIsRunningError) Error() string {
	return e.Msg
}


type BackingFileDoesNotExistError struct {
	DiskImagePath string
	Msg    string
}

func (e *BackingFileDoesNotExistError) Error() string {
	return e.Msg
}

