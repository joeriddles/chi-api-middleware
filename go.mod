module chidemo

go 1.21.3

// See https://stackoverflow.com/a/72312461 for replacing with a fork
replace github.com/go-chi/chi/v5 => github.com/joeriddles/chi/v5 v5.0.0-20231102191906-6f38e5802ad4

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
