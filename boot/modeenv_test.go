// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2019 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package boot_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	. "gopkg.in/check.v1"

	"github.com/snapcore/snapd/boot"
	"github.com/snapcore/snapd/dirs"
	"github.com/snapcore/snapd/testutil"
)

// baseBootSuite is used to setup the common test environment
type modeenvSuite struct {
	testutil.BaseTest

	tmpdir string
}

var _ = Suite(&modeenvSuite{})

func (s *modeenvSuite) SetUpTest(c *C) {
	s.tmpdir = c.MkDir()
}

func (s *modeenvSuite) TestReadEmptyErrors(c *C) {
	modeenv, err := boot.ReadModeenv("/no/such/file")
	c.Assert(os.IsNotExist(err), Equals, true)
	c.Assert(modeenv, IsNil)
}

func (s *modeenvSuite) makeMockModeenvFile(c *C, content string) {
	mockModeenvPath := filepath.Join(s.tmpdir, dirs.SnapModeenvFile)
	err := os.MkdirAll(filepath.Dir(mockModeenvPath), 0755)
	c.Assert(err, IsNil)
	err = ioutil.WriteFile(mockModeenvPath, []byte(content), 0644)
	c.Assert(err, IsNil)
}

func (s *modeenvSuite) TestReadEmpty(c *C) {
	s.makeMockModeenvFile(c, "")

	modeenv, err := boot.ReadModeenv(s.tmpdir)
	c.Assert(err, IsNil)
	c.Check(modeenv.Mode, Equals, "")
	c.Check(modeenv.RecoverySystem, Equals, "")
}

func (s *modeenvSuite) TestReadMode(c *C) {
	s.makeMockModeenvFile(c, "mode: run")

	modeenv, err := boot.ReadModeenv(s.tmpdir)
	c.Assert(err, IsNil)
	c.Check(modeenv.Mode, Equals, "run")
	c.Check(modeenv.RecoverySystem, Equals, "")
}

func (s *modeenvSuite) TestReadModeWithRecoverySystem(c *C) {
	s.makeMockModeenvFile(c, `mode: recovery
recovery_system: 20191126
`)

	modeenv, err := boot.ReadModeenv(s.tmpdir)
	c.Assert(err, IsNil)
	c.Check(modeenv.Mode, Equals, "recovery")
	c.Check(modeenv.RecoverySystem, Equals, "20191126")
}

func (s *modeenvSuite) TestWriteExisting(c *C) {
	s.makeMockModeenvFile(c, "mode: run")

	modeenv, err := boot.ReadModeenv(s.tmpdir)
	c.Assert(err, IsNil)
	modeenv.Mode = "recovery"
	err = modeenv.Write(s.tmpdir)
	c.Assert(err, IsNil)
}
