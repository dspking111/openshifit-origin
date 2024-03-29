package validation

import (
	"strings"
	"testing"

	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/util/validation/field"

	buildapi "github.com/openshift/origin/pkg/build/api"
)

func TestBuildValidationSuccess(t *testing.T) {
	build := &buildapi.Build{
		ObjectMeta: kapi.ObjectMeta{Name: "buildid", Namespace: "default"},
		Spec: buildapi.BuildSpec{
			Source: buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "http://github.com/my/repository",
				},
				ContextDir: "context",
			},
			Strategy: buildapi.BuildStrategy{
				DockerStrategy: &buildapi.DockerBuildStrategy{},
			},
			Output: buildapi.BuildOutput{
				To: &kapi.ObjectReference{
					Kind: "DockerImage",
					Name: "repository/data",
				},
			},
		},
		Status: buildapi.BuildStatus{
			Phase: buildapi.BuildPhaseNew,
		},
	}
	if result := ValidateBuild(build); len(result) > 0 {
		t.Errorf("Unexpected validation error returned %v", result)
	}
}

func TestBuildValidationFailure(t *testing.T) {
	build := &buildapi.Build{
		ObjectMeta: kapi.ObjectMeta{Name: "", Namespace: ""},
		Spec: buildapi.BuildSpec{
			Source: buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "http://github.com/my/repository",
				},
				ContextDir: "context",
			},
			Strategy: buildapi.BuildStrategy{
				DockerStrategy: &buildapi.DockerBuildStrategy{},
			},
			Output: buildapi.BuildOutput{
				To: &kapi.ObjectReference{
					Kind: "DockerImage",
					Name: "repository/data",
				},
			},
		},
		Status: buildapi.BuildStatus{
			Phase: buildapi.BuildPhaseNew,
		},
	}
	if result := ValidateBuild(build); len(result) != 2 {
		t.Errorf("Unexpected validation result: %v", result)
	}
}

func newDefaultParameters() buildapi.BuildSpec {
	return buildapi.BuildSpec{
		Source: buildapi.BuildSource{
			Git: &buildapi.GitBuildSource{
				URI: "http://github.com/my/repository",
			},
			ContextDir: "context",
		},
		Strategy: buildapi.BuildStrategy{
			DockerStrategy: &buildapi.DockerBuildStrategy{},
		},
		Output: buildapi.BuildOutput{
			To: &kapi.ObjectReference{
				Kind: "DockerImage",
				Name: "repository/data",
			},
		},
	}
}

func newNonDefaultParameters() buildapi.BuildSpec {
	o := newDefaultParameters()
	o.Source.Git.URI = "changed"
	return o
}

func TestValidateBuildUpdate(t *testing.T) {
	old := &buildapi.Build{
		ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
		Spec:       newDefaultParameters(),
		Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseRunning},
	}

	errs := ValidateBuildUpdate(
		&buildapi.Build{
			ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
			Spec:       newDefaultParameters(),
			Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseComplete},
		},
		old,
	)
	if len(errs) != 0 {
		t.Errorf("expected success: %v", errs)
	}

	errorCases := map[string]struct {
		Old    *buildapi.Build
		Update *buildapi.Build
		T      field.ErrorType
		F      string
	}{
		"changed spec": {
			Old: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
			},
			Update: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newNonDefaultParameters(),
			},
			T: field.ErrorTypeInvalid,
			F: "spec",
		},
		"update from terminal1": {
			Old: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseComplete},
			},
			Update: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseRunning},
			},
			T: field.ErrorTypeInvalid,
			F: "status.phase",
		},
		"update from terminal2": {
			Old: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseCancelled},
			},
			Update: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseRunning},
			},
			T: field.ErrorTypeInvalid,
			F: "status.phase",
		},
		"update from terminal3": {
			Old: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseError},
			},
			Update: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseRunning},
			},
			T: field.ErrorTypeInvalid,
			F: "status.phase",
		},
		"update from terminal4": {
			Old: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseFailed},
			},
			Update: &buildapi.Build{
				ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault, Name: "my-build", ResourceVersion: "1"},
				Spec:       newDefaultParameters(),
				Status:     buildapi.BuildStatus{Phase: buildapi.BuildPhaseRunning},
			},
			T: field.ErrorTypeInvalid,
			F: "status.phase",
		},
	}

	for k, v := range errorCases {
		errs := ValidateBuildUpdate(v.Update, v.Old)
		if len(errs) == 0 {
			t.Errorf("expected failure %s for %v", k, v.Update)
			continue
		}
		for i := range errs {
			if errs[i].Type != v.T {
				t.Errorf("%s: expected errors to have type %s: %v", k, v.T, errs[i])
			}
			if errs[i].Field != v.F {
				t.Errorf("%s: expected errors to have field %s: %v", k, v.F, errs[i])
			}
		}
	}
}

func TestBuildConfigGitSourceWithProxyFailure(t *testing.T) {
	proxyAddress := "127.0.0.1:3128"
	buildConfig := &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{Name: "config-id", Namespace: "namespace"},
		Spec: buildapi.BuildConfigSpec{
			BuildSpec: buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI:        "git://github.com/my/repository",
						HTTPProxy:  &proxyAddress,
						HTTPSProxy: &proxyAddress,
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
	}
	errors := ValidateBuildConfig(buildConfig)
	if len(errors) != 1 {
		t.Errorf("Expected one error, got %d", len(errors))
	}
	err := errors[0]
	if err.Type != field.ErrorTypeInvalid {
		t.Errorf("Expected invalid value validation error, got %q", err.Type)
	}
	if err.Detail != "only http:// and https:// GIT protocols are allowed with HTTP or HTTPS proxy set" {
		t.Errorf("Exptected git:// protocol with proxy validation error, got: %q", err.Detail)
	}
}

// TestBuildConfigDockerStrategyImageChangeTrigger ensures that it is invalid to
// have a BuildConfig with Docker strategy and an ImageChangeTrigger where
// neither DockerStrategy.From nor ImageChange.From are defined.
func TestBuildConfigDockerStrategyImageChangeTrigger(t *testing.T) {
	buildConfig := &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{Name: "config-id", Namespace: "namespace"},
		Spec: buildapi.BuildConfigSpec{
			BuildSpec: buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
			Triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
			},
		},
	}
	errors := ValidateBuildConfig(buildConfig)
	switch len(errors) {
	case 0:
		t.Errorf("Expected validation error, got nothing")
	case 1:
		err := errors[0]
		if err.Type != field.ErrorTypeRequired {
			t.Errorf("Expected error to be '%v', got '%v'", field.ErrorTypeRequired, err.Type)
		}
	default:
		t.Errorf("Expected a single validation error, got %v", errors)
	}
}

func TestBuildConfigValidationFailureRequiredName(t *testing.T) {
	buildConfig := &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{Name: "", Namespace: "foo"},
		Spec: buildapi.BuildConfigSpec{
			BuildSpec: buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
	}
	errors := ValidateBuildConfig(buildConfig)
	if len(errors) != 1 {
		t.Fatalf("Unexpected validation errors %v", errors)
	}
	err := errors[0]
	if err.Type != field.ErrorTypeRequired {
		t.Errorf("Unexpected error type, expected %s, got %s", field.ErrorTypeRequired, err.Type)
	}
	if err.Field != "metadata.name" {
		t.Errorf("Unexpected field name expected metadata.name, got %s", err.Field)
	}
}

func TestBuildConfigImageChangeTriggers(t *testing.T) {
	tests := []struct {
		name        string
		triggers    []buildapi.BuildTriggerPolicy
		expectError bool
		errorType   field.ErrorType
	}{
		{
			name: "valid default trigger",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
			},
			expectError: false,
		},
		{
			name: "more than one default trigger",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
		{
			name: "missing image change struct",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type: buildapi.ImageChangeBuildTriggerType,
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeRequired,
		},
		{
			name: "only one default image change trigger",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "myimage:tag",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "invalid reference kind for trigger",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind: "DockerImage",
							Name: "myimage:tag",
						},
					},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
		{
			name: "empty reference kind for trigger",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Name: "myimage:tag",
						},
					},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
		{
			name: "duplicate imagestreamtag references",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "myimage:tag",
						},
					},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "myimage:tag",
						},
					},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
		{
			name: "duplicate imagestreamtag - same as strategy ref",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type:        buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "builderimage:latest",
						},
					},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
		{
			name: "imagestreamtag references with same name, different ns",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      "myimage:tag",
							Namespace: "ns1",
						},
					},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      "myimage:tag",
							Namespace: "ns2",
						},
					},
				},
			},
			expectError: false,
		},
		{
			name: "imagestreamtag references with same name, same ns",
			triggers: []buildapi.BuildTriggerPolicy{
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      "myimage:tag",
							Namespace: "ns",
						},
					},
				},
				{
					Type: buildapi.ImageChangeBuildTriggerType,
					ImageChange: &buildapi.ImageChangeTrigger{
						From: &kapi.ObjectReference{
							Kind:      "ImageStreamTag",
							Name:      "myimage:tag",
							Namespace: "ns",
						},
					},
				},
			},
			expectError: true,
			errorType:   field.ErrorTypeInvalid,
		},
	}

	for _, tc := range tests {
		buildConfig := &buildapi.BuildConfig{
			ObjectMeta: kapi.ObjectMeta{Name: "bar", Namespace: "foo"},
			Spec: buildapi.BuildConfigSpec{
				BuildSpec: buildapi.BuildSpec{
					Source: buildapi.BuildSource{
						Git: &buildapi.GitBuildSource{
							URI: "http://github.com/my/repository",
						},
						ContextDir: "context",
					},
					Strategy: buildapi.BuildStrategy{
						SourceStrategy: &buildapi.SourceBuildStrategy{
							From: kapi.ObjectReference{
								Kind: "ImageStreamTag",
								Name: "builderimage:latest",
							},
						},
					},
					Output: buildapi.BuildOutput{
						To: &kapi.ObjectReference{
							Kind: "DockerImage",
							Name: "repository/data",
						},
					},
				},
				Triggers: tc.triggers,
			},
		}
		errors := ValidateBuildConfig(buildConfig)
		// Check whether an error was returned
		if hasError := len(errors) > 0; hasError != tc.expectError {
			t.Errorf("%s: did not get expected result: %#v", tc.name, errors)
		}
		// Check whether it's the expected error type
		if len(errors) > 0 && tc.expectError && tc.errorType != "" {
			verr := errors[0]
			if verr.Type != tc.errorType {
				t.Errorf("%s: unexpected error type. Expected: %s. Got: %s", tc.name, tc.errorType, verr.Type)
			}
		}
	}
}

func TestBuildConfigValidationOutputFailure(t *testing.T) {
	buildConfig := &buildapi.BuildConfig{
		ObjectMeta: kapi.ObjectMeta{Name: ""},
		Spec: buildapi.BuildConfigSpec{
			BuildSpec: buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Name: "other",
					},
				},
			},
		},
	}
	if result := ValidateBuildConfig(buildConfig); len(result) != 3 {
		for _, e := range result {
			t.Errorf("Unexpected validation result %v", e)
		}
	}
}

func TestValidateBuildRequest(t *testing.T) {
	testCases := map[string]*buildapi.BuildRequest{
		string(field.ErrorTypeRequired) + "metadata.namespace": {ObjectMeta: kapi.ObjectMeta{Name: "requestName"}},
		string(field.ErrorTypeRequired) + "metadata.name":      {ObjectMeta: kapi.ObjectMeta{Namespace: kapi.NamespaceDefault}},
	}

	for desc, tc := range testCases {
		errors := ValidateBuildRequest(tc)
		if len(desc) == 0 && len(errors) > 0 {
			t.Errorf("%s: Unexpected validation result: %v", desc, errors)
		}
		if len(desc) > 0 && len(errors) != 1 {
			t.Errorf("%s: Unexpected validation result: %v", desc, errors)
		}
		if len(desc) > 0 {
			err := errors[0]
			errDesc := string(err.Type) + err.Field
			if desc != errDesc {
				t.Errorf("Unexpected validation result for %s: expected %s, got %s", err.Field, desc, errDesc)
			}
		}
	}
}

func TestValidateSource(t *testing.T) {
	dockerfile := "FROM something"
	invalidProxyAddress := "some!@#$%^&*()url"
	errorCases := []struct {
		t        field.ErrorType
		path     string
		source   *buildapi.BuildSource
		ok       bool
		multiple bool
	}{
		// 0
		{
			t:    field.ErrorTypeRequired,
			path: "git.uri",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "",
				},
			},
		},
		// 1
		{
			t:    field.ErrorTypeInvalid,
			path: "git.uri",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "::",
				},
			},
		},
		// 2
		{
			t:    field.ErrorTypeInvalid,
			path: "contextDir",
			source: &buildapi.BuildSource{
				Dockerfile: &dockerfile,
				ContextDir: "../file",
			},
		},
		// 3
		{
			t:    field.ErrorTypeInvalid,
			path: "git",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "https://example.com/repo.git",
				},
				Binary: &buildapi.BinaryBuildSource{},
			},
			multiple: true,
		},
		// 4
		{
			t:    field.ErrorTypeInvalid,
			path: "binary.asFile",
			source: &buildapi.BuildSource{
				Binary: &buildapi.BinaryBuildSource{AsFile: "/a/path"},
			},
		},
		// 5
		{
			t:    field.ErrorTypeInvalid,
			path: "binary.asFile",
			source: &buildapi.BuildSource{
				Binary: &buildapi.BinaryBuildSource{AsFile: "/"},
			},
		},
		// 6
		{
			t:    field.ErrorTypeInvalid,
			path: "binary.asFile",
			source: &buildapi.BuildSource{
				Binary: &buildapi.BinaryBuildSource{AsFile: "a\\b"},
			},
		},
		// 7
		{
			source: &buildapi.BuildSource{
				Binary: &buildapi.BinaryBuildSource{AsFile: "/././file"},
			},
			ok: true,
		},
		// 8
		{
			source: &buildapi.BuildSource{
				Binary:     &buildapi.BinaryBuildSource{AsFile: "/././file"},
				Dockerfile: &dockerfile,
			},
			ok: true,
		},
		// 9
		{
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "https://example.com/repo.git",
				},
				Dockerfile: &dockerfile,
			},
			ok: true,
		},
		// 10
		{
			source: &buildapi.BuildSource{
				Dockerfile: &dockerfile,
			},
			ok: true,
		},
		// 11
		{
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "https://example.com/repo.git",
				},
				ContextDir: "contextDir",
			},
			ok: true,
		},
		// 12
		{
			t:    field.ErrorTypeRequired,
			path: "sourceSecret.name",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI: "http://example.com/repo.git",
				},
				SourceSecret: &kapi.LocalObjectReference{},
				ContextDir:   "contextDir/../somedir",
			},
		},
		// 13
		{
			t:    field.ErrorTypeInvalid,
			path: "git.httpproxy",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI:       "https://example.com/repo.git",
					HTTPProxy: &invalidProxyAddress,
				},
				ContextDir: "contextDir",
			},
		},
		// 14
		{
			t:    field.ErrorTypeInvalid,
			path: "git.httpsproxy",
			source: &buildapi.BuildSource{
				Git: &buildapi.GitBuildSource{
					URI:        "https://example.com/repo.git",
					HTTPSProxy: &invalidProxyAddress,
				},
				ContextDir: "contextDir",
			},
		},
		// 15
		{
			ok: true,
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "test/dir",
							},
						},
					},
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "test/dir",
							},
						},
					},
				},
			},
		},
		// 16
		{
			t:    field.ErrorTypeRequired,
			path: "images[0].paths",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
					},
				},
			},
		},
		// 17 - destinationdir is not relative.
		{
			t:    field.ErrorTypeInvalid,
			path: "images[0].paths[0].destinationDir",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "/test/dir",
							},
						},
					},
				},
			},
		},
		// 18 - sourcepath is not absolute.
		{
			t:    field.ErrorTypeInvalid,
			path: "images[0].paths[0].sourcePath",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "some/path",
								DestinationDir: "test/dir",
							},
						},
					},
				},
			},
		},
		// 19 - destinationdir backsteps above basedir
		{
			t:    field.ErrorTypeInvalid,
			path: "images[0].paths[0].destinationDir",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "test/../../dir",
							},
						},
					},
				},
			},
		},
		// 20
		{
			t:    field.ErrorTypeInvalid,
			path: "images[0].from.kind",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "InvalidKind",
							Name: "my-image:latest",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "test/dir",
							},
						},
					},
				},
			},
		},
		// 21
		{
			t:    field.ErrorTypeRequired,
			path: "images[0].pullSecret.name",
			source: &buildapi.BuildSource{
				Images: []buildapi.ImageSource{
					{
						From: kapi.ObjectReference{
							Kind: "DockerImage",
							Name: "my-image:latest",
						},
						PullSecret: &kapi.LocalObjectReference{
							Name: "",
						},
						Paths: []buildapi.ImageSourcePath{
							{
								SourcePath:     "/some/path",
								DestinationDir: "test/dir",
							},
						},
					},
				},
			},
		},
	}
	for i, tc := range errorCases {
		errors := validateSource(tc.source, false, false, nil)
		switch len(errors) {
		case 0:
			if !tc.ok {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
			}
			continue
		case 1:
			if tc.ok || tc.multiple {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
				continue
			}
		default:
			if tc.ok || !tc.multiple {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
				continue
			}
		}
		err := errors[0]
		if err.Type != tc.t {
			t.Errorf("%d: Expected error type %s, got %s", i, tc.t, err.Type)
		}
		if err.Field != tc.path {
			t.Errorf("%d: Expected error path %s, got %s", i, tc.path, err.Field)
		}
	}

	errorCases[11].source.ContextDir = "."
	validateSource(errorCases[11].source, false, false, nil)
	if len(errorCases[11].source.ContextDir) != 0 {
		t.Errorf("ContextDir was not cleaned: %s", errorCases[11].source.ContextDir)
	}
}

func TestValidateStrategy(t *testing.T) {
	errorCases := []struct {
		t        field.ErrorType
		path     string
		strategy *buildapi.BuildStrategy
		ok       bool
		multiple bool
	}{
		// 0
		{
			t:    field.ErrorTypeInvalid,
			path: "",
			strategy: &buildapi.BuildStrategy{
				SourceStrategy: &buildapi.SourceBuildStrategy{},
				DockerStrategy: &buildapi.DockerBuildStrategy{},
				CustomStrategy: &buildapi.CustomBuildStrategy{},
			},
		},
	}
	for i, tc := range errorCases {
		errors := validateStrategy(tc.strategy, nil)
		switch len(errors) {
		case 0:
			if !tc.ok {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
			}
			continue
		case 1:
			if tc.ok || tc.multiple {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
				continue
			}
		default:
			if tc.ok || !tc.multiple {
				t.Errorf("%d: Unexpected validation result: %v", i, errors)
				continue
			}
		}
		err := errors[0]
		if err.Type != tc.t {
			t.Errorf("%d: Unexpected error type: %s", i, err.Type)
		}
		if err.Field != tc.path {
			t.Errorf("%d: Unexpected error path: %s", i, err.Field)
		}
	}
}

func TestValidateBuildSpec(t *testing.T) {
	zero := int64(0)
	longString := strings.Repeat("1234567890", 100*61)
	//shortString := "FROM foo"
	errorCases := []struct {
		err string
		*buildapi.BuildSpec
	}{
		// 0
		{
			string(field.ErrorTypeInvalid) + "output.to.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "some/long/value/with/no/meaning",
					},
				},
			},
		},
		// 1
		{
			string(field.ErrorTypeInvalid) + "output.to.kind",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "Foo",
						Name: "test",
					},
				},
			},
		},
		// 2
		{
			string(field.ErrorTypeRequired) + "output.to.kind",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{},
				},
			},
		},
		// 3
		{
			string(field.ErrorTypeRequired) + "output.to.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "ImageStreamTag",
					},
				},
			},
		},
		// 4
		{
			string(field.ErrorTypeInvalid) + "output.to.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind:      "ImageStreamTag",
						Name:      "missingtag",
						Namespace: "subdomain",
					},
				},
			},
		},
		// 5
		{
			string(field.ErrorTypeInvalid) + "output.to.namespace",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind:      "ImageStreamTag",
						Name:      "test:tag",
						Namespace: "not_a_valid_subdomain",
					},
				},
			},
		},
		// 6
		// invalid because from is not specified in the
		// sti strategy definition
		{
			string(field.ErrorTypeRequired) + "strategy.sourceStrategy.from.kind",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					SourceStrategy: &buildapi.SourceBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 7
		// Invalid because from.name is not specified
		{
			string(field.ErrorTypeRequired) + "strategy.sourceStrategy.from.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					SourceStrategy: &buildapi.SourceBuildStrategy{
						From: kapi.ObjectReference{
							Kind: "DockerImage",
						},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 8
		// invalid because from name is a bad format
		{
			string(field.ErrorTypeInvalid) + "strategy.sourceStrategy.from.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					SourceStrategy: &buildapi.SourceBuildStrategy{
						From: kapi.ObjectReference{Kind: "ImageStreamTag", Name: "bad format"},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 9
		// invalid because from is not specified in the
		// custom strategy definition
		{
			string(field.ErrorTypeRequired) + "strategy.customStrategy.from.kind",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					CustomStrategy: &buildapi.CustomBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 10
		// invalid because from.name is not specified in the
		// custom strategy definition
		{
			string(field.ErrorTypeInvalid) + "strategy.customStrategy.from.name",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					CustomStrategy: &buildapi.CustomBuildStrategy{
						From: kapi.ObjectReference{Kind: "ImageStreamTag", Name: "bad format"},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 11
		{
			string(field.ErrorTypeInvalid) + "source.dockerfile",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Dockerfile: &longString,
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
			},
		},
		// 12
		{
			string(field.ErrorTypeInvalid) + "source.dockerfile",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Dockerfile: &longString,
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
			},
		},
		// 13
		// invalid because CompletionDeadlineSeconds <= 0
		{
			string(field.ErrorTypeInvalid) + "completionDeadlineSeconds",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
				CompletionDeadlineSeconds: &zero,
			},
		},
		// 14
		// must provide some source input
		{
			string(field.ErrorTypeInvalid) + "source",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 15
		// dockerfilePath can't be an absolute path
		{
			string(field.ErrorTypeInvalid) + "strategy.dockerStrategy.dockerfilePath",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{
						DockerfilePath: "/myDockerfile",
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 16
		// dockerfilePath can't start with ..
		{
			string(field.ErrorTypeInvalid) + "strategy.dockerStrategy.dockerfilePath",
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
					ContextDir: "context",
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{
						DockerfilePath: "../someDockerfile",
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		}}

	for count, config := range errorCases {
		errors := validateBuildSpec(config.BuildSpec, nil)
		if len(errors) != 1 {
			t.Errorf("Test[%d] %s: Unexpected validation result: %v", count, config.err, errors)
			continue
		}
		err := errors[0]
		errDesc := string(err.Type) + err.Field
		if config.err != errDesc {
			t.Errorf("Test[%d] Unexpected validation result for %s: expected %s, got %s", count, err.Field, config.err, errDesc)
		}
	}
}

func TestValidateBuildSpecSuccess(t *testing.T) {
	shortString := "FROM foo"
	testCases := []struct {
		*buildapi.BuildSpec
	}{
		// 0
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					SourceStrategy: &buildapi.SourceBuildStrategy{
						From: kapi.ObjectReference{
							Kind: "DockerImage",
							Name: "reponame",
						},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 1
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					CustomStrategy: &buildapi.CustomBuildStrategy{
						From: kapi.ObjectReference{
							Kind: "ImageStreamTag",
							Name: "imagestreamname:tag",
						},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 2
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 3
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamImage",
							Name: "imagestreamimage",
						},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 4
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Dockerfile: &shortString,
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamImage",
							Name: "imagestreamimage",
						},
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
		// 5
		{
			&buildapi.BuildSpec{
				Source: buildapi.BuildSource{
					Git: &buildapi.GitBuildSource{
						URI: "http://github.com/my/repository",
					},
				},
				Strategy: buildapi.BuildStrategy{
					DockerStrategy: &buildapi.DockerBuildStrategy{
						From: &kapi.ObjectReference{
							Kind: "ImageStreamImage",
							Name: "imagestreamimage",
						},
						DockerfilePath: "dockerfiles/firstDockerfile",
					},
				},
				Output: buildapi.BuildOutput{
					To: &kapi.ObjectReference{
						Kind: "DockerImage",
						Name: "repository/data",
					},
				},
			},
		},
	}

	for count, config := range testCases {
		errors := validateBuildSpec(config.BuildSpec, nil)
		if len(errors) != 0 {
			t.Errorf("Test[%d] Unexpected validation error: %v", count, errors)
		}
	}

}

func TestValidateDockerfilePath(t *testing.T) {
	tests := []struct {
		strategy               *buildapi.DockerBuildStrategy
		expectedDockerfilePath string
	}{
		{
			strategy: &buildapi.DockerBuildStrategy{
				DockerfilePath: ".",
			},
			expectedDockerfilePath: "",
		},
		{
			strategy: &buildapi.DockerBuildStrategy{
				DockerfilePath: "somedir/..",
			},
			expectedDockerfilePath: "",
		},
		{
			strategy: &buildapi.DockerBuildStrategy{
				DockerfilePath: "somedir/../somedockerfile",
			},
			expectedDockerfilePath: "somedockerfile",
		},
		{
			strategy: &buildapi.DockerBuildStrategy{
				DockerfilePath: "somedir/somedockerfile",
			},
			expectedDockerfilePath: "somedir/somedockerfile",
		},
	}

	for count, test := range tests {
		errors := validateDockerStrategy(test.strategy, nil)
		if len(errors) != 0 {
			t.Errorf("Test[%d] Unexpected validation error: %v", count, errors)
		}
		if test.strategy.DockerfilePath != test.expectedDockerfilePath {
			t.Errorf("Test[%d] Unexpected DockerfilePath: %v (expected: %s)", count, test.strategy.DockerfilePath, test.expectedDockerfilePath)
		}
	}
}

func TestValidateTrigger(t *testing.T) {
	tests := map[string]struct {
		trigger  buildapi.BuildTriggerPolicy
		expected []*field.Error
	}{
		"trigger without type": {
			trigger:  buildapi.BuildTriggerPolicy{},
			expected: []*field.Error{field.Required(field.NewPath("type"))},
		},
		"trigger with unknown type": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: "UnknownTriggerType",
			},
			expected: []*field.Error{field.Invalid(field.NewPath("type"), "", "")},
		},
		"GitHub type with no github webhook": {
			trigger:  buildapi.BuildTriggerPolicy{Type: buildapi.GitHubWebHookBuildTriggerType},
			expected: []*field.Error{field.Required(field.NewPath("github"))},
		},
		"GitHub trigger with no secret": {
			trigger: buildapi.BuildTriggerPolicy{
				Type:          buildapi.GitHubWebHookBuildTriggerType,
				GitHubWebHook: &buildapi.WebHookTrigger{},
			},
			expected: []*field.Error{field.Required(field.NewPath("github", "secret"))},
		},
		"GitHub trigger with generic webhook": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.GitHubWebHookBuildTriggerType,
				GenericWebHook: &buildapi.WebHookTrigger{
					Secret: "secret101",
				},
			},
			expected: []*field.Error{field.Required(field.NewPath("github"))},
		},
		"Generic trigger with no generic webhook": {
			trigger:  buildapi.BuildTriggerPolicy{Type: buildapi.GenericWebHookBuildTriggerType},
			expected: []*field.Error{field.Required(field.NewPath("generic"))},
		},
		"Generic trigger with no secret": {
			trigger: buildapi.BuildTriggerPolicy{
				Type:           buildapi.GenericWebHookBuildTriggerType,
				GenericWebHook: &buildapi.WebHookTrigger{},
			},
			expected: []*field.Error{field.Required(field.NewPath("generic", "secret"))},
		},
		"Generic trigger with github webhook": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.GenericWebHookBuildTriggerType,
				GitHubWebHook: &buildapi.WebHookTrigger{
					Secret: "secret101",
				},
			},
			expected: []*field.Error{field.Required(field.NewPath("generic"))},
		},
		"ImageChange trigger without params": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.ImageChangeBuildTriggerType,
			},
			expected: []*field.Error{field.Required(field.NewPath("imageChange"))},
		},
		"valid GitHub trigger": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.GitHubWebHookBuildTriggerType,
				GitHubWebHook: &buildapi.WebHookTrigger{
					Secret: "secret101",
				},
			},
		},
		"valid Generic trigger": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.GenericWebHookBuildTriggerType,
				GenericWebHook: &buildapi.WebHookTrigger{
					Secret: "secret101",
				},
			},
		},
		"valid ImageChange trigger": {
			trigger: buildapi.BuildTriggerPolicy{
				Type: buildapi.ImageChangeBuildTriggerType,
				ImageChange: &buildapi.ImageChangeTrigger{
					LastTriggeredImageID: "asdf1234",
				},
			},
		},
		"valid ImageChange trigger with empty fields": {
			trigger: buildapi.BuildTriggerPolicy{
				Type:        buildapi.ImageChangeBuildTriggerType,
				ImageChange: &buildapi.ImageChangeTrigger{},
			},
		},
	}
	for desc, test := range tests {
		errors := validateTrigger(&test.trigger, nil)
		if len(test.expected) == 0 {
			if len(errors) != 0 {
				t.Errorf("%s: Got unexpected validation errors: %#v", desc, errors)
			}
			continue
		}
		if len(errors) != 1 {
			t.Errorf("%s: Expected one validation error, got %d", desc, len(errors))
			for i, err := range errors {
				validationError := err
				t.Errorf("  %d. %v", i+1, validationError)
			}
			continue
		}
		err := errors[0]
		validationError := err
		if validationError.Type != test.expected[0].Type {
			t.Errorf("%s: Expected error type %s, got %s", desc, test.expected[0].Type, validationError.Type)
		}
		if validationError.Field != test.expected[0].Field {
			t.Errorf("%s: Expected error field %s, got %s", desc, test.expected[0].Field, validationError.Field)
		}
	}
}

func TestValidateToImageReference(t *testing.T) {
	o := &kapi.ObjectReference{
		Name:      "somename",
		Namespace: "somenamespace",
		Kind:      "DockerImage",
	}
	errs := validateToImageReference(o, nil)
	if len(errs) != 1 {
		t.Errorf("Wrong number of errors: %v", errs)
	}
	err := errs[0]
	if err.Type != field.ErrorTypeInvalid {
		t.Errorf("Wrong error type, expected %v, got %v", field.ErrorTypeInvalid, err.Type)
	}
	if err.Field != "namespace" {
		t.Errorf("Error on wrong field, expected %s, got %s", "namespace", err.Field)
	}
}

func TestValidateStrategyEnvVars(t *testing.T) {
	tests := []struct {
		env         []kapi.EnvVar
		errExpected bool
		errField    string
		errType     field.ErrorType
	}{
		// 0: missing Env variable name
		{
			env: []kapi.EnvVar{
				{
					Name:  "",
					Value: "test",
				},
			},
			errExpected: true,
			errField:    "env[0].name",
			errType:     field.ErrorTypeRequired,
		},
		// 1: invalid Env variable name
		{
			env: []kapi.EnvVar{
				{
					Name:  " invalid,name",
					Value: "test",
				},
			},
			errExpected: true,
			errField:    "env[0].name",
			errType:     field.ErrorTypeInvalid,
		},
		// 2: valueFrom present in env var
		{
			env: []kapi.EnvVar{
				{
					Name:      "name",
					Value:     "test",
					ValueFrom: &kapi.EnvVarSource{},
				},
			},
			errExpected: true,
			errField:    "env[0].valueFrom",
			errType:     field.ErrorTypeInvalid,
		},
		// 3: valid env
		{
			env: []kapi.EnvVar{
				{
					Name:  "VAR1",
					Value: "value1",
				},
				{
					Name:  "VAR2",
					Value: "value2",
				},
			},
		},
	}

	for i, tc := range tests {
		errs := ValidateStrategyEnv(tc.env, field.NewPath("env"))
		if !tc.errExpected {
			if len(errs) > 0 {
				t.Errorf("%d: unexpected error: %v", i, errs.ToAggregate())
			}
			continue
		}
		if tc.errExpected && len(errs) == 0 {
			t.Errorf("%d: expected error. Got none.", i)
			continue
		}
		err := errs[0]
		if err.Field != tc.errField {
			t.Errorf("%d: unexpected error field: %s", i, err.Field)
		}
		if err.Type != tc.errType {
			t.Errorf("%d: unexpected error type: %s", i, err.Type)
		}
	}
}
