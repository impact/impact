package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/xogeny/impact/graph"
	"github.com/xogeny/impact/utils"

	"github.com/blang/semver"
	"github.com/opesun/copyrecur"
	"github.com/pierrre/archivefile/zip"
	"github.com/wsxiaoys/terminal/color"
)

/* Define a struct listing all command line options for 'install' */
type InstallCommand struct {
	Verbose bool `short:"v" login:"verbose" description:"Turn on verbose output"`
	DryRun  bool `short:"d" login:"dryrun" description:"Resolve dependencies but don't install"`
}

func ParseVersion(libver string, index utils.Index) (libname utils.LibraryName,
	ver utils.VersionString, err error) {
	parts := strings.Split(libver, "#")

	libname = utils.LibraryName(parts[0])

	lib, ok := index[libname]
	if !ok {
		err = utils.MissingLibraryError{Name: libname}
		return
	}

	if len(parts) == 1 {
		version, lerr := lib.Latest()
		if lerr != nil {
			err = lerr
			return
		}
		ver = version.Version
	} else if len(parts) == 2 {
		ver = utils.VersionString(parts[1])
	} else if len(parts) > 2 {
		err = errors.New("Invalid version specification: " + libver +
			" (must be libraryName#version)")
		return
	}
	return
}

func (x *InstallCommand) Execute(args []string) error {
	/* Check to make sure we have something to install */
	if len(args) == 0 {
		return errors.New("No libraries requested for installation")
	}

	/* Build the index */
	index := utils.DownloadIndex()

	/* Create an empty set of libraries */
	var resolver graph.Resolver = graph.NewLibraryGraph()

	for libname, lib := range index {
		for _, ver := range lib.Versions {
			// TODO: If lib.Versions already used semver, this wouldn't be necessary
			// Should probably make a 'type Version ...' globally and utility functions
			// to parse and compare.  Then the implementation is a behind the scenes
			// detail we can change globally.
			v := semver.Version{
				Major: uint64(ver.Major),
				Minor: uint64(ver.Minor),
				Patch: uint64(ver.Patch),
			}
			resolver.AddLibrary(graph.LibraryName(libname), &v)
			for _, dep := range ver.Dependencies {
				dlib, err := index.Find(dep.Name, dep.Version)
				if err != nil {
					//log.Printf("Unable to find version %s of library %s",
					//  dep.Version, dep.Name)
					continue
				}
				dv := semver.Version{
					Major: uint64(dlib.Major),
					Minor: uint64(dlib.Minor),
					Patch: uint64(dlib.Patch),
				}

				//log.Printf("%s %s -> %s %s", libname, v.String(), dep.Name, dv.String())

				deplib := graph.LibraryName(dep.Name)
				depver := &dv

				if !resolver.Contains(deplib, depver) {
					resolver.AddLibrary(deplib, depver)
				}

				err = resolver.AddDependency(graph.LibraryName(libname), &v, deplib, depver)
				if err != nil {
					return err
				}
			}
		}
	}

	// TODO: Add a universal LibraryName alias
	libnames := []graph.LibraryName{}
	for _, n := range args {
		libnames = append(libnames, graph.LibraryName(n))
	}

	config, err := resolver.Resolve(libnames...)
	if err != nil {
		log.Printf("Error resolving libraries: %v", err)
	}
	if x.Verbose {
		color.Println("Libraries to be installed:")
		for name, ver := range config {
			color.Printf("  @{c}%s %v\n", name, ver)
		}
		color.Printf("\n")
	}

	/* Loop over all the libraries we have identified for installation and install them */
	for ln, v := range config {
		lib, exists := index[utils.LibraryName(ln)]
		if !exists {
			fmt.Printf("Unable to locate library named %s (this should not happen)", ln)
			return fmt.Errorf("Unable to locate library named %s (this should not happen)", ln)
		}
		var lv utils.Version
		found := false
		for _, ver := range lib.Versions {
			if uint64(ver.Major) == v.Major && uint64(ver.Minor) == v.Minor && uint64(ver.Patch) == v.Patch {
				found = true
				lv = ver
				break
			}
		}

		if !found {
			fmt.Printf("Unable to locate version %s of library %s",
				v.String(), ln)
			return fmt.Errorf("Unable to locate version %s of library %s",
				v.String(), ln)
		}
		if x.Verbose {
			color.Println("@!Installing @{c}%s, version %s...", string(ln), string(lv.Version))
		}

		// If this is just a DryRun, don't actually install
		if !x.DryRun {
			ierr := Install(lv, index, ".", x.Verbose)
			if ierr != nil {
				color.Println("@{r}Error: " + ierr.Error())
			}
		}

		if x.Verbose {
			color.Println("...@{!g}done")
		}

	}

	if x.Verbose {
		color.Println("\n@{!g}All libraries installed.\n")
	}

	return nil
}

func Install(ver utils.Version, index utils.Index, target string, verbose bool) error {
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
