package definitions

import (
	"github.com/iddqdeika/pim"
)

type InfomodelProvider interface {
	GetByIdentifier(identifier string) (*pim.StructureGroup, error)
}

type WebService interface {
	Run() error
}

// предоставляет конфигурации по данному пути
// по умолчанию разделитель - точка
type Config interface {
	GetString(path string) (string, error)
	GetInt(path string) (int, error)
	GetArray(path string) ([]Config, error)
	Child(path string) Config
}

type JsonInfomodelDTO struct {
	Identifier string                 `json:"identifier"`
	Features   map[string]*FeatureDTO `json:"features"`
}

type XmlInfomodelDTO struct {
	Identifier string       `xml:"identifier"`
	Features   []FeatureDTO `xml:"features"`
}

type FeatureDTO struct {
	Name         string   `json:"name" xml:"name"`
	DataType     string   `json:"data_type" xml:"data_type"`
	PresetValues []string `json:"preset_values" xml:"preset_values"`
	Mandatory    bool     `json:"mandatory" xml:"mandatory"`
	Multivalued  bool     `json:"multivalued" xml:"multivalued"`
}
