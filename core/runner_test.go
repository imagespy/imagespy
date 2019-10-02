package core

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/imagespy/imagespy/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustExec(t *testing.T, name string, args ...string) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("command '%s %s' failed: %s", name, strings.Join(args, " "), out)
		t.Fatal(err)
	}
}

func TestRunner_Run(t *testing.T) {
	testcases := []struct {
		name            string
		discovery       *discovery.Input
		expectedMetrics []string
	}{
		{
			name: "When a newer version of an image is available it exports the status of the image as needs-update",
			discovery: &discovery.Input{
				Name: "test",
				Images: []*discovery.Image{
					{Digest: "sha256:55f250f8bc296f15478819abd7439a70c08f9864ad2fde20be55a39341e58c93", Repository: "127.0.0.1:52854/redis", Source: "ttt", Tag: "4.0.14-alpine"},
				},
			},
			expectedMetrics: []string{
				"tinker_up 1",
				`tinker_image_status{current_digest="55f250f8",current_tag="4.0.14-alpine",input="test",latest_digest="e1cd649a",latest_tag="5.0.6-alpine",repository="127.0.0.1:52854/redis",source="ttt"} 1`,
			},
		},
		{
			name: "When the image is the latest version it exports the status of the image as no-update",
			discovery: &discovery.Input{
				Name: "test",
				Images: []*discovery.Image{
					{Digest: "sha256:e1cd649ac85b0b170d70ce695644999419764621de5208f0fb00283aef0fdc2f", Repository: "127.0.0.1:52854/redis", Source: "ttt", Tag: "5.0.6-alpine"},
				},
			},
			expectedMetrics: []string{
				"tinker_up 1",
				`tinker_image_status{current_digest="e1cd649a",current_tag="5.0.6-alpine",input="test",latest_digest="e1cd649a",latest_tag="5.0.6-alpine",repository="127.0.0.1:52854/redis",source="ttt"} 0`,
			},
		},
	}

	tmpDir, err := ioutil.TempDir("", "testrunner_run")
	require.NoError(t, err, "create temp dir")

	r, err := NewRunnerFromConfig("fixtures/TestRunner_Run/config.yaml")
	require.NoError(t, err, "create runner from config")
	r.cfg.DiscoveryDirectory = tmpDir
	go func() { r.Run() }()
	defer r.Stop()

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			disc, _ := json.Marshal(tc.discovery)
			ioutil.WriteFile(path.Join(tmpDir, "test.json"), disc, 0644)

			resp, err := http.Get("http://127.0.0.1:8567/metrics")
			require.NoError(t, err, "http get metrics")
			b, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err, "read http response body")
			defer resp.Body.Close()

			for _, em := range tc.expectedMetrics {
				assert.Regexp(t, regexp.MustCompile(em), string(b))
			}
		})
	}
}
