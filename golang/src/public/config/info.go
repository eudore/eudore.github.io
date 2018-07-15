package config

type configinfo struct {
	Help		bool		`comment:"Show help"`
	Test		bool		`comment:"Show config "`
	Workdir		string		`comment:"Show Workdir"`
	Enable		[]string	`comment:"Enable modes"`
	Disable		[]string	`comment:"Disable modes"`
	Mode		map[string]interface{}
}

func (ci *configinfo) getmode() []string {
	// set Enable
	var d []string
	for _,ve := range ci.Enable {
		b := true
		for _,vd := range ci.Disable {
			if ve==vd {
				b = false
				break
			}
		}
		if b{
			d=append(d,ve)
		}
	}
	return d
}
	