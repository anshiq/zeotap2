package models

type ConnectionRequest struct {
	SourceType       string           `json:"sourceType"`
	ClickHouseConfig ClickHouseConfig `json:"clickHouseConfig"`
	FlatFileConfig   FlatFileConfig   `json:"flatFileConfig"`
}

type ClickHouseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	JWTToken string `json:"jwtToken"`
	Secure   bool   `json:"secure"`
}

type FlatFileConfig struct {
	FilePath  string `json:"filePath"`
	Delimiter string `json:"delimiter"`
}

type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type PreviewRequest struct {
	SourceType       string           `json:"sourceType"`
	Table            string           `json:"table"`
	Columns          []string         `json:"columns"`
	ClickHouseConfig ClickHouseConfig `json:"clickHouseConfig"`
	FlatFileConfig   FlatFileConfig   `json:"flatFileConfig"`
}

type IngestRequest struct {
	Direction        string           `json:"direction"`
	Table            string           `json:"table"`
	Columns          []string         `json:"columns"`
	ClickHouseConfig ClickHouseConfig `json:"clickHouseConfig"`
	FlatFileConfig   FlatFileConfig   `json:"flatFileConfig"`
}
