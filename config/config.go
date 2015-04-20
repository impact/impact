package config

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/xogeny/denada-go"

	"github.com/mitchellh/go-homedir"

	"github.com/xogeny/impact/crawl"
)

// If this environment variable is set, we use the file it
// points to as the name of the user's settings file
var envvar = "IMPACT_CONFIG_FILE"

func SettingsFile() string {
	// Look to see if IMPACT_CONFIG_FILE is set.  If so, read the
	// file it points to.
	if os.Getenv(envvar) != "" {
		return os.Getenv(envvar)
	}

	home, _ := homedir.Dir()

	// Otherwise, find out what platform we are on...
	platform := runtime.GOOS

	datadir := ""

	switch platform {
	case "windows":
		// On windows, check to see if APPDATA is defined...
		datadir = os.Getenv("APPDATA")
		if datadir == "" {
			datadir = path.Join(home, ".config")
		}
	case "linux":
		// On windows, check to see if APPDATA is defined...
		datadir = os.Getenv("XDG_CONFIG_HOME")
		if datadir == "" {
			datadir = path.Join(home, ".config")
		}
	case "darwin":
		datadir = path.Join(home, "Library", "Preferences")
	default:
		log.Printf("Unknown platform %v", platform)
	}

	return path.Join(datadir, "impact", "impactrc")
}

var syntax = `
index = "$string" "indices*";
github source = "$string" "sources*";
`

func ReadSettings() (Settings, error) {
	blank := MakeSettings()
	sfile := SettingsFile()

	settings, err := denada.ParseFile(sfile)
	if err != nil {
		log.Printf("Ignoring settings file, got error while parsing: %v", err)
		settings = denada.ElementList{}
	}

	grammar, err := denada.ParseString(syntax)
	if err != nil {
		return blank,
			fmt.Errorf("Error parsing settings file syntax specification: %v", err)
	} else {
		err = denada.Check(settings, grammar, false)
		if err != nil {
			return blank,
				fmt.Errorf("Error in settings file: %v", err)
		}
	}

	ret := MakeSettings()

	// Parse index sources...
	indices := settings.OfRule("indices", true)
	for _, index := range indices {
		url := index.StringValueOf("")
		ret.Indices = append(ret.Indices, url)
	}

	// If no indices were specified, use the default:
	if len(ret.Indices) == 0 {
		ret.Indices = []string{"https://impact.modelica.org/impact_data2.json"}
	}

	mo := denada.NewDeclaration("source", "", "github")
	mo.SetValue("modelica/.+")
	mo3 := denada.NewDeclaration("source", "", "github")
	mo3.SetValue("modelica-3rdparty/.+")

	// Now parse sources
	sources := settings.OfRule("sources", true)

	if len(sources) == 0 {
		sources = append(sources, mo, mo3)
	}

	for _, source := range sources {
		val := source.StringValueOf("")

		scheme := source.Qualifiers[0]

		switch scheme {
		case "github":
			path := strings.Split(val, "/")
			switch len(path) {
			case 1:
				c, err := crawl.MakeGitHubCrawler(path[0], "", "")
				if err != nil {
					return blank,
						fmt.Errorf("Unable to create GitHub crawler from %s: %v",
							val, err)
				}
				ret.Sources = append(ret.Sources, c)
			case 2:
				c, err := crawl.MakeGitHubCrawler(path[0], path[1], "")
				if err != nil {
					return blank,
						fmt.Errorf("Unable to create GitHub crawler from %s: %v",
							val, err)
				}
				ret.Sources = append(ret.Sources, c)
			default:
				return blank,
					fmt.Errorf("GitHub source syntax: repo{/pattern}, found %s",
						val)
			}

		default:
			return blank,
				fmt.Errorf("Unrecognized scheme in source %s, expected 'github'",
					val)
		}
	}

	return ret, nil
}
