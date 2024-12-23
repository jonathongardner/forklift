# Extractors
To Add a new extractor, create a new subpackage in extractors:
```go
package NAME

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	// log "github.com/sirupsen/logrus"

	"github.com/jonathongardner/forklift/fs"
)

func ExtractArchive(virtualFS *fs.Virtual) error {
	f, err := virtualFS.RootOpen()
	if err != nil {
		return err
	}
	defer f.Close()

	...

	return nil
}
```

Then add it to the `extractors.go` file
```go
package extractors

import (
	...
	"github.com/jonathongardner/forklift/extractors/NAME"
)
...




func init() {
	//-----------------NAME-----------------
	{"some/content/type", NAME.ExtractArchive},
	//-----------------NAME-----------------
```