package config

type ConsumerCfg struct {
	Name   string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}

type Storage struct {
	Type   string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}
