package main

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/lxc/lxd/shared"
)

func mockStartDaemon() (*Daemon, error) {
	d := &Daemon{
		MockMode:              true,
		imagesDownloading:     map[string]chan bool{},
		imagesDownloadingLock: sync.RWMutex{},
	}

	if err := d.Init(); err != nil {
		return nil, err
	}

	d.IdmapSet = &shared.IdmapSet{Idmap: []shared.IdmapEntry{
		{Isuid: true, Hostid: 100000, Nsid: 0, Maprange: 500000},
		{Isgid: true, Hostid: 100000, Nsid: 0, Maprange: 500000},
	}}

	// Call this after Init so we have a log object.
	storageConfig := make(map[string]interface{})
	d.Storage = &storageLogWrapper{w: &storageMock{d: d}}
	if _, err := d.Storage.Init(storageConfig); err != nil {
		return nil, err
	}

	return d, nil
}

type lxdTestSuite struct {
	suite.Suite
	d      *Daemon
	Req    *require.Assertions
	tmpdir string
}

func (suite *lxdTestSuite) SetupSuite() {
	tmpdir, err := ioutil.TempDir("", "lxd_testrun_")
	if err != nil {
		os.Exit(1)
	}
	suite.tmpdir = tmpdir

	if err := os.Setenv("LXD_DIR", suite.tmpdir); err != nil {
		os.Exit(1)
	}

	suite.d, err = mockStartDaemon()
	if err != nil {
		os.Exit(1)
	}
}

func (suite *lxdTestSuite) TearDownSuite() {
	suite.d.Stop()

	err := os.RemoveAll(suite.tmpdir)
	if err != nil {
		os.Exit(1)
	}
}

func (suite *lxdTestSuite) SetupTest() {
	suite.Req = require.New(suite.T())
}

func TestLxdTestSuite(t *testing.T) {
	suite.Run(t, new(lxdTestSuite))
}
