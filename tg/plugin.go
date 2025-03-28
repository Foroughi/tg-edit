package TG

type Plugin interface {
	Init(api *Api)
	Name() string
}
