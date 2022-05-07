package cfg

import (
	"context"
	"fmt"
	"os"

	"github.com/mediocregopher/mediocre-go-lib/v2/mctx"
)

// DataDir manages the blog's data directory.
type DataDir struct {
	Path string

	deleteOnClose bool
}

// Init initializes the data directory, creating the directory named at path if
// it doesn't exist.
//
// If Path is not set, then a temporary directory will be created and its path
// set to the Path field. This directory will be removed when Close is called.
func (d *DataDir) Init() error {
	if d.Path == "" {

		d.deleteOnClose = true
		var err error

		if d.Path, err = os.MkdirTemp("", "mediocre-blog-data-*"); err != nil {
			return fmt.Errorf("creating temporary directory: %w", err)
		}

		return nil
	}

	if err := os.MkdirAll(d.Path, 0700); err != nil {
		return fmt.Errorf(
			"creating directory (and parents) of %q: %w",
			d.Path,
			err,
		)
	}

	return nil
}

// SetupCfg implement the cfg.Cfger interface.
func (d *DataDir) SetupCfg(cfg *Cfg) {

	cfg.StringVar(&d.Path, "data-dir", "", "Directory to use for persistent storage. If unset a temp directory will be created, and will be deleted when the process exits.")

	cfg.OnInit(func(ctx context.Context) error {
		return d.Init()
	})
}

// Annotate implements mctx.Annotator interface.
func (d *DataDir) Annotate(a mctx.Annotations) {
	a["dataDirPath"] = d.Path
}

// Close cleans up any temporary state created by DataDir.
func (d *DataDir) Close() error {

	if !d.deleteOnClose {
		return nil
	}

	if err := os.RemoveAll(d.Path); err != nil {
		return fmt.Errorf("removing temp dir %q: %w", d.Path, err)
	}

	return nil
}
