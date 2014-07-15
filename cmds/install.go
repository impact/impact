package cmds

import "os"
import "io"
import "fmt"
import "path"
import "errors"
import "strings"
import "net/http"
import "io/ioutil"

import "xogeny/gimpact/utils"

import "github.com/wsxiaoys/terminal/color"
import "github.com/pierrre/archivefile/zip"
import "github.com/opesun/copyrecur"

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
		ierr := install(d.Name, d.Version, index, ".", x.Verbose);
		if (ierr!=nil) { color.Println("@{r}Error: "+ierr.Error()) };
	}
	
	return nil;
}

func install(name utils.LibraryName, version utils.VersionString,
	index utils.Index, target string, verbose bool) error {
	// Create temporary directory
	tdir, err := ioutil.TempDir("", "gimpact");
	defer func() {
		//os.RemoveAll(string(tdir));
	}()
	if (err!=nil) { return err; }

	// TODO: Keep a cache

	// Extract Version data
	lib, ok := index[name];
	if (!ok) { return utils.MissingLibraryError{Name:name}; }
	ver, ok := lib.Versions[version];
	if (!ok) { return utils.MissingVersionError{Name:name, Version:version}; }
	if (verbose) { fmt.Println("Installing "+string(name)+", version "+string(version)); }

	// Identify Zip file
	zipfile := ver.Zipball;

	// Download Zip file
	if (verbose) { fmt.Println("  Downloading source from: "+string(zipfile)); }
	resp, err := http.Get(zipfile);
	if (err!=nil) { return err; }
	//fmt.Println(resp);
	defer resp.Body.Close();
	tzf, err := ioutil.TempFile("", "gimpact");
	defer func() {
		tzf.Close();
		//os.Remove(tzf.Name());
	}();
	zsize, err := io.Copy(tzf, resp.Body);
	if (err!=nil) { return err; }
	resp.Body.Close();
	if (verbose) { fmt.Println("  Temporary zip file stored at: "+tzf.Name()); }

	// Extract Zip file
	var adir string = "";
	if (verbose) { fmt.Println("  Extracting to: "+string(tdir)); }
	err = zip.Unarchive(tzf, zsize, string(tdir), func(x string) {
		if (adir=="") { adir = strings.Split(x, "/")[0]; }
		//fmt.Println("    Extracting: "+x)
	})
	if (err!=nil) { return err; }

	// Copy the 'path' content to the install directory
	keep := path.Join(string(tdir), adir, ver.Path);

	f, err := os.Open(keep);
	defer f.Close();
	if (err!=nil) { return err; }
	fi, err := f.Stat();
	if (err!=nil) { return err; }
	f.Close();
	if (fi.IsDir()) {
		if (verbose) {
			fmt.Println("  Copying directory "+keep+" to "+path.Join(target, fi.Name()));
		}
		copyrecur.CopyDir(keep, path.Join(target, fi.Name()));
	} else {
		copyrecur.CopyFile(keep, target);
	}

	// Get rid of temporary directory
	//os.RemoveAll(string(tdir));

	// Get rid of temporary file
	tzf.Close();
	//os.Remove(tzf.Name());

	return nil;
}
