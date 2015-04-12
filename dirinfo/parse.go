package dirinfo

import (
	"encoding/json"
)

func Parse(str string) (DirectoryInfo, error) {
	ret := MakeDirectoryInfo()
	blank := MakeDirectoryInfo()

	err := json.Unmarshal([]byte(str), &ret)
	if err != nil {
		return blank, err
	}

	return ret, nil
}
