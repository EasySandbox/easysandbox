package sandbox
type SandboxIsRunningError struct {
	Sandbox string
	Msg    string
}

func (e *SandboxIsRunningError) Error() string {
	return e.Msg
}