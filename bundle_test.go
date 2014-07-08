// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charm_test

import (
	"io/ioutil"
	"path/filepath"

	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "launchpad.net/gocheck"

	"gopkg.in/juju/charm.v2"
	charmtesting "gopkg.in/juju/charm.v2/testing"
)

var _ = gc.Suite(&BundleSuite{})

type BundleSuite struct {
	testing.IsolationSuite
}

func (*BundleSuite) TestReadBundleDir(c *gc.C) {
	path := charmtesting.Charms.BundleDirPath("wordpress")
	b, err := charm.ReadBundle(path, func(string) error { return nil })
	c.Assert(err, gc.IsNil)
	c.Assert(b, gc.FitsTypeOf, (*charm.BundleDir)(nil))
	checkWordpressBundle(c, b, path)
}

func (*BundleSuite) TestReadBundleArchive(c *gc.C) {
	path := filepath.Join(c.MkDir(), "something")
	err := ioutil.WriteFile(path, []byte("oh yeah"), 0644)
	c.Assert(err, gc.IsNil)
	c.Assert(func() {
		charm.ReadBundle(path, func(string) error { return nil })
	}, gc.PanicMatches, "unimplemented")
}

func checkWordpressBundle(c *gc.C, b charm.Bundle, path string) {
	// Load the charms required by the bundle.
	wordpressCharm := charmtesting.Charms.CharmDir("wordpress")
	mysqlCharm := charmtesting.Charms.CharmDir("mysql")

	bd := b.Data()
	c.Assert(bd.RequiredCharms(), jc.DeepEquals, []string{"mysql", "wordpress"})

	charms := map[string]charm.Charm{
		"wordpress": wordpressCharm,
		"mysql":     mysqlCharm,
	}
	err := bd.VerifyWithCharms(func(string) error { return nil }, charms)
	c.Assert(err, gc.IsNil)

	c.Assert(bd.Services, jc.DeepEquals, map[string]*charm.ServiceSpec{
		"wordpress": {
			Charm:    "wordpress",
			NumUnits: 1,
		},
		"mysql": {
			Charm:    "mysql",
			NumUnits: 1,
		},
	})
	c.Assert(bd.Relations, jc.DeepEquals, [][]string{
		{"wordpress:db", "mysql:server"},
	})
	c.Assert(b.ReadMe(), gc.Equals, "A dummy bundle\n")
	switch b := b.(type) {
	case *charm.BundleArchive:
		c.Assert(b.Path, gc.Equals, path)
	case *charm.BundleDir:
		c.Assert(b.Path, gc.Equals, path)
	}
}
