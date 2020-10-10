/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package octo

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/paketo-buildpacks/pipeline-builder/octo/actions"
	"github.com/paketo-buildpacks/pipeline-builder/octo/actions/event"
	"github.com/paketo-buildpacks/pipeline-builder/octo/internal"
)

func ContributeTest(descriptor Descriptor) (*Contribution, error) {
	if descriptor.OfflinePackages != nil {
		return nil, nil
	}

	w := actions.Workflow{
		Name: "Tests",
		On: map[event.Type]event.Event{
			event.PullRequestType: event.PullRequest{},
			event.PushType: event.Push{
				Branches: []string{"main"},
			},
		},
		Jobs: map[string]actions.Job{},
	}

	if f, err := internal.Find(descriptor.Path, regexp.MustCompile(`.+\.go`).MatchString); err != nil {
		return nil, fmt.Errorf("unable to find .go files in %s\n%w", descriptor.Path, err)
	} else if len(f) > 0 {
		w.Jobs["unit"] = actions.Job{
			Name:   "Unit Test",
			RunsOn: []actions.VirtualEnvironment{actions.UbuntuLatest},
			Steps: []actions.Step{
				{
					Uses: "actions/checkout@v2",
				},
				{
					Uses: "actions/cache@v2",
					With: map[string]interface{}{
						"path":         "${{ env.HOME }}/go/pkg/mod",
						"key":          "${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}",
						"restore-keys": "${{ runner.os }}-go-",
					},
				},
				{
					Uses: "actions/setup-go@v2",
					With: map[string]interface{}{"go-version": GoVersion},
				},
				{
					Name: "Install richgo",
					Run:  internal.StatikString("/install-richgo.sh"),
				},
				{
					Name: "Run Tests",
					Run:  internal.StatikString("/run-tests.sh"),
					Env:  map[string]string{"RICHGO_FORCE_COLOR": "1"},
				},
			},
		}
	}

	if descriptor.Builder != nil {
		j := actions.Job{
			Name:   "Create Builder Test",
			RunsOn: []actions.VirtualEnvironment{actions.UbuntuLatest},
			Steps: []actions.Step{
				{
					Uses: "actions/checkout@v2",
				},
				{
					Name: "Install pack",
					Run:  internal.StatikString("/install-pack.sh"),
					Env:  map[string]string{"PACK_VERSION": PackVersion},
				},
				{
					Id:   "version",
					Name: "Compute Version",
					Run:  internal.StatikString("/compute-version.sh"),
				},
				{
					Name: "Create Builder",
					Run:  internal.StatikString("/create-builder.sh"),
					Env: map[string]string{
						"BUILDER": "test",
						"VERSION": "${{ steps.version.outputs.version }}",
					},
				},
			},
		}

		j.Steps = append(NewDockerLoginActions(descriptor.Credentials), j.Steps...)

		w.Jobs["create-builder"] = j
	}

	if descriptor.Package != nil {
		j := actions.Job{
			Name:   "Create Package Test",
			RunsOn: []actions.VirtualEnvironment{actions.UbuntuLatest},
			Steps: []actions.Step{
				{
					Uses: "actions/checkout@v2",
				},
				{
					Uses: "actions/cache@v2",
					With: map[string]interface{}{
						"path": strings.Join([]string{
							"${{ env.HOME }}/.pack",
							"${{ env.HOME }}/carton-cache",
						}, "\n"),
						"key":          "${{ runner.os }}-go-${{ hashFiles('**/buildpack.toml', '**/package.toml') }}",
						"restore-keys": "${{ runner.os }}-go-",
					},
				},
				{
					Uses: "actions/setup-go@v2",
					With: map[string]interface{}{"go-version": GoVersion},
				},
				{
					Name: "Install create-package",
					Run:  internal.StatikString("/install-create-package.sh"),
				},
				{
					Name: "Install pack",
					Run:  internal.StatikString("/install-pack.sh"),
					Env:  map[string]string{"PACK_VERSION": PackVersion},
				},
				{
					Id:   "version",
					Name: "Compute Version",
					Run:  internal.StatikString("/compute-version.sh"),
				},
				{
					Name: "Create Package",
					Run:  internal.StatikString("/create-package.sh"),
					Env: map[string]string{
						"INCLUDE_DEPENDENCIES": "true",
						"VERSION":              "${{ steps.version.outputs.version }}",
					},
				},
				{
					Name: "Package Buildpack",
					Run:  internal.StatikString("/package-buildpack.sh"),
					Env: map[string]string{
						"PACKAGE": "test",
						"VERSION": "${{ steps.version.outputs.version }}",
					},
				},
			},
		}

		j.Steps = append(NewDockerLoginActions(descriptor.Credentials), j.Steps...)

		w.Jobs["create-package"] = j
	}

	c, err := NewActionContribution(w)
	if err != nil {
		return nil, err
	}

	return &c, err
}