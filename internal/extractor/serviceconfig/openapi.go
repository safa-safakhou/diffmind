package serviceconfig

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/mohammad-safakhou/diffmind/internal/extractor/sourcefilter"
)

type SpringOpenAPISpec struct{ Path, APIPackage string }

// SpringOpenAPISpecs reads explicit local Maven generator inputs. It never
// executes Maven, follows remote references, or guesses inherited configuration.
func SpringOpenAPISpecs(root string) []SpringOpenAPISpec {
	cfg, err := Load(root)
	if err != nil {
		return nil
	}
	policy, err := sourcefilter.NewPolicy(cfg.Paths.Include, cfg.Paths.Exclude)
	if err != nil {
		return nil
	}
	body, err := os.ReadFile(filepath.Join(root, "pom.xml"))
	if err != nil || len(body) > 256<<10 {
		return nil
	}
	type config struct {
		InputSpec     string `xml:"inputSpec"`
		Generator     string `xml:"generatorName"`
		APIPackage    string `xml:"apiPackage"`
		GenerateAPIs  string `xml:"generateApis"`
		InterfaceOnly string `xml:"configOptions>interfaceOnly"`
	}
	var pom struct {
		Plugins []struct {
			Group      string `xml:"groupId"`
			Artifact   string `xml:"artifactId"`
			Config     config `xml:"configuration"`
			Executions []struct {
				Config config `xml:"configuration"`
			} `xml:"executions>execution"`
		} `xml:"build>plugins>plugin"`
	}
	if xml.Unmarshal(body, &pom) != nil {
		return nil
	}
	var out []SpringOpenAPISpec
	seen := map[SpringOpenAPISpec]bool{}
	for _, plugin := range pom.Plugins {
		if plugin.Group != "org.openapitools" || plugin.Artifact != "openapi-generator-maven-plugin" {
			continue
		}
		configs := []config{plugin.Config}
		for _, execution := range plugin.Executions {
			c := execution.Config
			if c.InputSpec == "" {
				c.InputSpec = plugin.Config.InputSpec
			}
			if c.Generator == "" {
				c.Generator = plugin.Config.Generator
			}
			if c.APIPackage == "" {
				c.APIPackage = plugin.Config.APIPackage
			}
			if c.GenerateAPIs == "" {
				c.GenerateAPIs = plugin.Config.GenerateAPIs
			}
			if c.InterfaceOnly == "" {
				c.InterfaceOnly = plugin.Config.InterfaceOnly
			}
			configs = append(configs, c)
		}
		for _, c := range configs {
			if strings.TrimSpace(c.Generator) != "spring" || strings.TrimSpace(c.InterfaceOnly) != "true" || strings.TrimSpace(c.GenerateAPIs) == "false" || strings.TrimSpace(c.APIPackage) == "" {
				continue
			}
			input := strings.TrimSpace(c.InputSpec)
			input = strings.TrimPrefix(input, "${project.basedir}/")
			input = strings.TrimPrefix(input, "${basedir}/")
			if input == "" || filepath.IsAbs(input) || strings.Contains(input, "${") || strings.Contains(input, "://") {
				continue
			}
			input = filepath.Clean(input)
			if input == ".." || strings.HasPrefix(input, ".."+string(filepath.Separator)) || !policy.Allows(filepath.ToSlash(input)) {
				continue
			}
			path := filepath.Join(root, input)
			real, err := filepath.EvalSymlinks(path)
			realRoot, rootErr := filepath.EvalSymlinks(root)
			relative, relErr := filepath.Rel(realRoot, real)
			if err != nil || rootErr != nil || relErr != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
				continue
			}
			spec := SpringOpenAPISpec{path, strings.TrimSpace(c.APIPackage)}
			if !seen[spec] {
				out = append(out, spec)
				seen[spec] = true
			}
		}
	}
	return out
}
