package main

import "fmt"
import "xogeny/gimpact/utils"
import "net/http"

func main() {
	var master = "http://impact.modelica.org/impact_data.json";

	resp, err := http.Get(master)
	if err != nil {
		fmt.Println("Unable to locate index file at "+master);
		return;
	}
	defer resp.Body.Close()

	index := utils.Index{};

	err = index.BuildIndex(resp.Body);
	if (err!=nil) {
		fmt.Println("Error reading index: "+err.Error());
	} else {
		fmt.Println(index);
	}
}
