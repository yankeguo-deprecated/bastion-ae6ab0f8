package nova // import "github.com/novakit/nova"

// Env empty-safe environment string
type Env string

const (
	// Production production env value
	Production = Env("production")

	// Development development env value
	Development = Env("development")

	// Test test env value
	Test = Env("test")
)

// IsProduction is production
func (e Env) IsProduction() bool {
	return e == Production
}

// IsDevelopment is development, if empty returns true
func (e Env) IsDevelopment() bool {
	if len(e) == 0 {
		return true
	}
	return e == Development
}

// IsTest is test
func (e Env) IsTest() bool {
	return e == Test
}

// String convert to string
func (e Env) String() string {
	return string(e)
}
