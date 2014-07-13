package cmds

import "github.com/wsxiaoys/terminal/color"
import "errors"
import "strings"
import "xogeny/gimpact/utils"

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" login:"verbose" description:"Turn on verbose output"`
}

var Install InstallCommand; // Instantiate option struct

func ParseVersion(libver string, index utils.Index) (libname utils.LibraryName,
	ver utils.VersionString, err error) {
	parts := strings.Split(libver, "#");

	libname = utils.LibraryName(parts[0]);

	lib, ok := index[libname];
	if (!ok) {
		err = utils.MissingLibraryError{Name:libname};
		return;
	}

	if (len(parts)==1) {
		version, lerr := lib.Latest();
		if (lerr!=nil) { err = lerr; return; }
		ver = version.Version;
	} else if (len(parts)==2) {
		ver = utils.VersionString(parts[1]);
	} else if (len(parts)>2) {
		err = errors.New("Invalid version specification: "+libver+
			" (must be libraryName#version)");
		return;
	}
	return;
}

func (x *InstallCommand) Execute(args []string) error {
	if (len(args)==0) {
		return errors.New("No libraries requested for installation");
	}

	todo := []utils.Dependency{};

	index := buildIndex();

	for _, arg := range(args) {
		var libname utils.LibraryName;
		var ver utils.VersionString;

		parts := strings.Split(arg, "#");

		libname = utils.LibraryName(parts[0]);

		lib, ok := index[libname];
		if (!ok) {
			return utils.MissingLibraryError{Name:libname};
		}

		if (len(parts)==1) {
			version, err := lib.Latest();
			if (err!=nil) { return err; }
			ver = version.Version;
		} else if (len(parts)==2) {
			ver = utils.VersionString(parts[1]);
		} else if (len(parts)>2) {
			return errors.New("Invalid version specification: "+arg+
				" (must be libraryName#version)");
		}
		

		_, ok = lib.Versions[ver];
		if (!ok) {
			return utils.MissingVersionError{Name:libname, Version: ver};
		}

		todo = append(todo, utils.Dependency{Name:libname, Version: ver});
	}

	for _, d := range(todo) {
		if (x.Verbose) {
			color.Println("@!Installing @{c}"+string(d.Name)+", version "+string(d.Version));
		}
		install(d.Name, d.Version, index);
	}
	
	return nil;
}

func install(name utils.LibraryName, version utils.VersionString, index utils.Index) error {
	// Create temporary directory
	// Unzip file into temporary directory
	// Copy the 'path' content to the install directory
	return nil;
}
