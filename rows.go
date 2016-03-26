// Package rows provides Json-schema based database rows to arbitrary struct
// scanner generation
package generows

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"go/format"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/schema"
)

type Generator struct{}

// Generate generates and writes the errors of the schema
func (g *Generator) Generate(req *common.Req, res *common.Res) error {
	context := req.Context
	s := req.Schema

	if context == nil || context.Config == nil {
		return nil
	}

	if !common.IsIn("rows", context.Config.Generators...) {
		return nil
	}

	temp := template.New("rowscanner.tmpl").Funcs(context.TemplateFuncs)
	if _, err := temp.Parse(RowScannerTemplate); err != nil {
		return err
	}

	outputs := make([]common.Output, 0)

	for _, def := range common.SortedObjectSchemas(s.Definitions) {
		data := struct {
			Schema *schema.Schema
		}{
			Schema: def,
		}

		var buf bytes.Buffer

		if err := temp.ExecuteTemplate(&buf, "rowscanner.tmpl", data); err != nil {
			return err
		}

		f, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}

		outputs = append(outputs, common.Output{
			Content: f,
			Path: fmt.Sprintf(
				"%s%s_rowscanner.go",
				context.Config.Target,
				strings.ToLower(def.Title),
			),
		})
	}

	res.Output = outputs
	return nil
}
