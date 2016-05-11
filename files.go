package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func files() {
	d, _ := ioutil.ReadDir("./omega")
	for i := 0; i < len(d); i++ {
		file := d[i]
		r := strings.Replace(file.Name(), "omega.", "", -1)
		fmt.Println(r)
		os.Rename(fmt.Sprintf("./omega/%v", file.Name()), fmt.Sprintf("./omega/%v", strings.Replace(file.Name(), "omega.", "", -1)))
	}
}
