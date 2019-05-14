package model

type Program struct {
	ID, Name, Path       string
	Versions, Exceptions []string
}

func (program Program) Latest() (version string) {
	amount := len(program.Versions)
	latest := program.Versions[amount-1]

	return latest
}
