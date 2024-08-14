package customerrors

type ModuleNameAlreadyExists struct{}

func (e ModuleNameAlreadyExists) Error() string {
	return "Module name already exists!"
}

type TypeAssersionFailed struct{}

func (e TypeAssersionFailed) Error() string {
	return "Type assertion failed!"
}
