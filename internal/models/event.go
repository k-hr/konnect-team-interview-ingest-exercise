package models

// CDCEvent represents a Change Data Capture event from the source
type CDCEvent struct {
	Before interface{}     `json:"before"`
	After  CDCAfterPayload `json:"after"`
	Op     string         `json:"op"`
	TsMs   int64         `json:"ts_ms"`
}

// CDCAfterPayload represents the structure of the After field
type CDCAfterPayload struct {
	Key   string     `json:"key"`
	Value ValueType  `json:"value"`
}

// ValueType represents the value structure
type ValueType struct {
	Type   int         `json:"type"`
	Object interface{} `json:"object"`
}

// Service represents a Kong service entity
type Service struct {
	ID             string   `json:"id"`
	Host           string   `json:"host"`
	Name           string   `json:"name"`
	Path           string   `json:"path"`
	Port           int      `json:"port"`
	Tags           []string `json:"tags"`
	Enabled        bool     `json:"enabled"`
	Retries        int      `json:"retries"`
	Protocol       string   `json:"protocol"`
	CreatedAt      int64    `json:"created_at"`
	UpdatedAt      int64    `json:"updated_at"`
	ReadTimeout    int      `json:"read_timeout"`
	WriteTimeout   int      `json:"write_timeout"`
	ConnectTimeout int      `json:"connect_timeout"`
}

// Node represents a Kong node entity
type Node struct {
	ID              string                 `json:"id"`
	Type            string                 `json:"type"`
	Labels          map[string]string      `json:"labels,omitempty"`
	Version         string                 `json:"version"`
	Hostname        string                 `json:"hostname"`
	LastPing        int64                  `json:"last_ping"`
	CreatedAt       int64                  `json:"created_at"`
	UpdatedAt       int64                  `json:"updated_at"`
	ConfigHash      string                 `json:"config_hash"`
	ProcessConf     map[string]interface{} `json:"process_conf"`
	ConnectionState map[string]interface{} `json:"connection_state"`
	DataPlaneCertID string                `json:"data_plane_cert_id"`
}

// Upstream represents a Kong upstream entity
type Upstream struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Slots            int                    `json:"slots"`
	HashOn           string                 `json:"hash_on"`
	Algorithm        string                 `json:"algorithm"`
	CreatedAt        int64                 `json:"created_at"`
	UpdatedAt        int64                 `json:"updated_at"`
	Healthchecks     map[string]interface{} `json:"healthchecks"`
	UseSrvName       bool                   `json:"use_srv_name"`
	HashFallback     string                 `json:"hash_fallback"`
	HashOnCookiePath string                `json:"hash_on_cookie_path"`
}
