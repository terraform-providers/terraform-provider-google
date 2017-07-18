// Generates an initial version of a schema for a new resource type.
//
// This script draws heavily from https://github.com/radeksimko/terraform-gen,
// but uses GCP's discovery API instead of the struct definition to generate
// the schemas.
//
// This is not meant to be a definitive source of truth for resource schemas,
// just a starting point. It has some notable deficiencies, such as:
// 	* No way to differentiate between fields that are/are not updateable.
// 	* Required/Optional/Computed are set based on keywords in the description.
//
// Usage example (from root dir):
//
//   go run ./scripts/schemagen.go -api pubsub -resource Subscription -version v1
//
// This will output a file in the directory from which the script is run named `gen_resource_[api]_[resource].go`.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/discovery/v1"
)

func main() {
	api := flag.String("api", "", "api to query")
	resource := flag.String("resource", "", "resource to generate")
	version := flag.String("version", "v1", "api version to query")
	flag.Parse()

	if *api == "" || *resource == "" {
		log.Fatal("usage: go run schemagen.go -api $API -resource $RESOURCE -version $VERSION")
	}

	// Discovery API doesn't need authentication
	client, err := google.DefaultClient(oauth2.NoContext, []string{}...)
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating client: %v", err))
	}

	discoveryService, err := discovery.New(client)
	if err != nil {
		log.Fatal(fmt.Errorf("Error creating service: %v", err))
	}

	resp, err := discoveryService.Apis.GetRest(*api, *version).Fields("schemas").Do()
	if err != nil {
		log.Fatal(fmt.Errorf("Error reading API: %v", err))
	}

	fileName := fmt.Sprintf("gen_resource_%s_%s.go", *api, underscore(*resource))
	f, err := os.Create(fileName)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fields := generateFields(resp.Schemas, *resource)

	err = googleTemplate.Execute(f, struct {
		TypeName string
		Fields   map[string]string
	}{
		// Capitalize the first letter of the api name, then concatenate the resource name onto it.
		// e.g. compute, instance -> ComputeInstance
		TypeName: strings.ToUpper((*api)[0:1]) + (*api)[1:] + *resource,
		Fields:   fields,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = exec.Command("gofmt", "-w", fileName).Run()
	if err != nil {
		log.Fatal(err)
	}
}

func generateFields(schemas map[string]discovery.JsonSchema, property string) map[string]string {
	fields := make(map[string]string, 0)

	for k, v := range schemas[property].Properties {
		content, err := generateField(schemas, k, v, false)
		if err != nil {
			log.Printf("ERROR: %s", err)
		} else {
			fields[underscore(k)] = content
		}
	}

	return fields
}

func generateField(schemas map[string]discovery.JsonSchema, field string, v discovery.JsonSchema, isNested bool) (string, error) {
	s := &schema.Schema{
		Description: v.Description,
	}
	if field != "" {
		setProperties(v, s)
	}

	// JSON field types: https://tools.ietf.org/html/draft-zyp-json-schema-03#section-5.1
	switch v.Type {
	case "integer":
		s.Type = schema.TypeInt
	case "number":
		s.Type = schema.TypeFloat
	case "string":
		s.Type = schema.TypeString
	case "boolean":
		s.Type = schema.TypeBool
	case "array":
		s.Type = schema.TypeList
		elem, err := generateField(schemas, "", *v.Items, true)
		if err != nil {
			return "", fmt.Errorf("Unable to generate Elem for %q: %s", field, err)
		}
		s.Elem = elem
	case "object":
		s.Type = schema.TypeMap
	case "":
		s.Type = schema.TypeList
		s.MaxItems = 1

		elem := "&schema.Resource{\nSchema: map[string]*schema.Schema{\n"
		m := generateFields(schemas, v.Ref)
		fieldNames := []string{}
		for k, _ := range m {
			fieldNames = append(fieldNames, k)
		}
		sort.Strings(fieldNames)
		for _, k := range fieldNames {
			elem += fmt.Sprintf("%q: %s,\n", k, m[k])
		}
		elem += "},\n}"

		if isNested {
			return elem, nil
		}
		s.Elem = elem
	default:
		return "", fmt.Errorf("Unable to process: %s %s", field, v.Type)
	}

	return schemaCode(s, isNested)
}

func setProperties(v discovery.JsonSchema, s *schema.Schema) {
	if v.ReadOnly || strings.HasPrefix(v.Description, "Output-only") {
		s.Computed = true
	} else {
		if v.Required || strings.HasPrefix(v.Description, "Required") {
			s.Required = true
		} else {
			s.Optional = true
		}
	}

	s.ForceNew = true
}

func schemaCode(s *schema.Schema, isNested bool) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	err := schemaTemplate.Execute(buf, struct {
		Schema   *schema.Schema
		IsNested bool
	}{
		Schema:   s,
		IsNested: isNested,
	})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// shamelessly borrowed from https://github.com/radeksimko/terraform-gen/blob/master/internal/util/util.go
func underscore(name string) string {
	var words []string
	var camelCase = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

	for _, submatch := range camelCase.FindAllStringSubmatch(name, -1) {
		if submatch[1] != "" {
			words = append(words, submatch[1])
		}
		if submatch[2] != "" {
			words = append(words, submatch[2])
		}
	}

	return strings.ToLower(strings.Join(words, "_"))
}

var schemaTemplate = template.Must(template.New("schema").Parse(`{{if .IsNested}}&schema.Schema{{end}}{{"{"}}{{if not .IsNested}}
{{end}}Type: schema.{{.Schema.Type}},{{if ne .Schema.Description ""}}
Description: {{printf "%q" .Schema.Description}},{{end}}{{if .Schema.Required}}
Required: {{.Schema.Required}},{{end}}{{if .Schema.Optional}}
Optional: {{.Schema.Optional}},{{end}}{{if .Schema.ForceNew}}
ForceNew: {{.Schema.ForceNew}},{{end}}{{if .Schema.Computed}}
Computed: {{.Schema.Computed}},{{end}}{{if gt .Schema.MaxItems 0}}
MaxItems: {{.Schema.MaxItems}},{{end}}{{if .Schema.Elem}}
Elem: {{.Schema.Elem}},{{end}}{{if not .IsNested}}
{{end}}{{"}"}}`))

var googleTemplate = template.Must(template.New("google").Parse(`package google

import(
	"github.com/hashicorp/terraform/helper/schema"
)

func resource{{.TypeName}}() *schema.Resource {
	return &schema.Resource{
		Create: resource{{.TypeName}}Create,
		Read:   resource{{.TypeName}}Read,
		Update: resource{{.TypeName}}Update,
		Delete: resource{{.TypeName}}Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{ {{range $name, $schema := .Fields}}
			"{{ $name }}": {{ $schema }},
{{end}}
		},
	}
}

func resource{{.TypeName}}Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
}

func resource{{.TypeName}}Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
}

func resource{{.TypeName}}Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
}

func resource{{.TypeName}}Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
}
`))
