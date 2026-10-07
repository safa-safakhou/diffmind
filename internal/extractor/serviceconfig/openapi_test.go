package serviceconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSpringOpenAPIInputsAreLocalAndExplicit(t *testing.T) {
	for _, tc := range []struct {
		name, input, generator string
		want                   int
	}{{"local", "${project.basedir}/contract.yaml", "spring", 1}, {"remote", "https://example.com/contract.yaml", "spring", 0}, {"outside", "../contract.yaml", "spring", 0}, {"client", "contract.yaml", "typescript-axios", 0}, {"unresolved", "${custom.path}/contract.yaml", "spring", 0}} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			os.WriteFile(filepath.Join(root, "contract.yaml"), []byte("openapi: 3.0.3\n"), 0644)
			pom := `<project><build><plugins><plugin><groupId>org.openapitools</groupId><artifactId>openapi-generator-maven-plugin</artifactId><configuration><inputSpec>INPUT</inputSpec><generatorName>GENERATOR</generatorName><apiPackage>example.api</apiPackage><configOptions><interfaceOnly>true</interfaceOnly></configOptions></configuration><executions><execution/></executions></plugin></plugins></build></project>`
			pom = strings.ReplaceAll(strings.ReplaceAll(pom, "INPUT", tc.input), "GENERATOR", tc.generator)
			os.WriteFile(filepath.Join(root, "pom.xml"), []byte(pom), 0644)
			if got := SpringOpenAPISpecs(root); len(got) != tc.want {
				t.Fatalf("specs %+v", got)
			}
		})
	}
}
