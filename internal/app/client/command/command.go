package command

type ClientCommand interface {
	Name() string
	Desc() string
	Usage() string
	Execute() error
}
