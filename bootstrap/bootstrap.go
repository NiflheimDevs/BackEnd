package bootstrap

type Di struct {
	Env   *Env
	Const *Constants
}

func Get() *Di {
	di := &Di{}
	di.Env = NewEnvironment()
	di.Const = NewConstant()
	return di
}
