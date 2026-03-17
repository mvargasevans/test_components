package main

// Contract is the top-level YAML schema.
type Contract struct {
	Name         string        `yaml:"name"`
	Consumer     string        `yaml:"consumer"`
	Provider     string        `yaml:"provider"`
	Interactions []Interaction `yaml:"interactions"`
}

// Interaction describes one protocol exchange.
type Interaction struct {
	Name     string  `yaml:"name"`
	Type     string  `yaml:"type"` // "request_reply" or "publish"
	Request  []Field `yaml:"request"`
	Response []Field `yaml:"response"`
	Fields   []Field `yaml:"fields"` // used for publish
}

// Field is a named, typed message field.
type Field struct {
	Name    string `yaml:"name"`
	Type    string `yaml:"type"`    // string | int | float64 | bool
	Example string `yaml:"example"` // optional example value for contract assertions
}
