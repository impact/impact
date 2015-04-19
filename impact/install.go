package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/opesun/copyrecur"
	"github.com/pierrre/archivefile/zip"
	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/graph"
	"github.com/xogeny/impact/index"
)

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" login:"verbose" description:"Turn on verbose output"`
	DryRun  bool `short:"d" login:"dryrun" description:"Resolve dependencies but don't install"`
}

/*
// TODO: This probably needs to be refactored
func ParseVersion(libver string, ind index.Index) (libname index.LibraryName,
	ver index.VersionString, err error) {
	parts := strings.Split(libver, "#")

	libname = index.LibraryName(parts[0])

	lib, ok := ind[libname]
	if !ok {
		err = index.MissingLibraryError{Name: libname}
		return
	}

	if len(parts) == 1 {
		version, lerr := lib.Latest()
		if lerr != nil {
			err = lerr
			return
		}
		ver = index.VersionString(version.Version.String())
	} else if len(parts) == 2 {
		ver = index.VersionString(parts[1])
	} else if len(parts) > 2 {
		err = errors.New("Invalid version specification: " + libver +
			" (must be libraryName#version)")
		return
	}
	return
}
*/

func (x *InstallCommand) Execute(args []string) error {
	//Check to make sure we have something to install
	if len(args) == 0 {
		return errors.New("No libraries requested for installation")
	}

	// Load index
	ind, err := index.LoadIndex()
	if err != nil {
		return fmt.Errorf("Error loading indices: %v", err)
	}

	// Build dependency graph from index
	resolver, err := ind.BuildGraph(x.Verbose)
	if err != nil {
		return fmt.Errorf("Error building dependency graph: %v", err)
	}

	// State root dependencies
	libnames := []graph.LibraryName{}
	for _, n := range args {
		libnames = append(libnames, graph.LibraryName(n))
	}

	// Resolve dependencies
	solution, err := resolver.Resolve(libnames...)
	if err != nil {
		return fmt.Errorf("Error resolving dependencies for %v: %v", libnames, err)
	}

	// Install dependencies
	color.Printf("@{y}Installing...\n")
	for name, version := range solution {
		color.Printf("  Library: @{g}%s\n", name)
		color.Printf("    Required version: @{g}%v\n", version)
		lv, err := ind.Find(string(name), version)
		if err != nil {
			return fmt.Errorf("Couldn't find version %v of library %s (this should not happen)",
				version, name)
		}
		if !x.DryRun {
			Install(lv, ind, ".", x.Verbose)
		}
	}

	return nil
}

func Install(ver index.VersionDetails, ind *index.Index, target string, verbose bool) error {
	/* Download the Zipball to a temporary file */
	if verbose {
		color.Println("  @{y}Downloading source from: @{!y}" + string(ver.Zipball))
	}

	/*   Do a GET request */
	resp, err := http.Get(ver.Zipball)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // Make sure this gets closed

	/*   Open a temporary file to direct the download into */
	tzf, err := ioutil.TempFile("", "impact")
	defer func() {
		tzf.Close()           // Make sure we close this file and...
		os.Remove(tzf.Name()) // ...delete it.
	}()
	/*   Copy the bytes to temporary file */
	zsize, err := io.Copy(tzf, resp.Body)
	if err != nil {
		return err
	}

	/* Create a temporary directory to extract into */
	tdir, err := ioutil.TempDir("", "impact")
	defer func() {
		os.RemoveAll(string(tdir)) // Make sure this gets removed in case of a panic
	}()
	if err != nil {
		return err
	}

	/* Extract the zip file into our temporary directory */
	var adir string = ""
	err = zip.Unarchive(tzf, zsize, string(tdir), func(x string) {
		if adir == "" {
			adir = strings.Split(x, "/")[0]
		}
	})
	if err != nil {
		return err
	}

	/* Figure out where the Modelica code is in our temporary directory */
	keep := path.Join(string(tdir), adir, ver.Path)

	/* Figure out whether we are dealing with a package stored as a file or diretory */
	fi, err := os.Stat(keep)
	if err != nil {
		return err
	}

	/* Copy the Modelica code to our target installation directory */
	if fi.IsDir() {
		if verbose {
			color.Println("  @{y}Copying  @{!y}" + ver.Path +
				"@{y} to @{!y}" + path.Join(target, fi.Name()))
		}
		copyrecur.CopyDir(keep, path.Join(target, fi.Name()))
	} else {
		if verbose {
			color.Println("  @{y}Copying  @{!y}" + ver.Path +
				"@{y} to @{!y}" + path.Join(target, fi.Name()))
		}
		copyrecur.CopyFile(keep, path.Join(target, fi.Name()))
	}

	return nil
}
