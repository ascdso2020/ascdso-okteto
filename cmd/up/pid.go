// Copyright 2022 The Okteto Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package up

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/okteto/okteto/pkg/config"
	oktetoLog "github.com/okteto/okteto/pkg/log"
)

// createPIDFile creates a PID file to track Up state and existence
func createPIDFile(ns, dpName string) error {
	filePath := filepath.Join(config.GetAppHome(ns, dpName), "okteto.pid")
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("unable to create PID file at %s", filePath)
	}
	defer func() {
		if err := file.Close(); err != nil {
			oktetoLog.Debugf("Error closing file %s: %s", filePath, err)
		}
	}()
	if _, err := file.WriteString(strconv.Itoa(os.Getpid())); err != nil {
		return fmt.Errorf("unable to write to PID file at %s", filePath)
	}
	return nil
}

// cleanPIDFile deletes PID file after Up finishes
func cleanPIDFile(ns, dpName string) {
	filePath := filepath.Join(config.GetAppHome(ns, dpName), "okteto.pid")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		oktetoLog.Infof("unable to delete PID file at %s", filePath)
	}
}
