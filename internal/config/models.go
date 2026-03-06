package config

type Configs struct {
	App      App      `json:"app"`
	Envelope Envelope `json:"envelope"`
}

type App struct {
	Server               Server        `json:"server"`
	ClientTimeoutSeconds int64         `json:"client_timeout_seconds"`
	DebugMode            bool          `json:"debug_mode"`
	DefaultParams        DefaultParams `json:"default_params"`
	ApiKey               string        `json:"api_key"`
	Processing           string        `json:"processing"`
}

type Server struct {
	Name    string            `json:"name"`
	Address string            `json:"address"`
	Host    string            `json:"host"`
	Port    string            `json:"port"`
	Token   string            `json:"token"`
	Extra   map[string]string `json:"extra,omitempty"`
}

type DefaultParams struct {
}

type Envelope struct {
	Soapenv     string `json:"soapenv"`
	Tw          string `json:"tw"`
	Tran        string `json:"tran"`
	Iss         string `json:"iss"`
	Sub         string `json:"sub"`
	Com         string `json:"com"`
	Xsi         string `json:"xsi"`
	Tran1       string `json:"tran1"`
	Con         string `json:"con"`
	Con1        string `json:"con1"`
	Con2        string `json:"con2"`
	GetRsTag    string `json:"getRsTag"`
	SetTag      string `json:"setTag"`
	TokensAdmin string `json:"tokensAdmin"`
	Tcmn        string `json:"tcmn"`
	NetworkRid  string `json:"networkRid"`
}
