//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cli

import (
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/scheme"
)

const logLevelEnvKey = "K8S_MANIFEST_SIGSTORE_LOG_LEVEL"

var logLevelMap = map[string]log.Level{
	"panic": log.PanicLevel,
	"fatal": log.FatalLevel,
	"error": log.ErrorLevel,
	"warn":  log.WarnLevel,
	"info":  log.InfoLevel,
	"debug": log.DebugLevel,
	"trace": log.TraceLevel,
}

var KOptions KubectlOptions

var RootCmd = &cobra.Command{
	Use:   "kubectl-sigstore",
	Short: "A command to sign/verify Kubernetes YAML manifests and resources on cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("kubectl sigstore cannot be invoked without a subcommand operation")
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("operation must be specified (e.g. kubectl sigstore sign)")
		}
		return nil
	},
}

func init() {
	KOptions = KubectlOptions{
		// generic options
		ConfigFlags: genericclioptions.NewConfigFlags(true),
		PrintFlags:  genericclioptions.NewPrintFlags("created").WithTypeSetter(scheme.Scheme),
	}

	RootCmd.AddCommand(NewCmdSign())
	RootCmd.AddCommand(NewCmdVerify())
	RootCmd.AddCommand(NewCmdApplyAfterVerify())
	RootCmd.AddCommand(NewCmdManifestBuild())
	RootCmd.AddCommand(NewCmdVersion())

	logLevelStr := os.Getenv(logLevelEnvKey)
	if logLevelStr == "" {
		logLevelStr = "info"
	}
	logLevel, ok := logLevelMap[logLevelStr]
	if !ok {
		logLevel = log.InfoLevel
	}

	log.SetLevel(logLevel)
}
