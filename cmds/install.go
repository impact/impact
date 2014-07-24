package cmds

import "os"
import "io"
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
	/* Check to make sure we have something to install */
	if (len(args)==0) {
		return errors.New("No libraries requested for installation");
	}

	/* Build the index */
	index := buildIndex();

	/* Create an empty set of libraries */
	todo := utils.Libraries{};

	/* Loop over all the libraries to be installed */
	for _, arg := range(args) {
		libname, ver, err := ParseVersion(arg, index);
		if (err!=nil) { return err; }
		
		/* Get Version objects for this library and all its dependencies */
		deps, err := index.Dependencies(libname, ver);
		if (err!=nil) { return err; }
	
		/* Merge them with the master list of libraries we are going to install */
		err = todo.Merge(deps);
		if (err!=nil) { return err; }
	}

	/* Loop over all the libraries we have identified for installation and install them */
	for ln, lv := range(todo) {
		if (x.Verbose) {
			color.Println("@!Installing @{c}"+string(ln)+", version "+string(lv.Version));
		}
		ierr := Install(lv, index, ".", x.Verbose);
		if (ierr!=nil) { color.Println("@{r}Error: "+ierr.Error()) };
	}
	
	return nil;
}

func Install(ver utils.Version,	index utils.Index, target string, verbose bool) error {
	/* Download the Zipball to a temporary file */
	if (verbose) { color.Println("  @{y}Downloading source from: @{!y}"+string(ver.Zipball)); }

	/*   Do a GET request */
	resp, err := http.Get(ver.Zipball);
	if (err!=nil) { return err; }
	defer resp.Body.Close(); // Make sure this gets closed

	/*   Open a temporary file to direct the download into */
	tzf, err := ioutil.TempFile("", "gimpact");
	defer func() {
		tzf.Close(); // Make sure we close this file and...
		os.Remove(tzf.Name()); // ...delete it.
	}();
	/*   Copy the bytes to temporary file */
	zsize, err := io.Copy(tzf, resp.Body);
	if (err!=nil) { return err; }

	/* Create a temporary directory to extract into */
	tdir, err := ioutil.TempDir("", "gimpact");
	defer func() {
		os.RemoveAll(string(tdir)); // Make sure this gets removed in case of a panic
	}()
	if (err!=nil) { return err; }

	/* Extract the zip file into our temporary directory */
	var adir string = "";
	err = zip.Unarchive(tzf, zsize, string(tdir), func(x string) {
		if (adir=="") { adir = strings.Split(x, "/")[0]; }
	})
	if (err!=nil) { return err; }

	/* Figure out where the Modelica code is in our temporary directory */
	keep := path.Join(string(tdir), adir, ver.Path);

	/* Figure out whether we are dealing with a package stored as a file or diretory */
	fi, err := os.Stat(keep);
	if (err!=nil) { return err; }

	/* Copy the Modelica code to our target installation directory */
	if (fi.IsDir()) {
		if (verbose) {
			color.Println("  @{y}Copying  @{!y}"+ver.Path+
				"@{y} to @{!y}"+path.Join(target, fi.Name()));
		}
		copyrecur.CopyDir(keep, path.Join(target, fi.Name()));
	} else {
		if (verbose) {
			color.Println("  @{y}Copying  @{!y}"+ver.Path+
				"@{y} to @{!y}"+path.Join(target, fi.Name()));
		}
		copyrecur.CopyFile(keep, path.Join(target, fi.Name()));
	}

	return nil;
}
