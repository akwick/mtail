// Copyright 2019 Google Inc. All Rights Reserved.
// This file is available under the Apache license.

package mtail_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/golang/glog"
	"github.com/google/mtail/internal/mtail"
	"github.com/google/mtail/internal/testutil"
)

func TestMultipleLinesInOneWrite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	for _, test := range mtail.LogWatcherTestTable {
		t.Run(fmt.Sprintf("%s %v", test.PollInterval, test.EnableFsNotify), func(t *testing.T) {
			tmpDir, rmTmpDir := testutil.TestTempDir(t)
			defer rmTmpDir()

			logDir := path.Join(tmpDir, "logs")
			progDir := path.Join(tmpDir, "progs")
			err := os.Mkdir(logDir, 0700)
			testutil.FatalIfErr(t, err)
			err = os.Mkdir(progDir, 0700)
			testutil.FatalIfErr(t, err)

			logFile := path.Join(logDir, "log")

			f := testutil.TestOpenFile(t, logFile)

			m, stopM := mtail.TestStartServer(t, test.PollInterval, test.EnableFsNotify, mtail.ProgramPath(progDir), mtail.LogPathPatterns(logDir+"/log"))
			defer stopM()

			{
				lineCountCheck := m.ExpectMetricDeltaWithDeadline("lines_total", 1)
				n, err := f.WriteString("line 1\n")
				testutil.FatalIfErr(t, err)
				glog.Infof("Wrote %d bytes", n)
				lineCountCheck()
			}

			{
				lineCountCheck := m.ExpectMetricDeltaWithDeadline("lines_total", 2)
				n, err := f.WriteString("line 2\nline 3\n")
				testutil.FatalIfErr(t, err)
				glog.Infof("Wrote %d bytes", n)
				lineCountCheck()
			}
		})
	}
}
