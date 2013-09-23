##About

This is a client library for [Snowizard](https://github.com/GeneralElectric/snowizard) 

##Installing

Using go get

    $ go get github.com/pims/gosnowizard/snowizard

After this command snowizard is ready to use. Its source will be in:

    $GOROOT/src/pkg/github.com/pims/gosnowizard/snowizard

You can use go get -u -a for update all installed packages.

## Example

    import (
		"fmt"
		"github.com/pims/gosnowizard/snowizard"
		"log"
    )

    func main() {
		hosts := make([]string, 2)
		hosts[0] = "snowizard-1.dev:6776"
		hosts[1] = "snowizard-2.dev:6776"

		timeout := time.Duration(2 * time.Second)
		client := snowizard.NewSnowizardTextClient(hosts, timeout)

		id, err := client.Next()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id)
	}

##To do

Write tests.