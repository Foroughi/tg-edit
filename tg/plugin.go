package TG

type Plugin interface {
	Init(tg *TG)
	OnInstall()
	OnUninstall()
	Name() string
}
