package TG

type Plugin interface {
	DependsOn() []string
	Init(tg *TG)
	OnInstall()
	OnUninstall()
	Name() string
}
