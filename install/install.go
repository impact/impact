package install

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/opesun/copyrecur"
	"github.com/pierrre/archivefile/zip"
	"github.com/wsxiaoys/terminal/color"

	"github.com/xogeny/impact/index"
)

func Install(libname string, ver index.VersionDetails, ind *index.Index,
	target string, verbose bool) error {
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

	dst := path.Join(target, libname)
	/* Copy the Modelica code to our target installation directory */
	if fi.IsDir() {
		if verbose {
			color.Printf("  @{y}Copying  @{!y}%s@{y} to @{!y}%s\n", ver.Path, dst)
		}
		copyrecur.CopyDir(keep, path.Join(target, libname))
	} else {
		if verbose {
			color.Printf("  @{y}Copying  @{!y}%s@{y} to @{!y}%s\n", ver.Path, dst)
		}
		copyrecur.CopyFile(keep, path.Join(target, fi.Name()))
	}

	return nil
}
