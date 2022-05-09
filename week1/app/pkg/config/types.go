package config

// Credentials include Mix 3 oauth settings
type Credentials struct {
	AuthDisabled  bool   `json:"auth_disabled"`
	TokenURL      string `json:"token_url,omitempty"`
	ClientID      string `json:"client_id,omitempty"`
	ClientSecret  string `json:"client_secret,omitempty"`
	Scope         string `json:"scope,omitempty"`
	AuthTimeoutMS int    `json:"auth_timeout_ms,omitempty"`
}

type Consumer struct {
	Name   string                 `json:"name,omitempty" yaml:"name,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}

type Transformer struct {
	Type   string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}

type Storage struct {
	Type   string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}

type Cache struct {
	Type   string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}
type MessageQueue struct {
	Type   string                 `json:"type,omitempty" yaml:"type,omitempty"`
	Config map[string]interface{} `json:",omitempty" yaml:",omitempty"`
}
