package cli

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/spf13/pflag"

	kcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
	"k8s.io/kubernetes/pkg/util/sets"

	"github.com/openshift/origin/pkg/cmd/util/clientcmd"
)

// MissingCommands is the list of commands we're already missing.
// NEVER ADD TO THIS LIST
// TODO kill this list
var MissingCommands = sets.NewString("namespace", "rolling-update", "cluster-info", "api-versions",
	"apply",     // we don't want to support this implementation
	"autoscale", // TODO
)

// WhitelistedCommands is the list of commands we're never going to have in oc
// defend each one with a comment
var WhitelistedCommands = sets.NewString()

func TestKubectlCompatibility(t *testing.T) {
	f := clientcmd.New(pflag.NewFlagSet("name", pflag.ContinueOnError))

	oc := NewCommandCLI("oc", "oc", &bytes.Buffer{}, ioutil.Discard, ioutil.Discard)
	kubectl := kcmd.NewKubectlCommand(f.Factory, nil, ioutil.Discard, ioutil.Discard)

kubectlLoop:
	for _, kubecmd := range kubectl.Commands() {
		for _, occmd := range oc.Commands() {
			if kubecmd.Name() == occmd.Name() {
				if MissingCommands.Has(kubecmd.Name()) {
					t.Errorf("%s was supposed to be missing", kubecmd.Name())
					continue
				}
				if WhitelistedCommands.Has(kubecmd.Name()) {
					t.Errorf("%s was supposed to be whitelisted", kubecmd.Name())
					continue
				}
				continue kubectlLoop
			}
		}
		if MissingCommands.Has(kubecmd.Name()) || WhitelistedCommands.Has(kubecmd.Name()) {
			continue
		}

		t.Errorf("missing %q in oc", kubecmd.Name())
	}

}

// this only checks one level deep for nested commands, but it does ensure that we've gotten several
// --validate flags.  Based on that we can reasonably assume we got them in the kube commands since they
// all share the same registration.
func TestValidateDisabled(t *testing.T) {
	f := clientcmd.New(pflag.NewFlagSet("name", pflag.ContinueOnError))

	oc := NewCommandCLI("oc", "oc", &bytes.Buffer{}, ioutil.Discard, ioutil.Discard)
	kubectl := kcmd.NewKubectlCommand(f.Factory, nil, ioutil.Discard, ioutil.Discard)

	for _, kubecmd := range kubectl.Commands() {
		for _, occmd := range oc.Commands() {
			if kubecmd.Name() == occmd.Name() {
				ocValidateFlag := occmd.Flags().Lookup("validate")
				if ocValidateFlag == nil {
					continue
				}

				if ocValidateFlag.Value.String() != "false" {
					t.Errorf("%s --validate is not defaulting to false", occmd.Name())
				}
			}
		}
	}

}
