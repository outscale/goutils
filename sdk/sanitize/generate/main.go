package main

import (
	"cmp"
	"flag"
	"fmt"
	"go/ast"
	"os"
	"slices"
	"strings"

	"github.com/fatih/structtag"
	"github.com/samber/lo"
	"golang.org/x/tools/go/packages"
)

type HTTPField struct {
	Name        string
	Sensitivity string
}

var (
	pkgName = flag.String("package", "", "package name")
	service = flag.String("service", "", "service name (as found in signature")
	output  = flag.String("output", "", "output file name")
	suffix  = flag.String("suffix", "", "struct name suffix")
	typ     = flag.String("type", "", "Request or Response")
)

func main() {
	flag.Parse()
	fields, err := fetchFields(*pkgName, *suffix)
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred:", err)
		os.Exit(1)
	}
	fd, err := os.Create(*output)
	if err == nil {
		_, err = fd.WriteString(
			`//nolint
package sanitize

func init() {
	serviceFields[` + *typ + `]["` + *service + `"] = map[string]string{
` +
				lo.Reduce(fields, func(agg string, f HTTPField, _ int) string {
					return agg + fmt.Sprintf(`		%q: %s,
`, f.Name, sensitivity[f.Sensitivity])
				}, "") + `	}
}`,
		)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "An error occurred:", err)
		os.Exit(1)
	}
}

var sensitivity = map[string]string{
	"pii":       "PII",
	"sensitive": "Sensitive",
}

func fetchFields(pkgName, suffix string) ([]HTTPField, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedSyntax | packages.LoadFiles,
	}, pkgName)
	if err != nil {
		return nil, err
	}
	fields := fetchFieldsIn(pkgs[0], func(name string) bool {
		return strings.HasSuffix(name, suffix)
	})
	slices.SortFunc(fields, func(a, b HTTPField) int {
		return cmp.Compare(a.Name, b.Name)
	})
	fields = lo.Uniq(fields)
	return fields, nil
}

func fetchFieldsIn(pkg *packages.Package, match func(name string) bool) []HTTPField {
	var fields []HTTPField
	done := map[string]bool{}
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			// Looking for `type [name]` where the type name matches
			if typ, ok := n.(*ast.TypeSpec); ok && match(typ.Name.Name) {
				fmt.Println("***", typ.Name.Name)
				// Looking for struct, with included fields
				if struc, ok := typ.Type.(*ast.StructType); ok && struc.Fields != nil {
					for _, field := range struc.Fields.List {
						if len(field.Names) == 0 {
							continue
						}
						typ := field.Type
						if star, ok := typ.(*ast.StarExpr); ok {
							typ = star.X
						}
						ident, ok := typ.(*ast.Ident)
						if !ok {
							continue
						}
						// Adding fields of the field type
						if !done[ident.Name] {
							fields = append(fields, fetchFieldsIn(pkg, func(name string) bool { return name == ident.Name })...)
							done[ident.Name] = true
						}
						if field.Tag == nil {
							continue
						}
						// Parsing tag (without quotes)
						tags, err := structtag.Parse(field.Tag.Value[1 : len(field.Tag.Value)-1])
						if err != nil {
							continue
						}
						if log, err := tags.Get("log"); err == nil {
							name := field.Names[0].Name
							if json, err := tags.Get("json"); err == nil {
								name = json.Name
							}
							fmt.Println(" -", name, log.Name)
							fields = append(fields, HTTPField{
								Name:        name,
								Sensitivity: log.Name,
							})
						}
					}
				}
			}
			return true
		})
	}
	return fields
}
