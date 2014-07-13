package cmds

import "net/http"
import "os"
import "fmt"

import "xogeny/gimpact/utils"

func buildIndex() utils.Index {
	var master = "http://impact.modelica.org/impact_data.json";

	resp, err := http.Get(master)
	if err != nil {
		fmt.Println("Unable to locate index file at "+master);
		os.Exit(1);
	}
	defer resp.Body.Close()

	index := utils.Index{};

	err = index.BuildIndex(resp.Body);
	if (err!=nil) {
		fmt.Println("Error reading index: "+err.Error());
		os.Exit(2);
	}

	return index;
}

