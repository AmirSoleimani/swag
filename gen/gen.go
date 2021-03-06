package gen

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/AmirSoleimani/swag"
	"github.com/ghodss/yaml"
)

// Gen presents a generate tool for swag.
type Gen struct {
}

// New creates a new Gen.
func New() *Gen {
	return &Gen{}
}

// Build builds swagger json file  for gived searchDir and mainAPIFile. Returns json
func (g *Gen) Build(searchDir, mainAPIFile, swaggerConfDir, propNamingStrategy string) (string, error) {
	log.Println("Generate swagger docs....")
	p := swag.New()
	p.PropNamingStrategy = propNamingStrategy
	p.ParseAPI(searchDir, mainAPIFile)
	swagger := p.GetSwagger()

	b, _ := json.MarshalIndent(swagger, "", "    ")

	os.MkdirAll(path.Join(searchDir, "docs"), os.ModePerm)
	docs, _ := os.Create(path.Join(searchDir, "docs", "docs.go"))
	defer docs.Close()

	os.Mkdir(swaggerConfDir, os.ModePerm)
	swaggerJSON, _ := os.Create(path.Join(swaggerConfDir, "swagger.json"))
	defer swaggerJSON.Close()
	swaggerJSON.Write(b)

	swaggerYAML, _ := os.Create(path.Join(swaggerConfDir, "swagger.yaml"))
	defer swaggerYAML.Close()
	y, err := yaml.JSONToYAML(b)
	if err != nil {
		log.Fatalf("can't swagger json covert to yaml err: %s", err)
	}
	swaggerYAML.Write(y)

	packageTemplate.Execute(docs, struct {
		Timestamp time.Time
		Doc       string
	}{
		Timestamp: time.Now(),
		Doc:       "`" + string(b) + "`",
	})

	log.Printf("create docs.go at  %+v", docs.Name())
	return string(b), nil
}

var packageTemplate = template.Must(template.New("").Parse(`// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// {{ .Timestamp }}

package docs

import (
	"github.com/AmirSoleimani/swag"
)

var doc = {{.Doc}}

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
`))
