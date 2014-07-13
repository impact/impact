package main

import "fmt"
import "encoding/json"
import "xogeny/gimpact/utils"

func main() {
	var ds = `{
                  "version": "3.2", 
                  "name": "Modelica"
              }`
	sample := []byte(ds);
	var dep utils.Dependency;
	dep = utils.Dependency{};
	json.Unmarshal(sample, &dep);
	fmt.Println(dep);

	index := utils.Index{};

	err := index.ReadIndex("sample.json");
	if (err!=nil) {
		fmt.Println("Error reading file: "+err.Error());
	} else {
		fmt.Println(index);
	}
}
