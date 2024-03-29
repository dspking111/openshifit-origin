package defaults

import (
	"testing"

	kapi "k8s.io/kubernetes/pkg/api"

	buildadmission "github.com/openshift/origin/pkg/build/admission"
	u "github.com/openshift/origin/pkg/build/admission/testutil"
)

func TestProxyDefaults(t *testing.T) {
	defaultsConfig := &BuildDefaultsConfig{
		GitHTTPProxy:  "http",
		GitHTTPSProxy: "https",
	}

	admitter := NewBuildDefaults(defaultsConfig)
	pod := u.Pod().WithBuild(t, u.Build().WithDockerStrategy().AsBuild(), "v1")
	err := admitter.Admit(pod.ToAttributes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	build, _, err := buildadmission.GetBuild(pod.ToAttributes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if build.Spec.Source.Git.HTTPProxy == nil || len(*build.Spec.Source.Git.HTTPProxy) == 0 || *build.Spec.Source.Git.HTTPProxy != "http" {
		t.Errorf("failed to find http proxy in git source")
	}
	if build.Spec.Source.Git.HTTPSProxy == nil || len(*build.Spec.Source.Git.HTTPSProxy) == 0 || *build.Spec.Source.Git.HTTPSProxy != "https" {
		t.Errorf("failed to find http proxy in git source")
	}
}

func TestEnvDefaults(t *testing.T) {
	defaultsConfig := &BuildDefaultsConfig{
		Env: []kapi.EnvVar{
			{
				Name:  "VAR1",
				Value: "VALUE1",
			},
			{
				Name:  "VAR2",
				Value: "VALUE2",
			},
		},
	}

	admitter := NewBuildDefaults(defaultsConfig)
	pod := u.Pod().WithBuild(t, u.Build().WithSourceStrategy().AsBuild(), "v1")
	err := admitter.Admit(pod.ToAttributes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	build, _, err := buildadmission.GetBuild(pod.ToAttributes())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := getBuildEnv(build)
	var1found, var2found := false, false
	for _, ev := range *env {
		if ev.Name == "VAR1" {
			if ev.Value != "VALUE1" {
				t.Errorf("unexpected value %s", ev.Value)
			}
			var1found = true
		}
		if ev.Name == "VAR2" {
			if ev.Value != "VALUE2" {
				t.Errorf("unexpected value %s", ev.Value)
			}
			var2found = true
		}
	}
	if !var1found {
		t.Errorf("VAR1 not found")
	}
	if !var2found {
		t.Errorf("VAR2 not found")
	}
}
