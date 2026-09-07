package system

// InputError describes malformed input or a missing request prerequisite.
// ValidationError describes a decoded character that violates domain constraints.
type InputError struct{ Err error }

func (e *InputError) Error() string { return e.Err.Error() }
func (e *InputError) Unwrap() error { return e.Err }
