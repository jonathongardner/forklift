# Extractors
To Add a new extractor, create a new subpackage:
```go
package NAME

import (
	"fmt"
	"io"

	"github.com/jonathongardner/forklift/extractors/helpers"
	<!-- log "github.com/sirupsen/logrus" -->

	"github.com/jonathongardner/forklift/fs"
)

func Add(add func(string, helpers.ExtratFunc)) {
		add("...", ExtractArchive)
}

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