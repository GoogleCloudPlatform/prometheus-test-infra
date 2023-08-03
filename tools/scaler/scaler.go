// Copyright 2019 The Prometheus Authors
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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
	appsV1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/prometheus/test-infra/pkg/provider/k8s"
)

type scale struct {
	k8sClient     *k8s.K8s
	min           int32
	max           int32
	interval      time.Duration
	patternName   string
	scalingFactor int32 // for step-like scaling (~ step height)
}

func newScaler() *scale {
	k, err := k8s.New(context.Background(), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error creating k8s client inside the k8s cluster"))
		os.Exit(2)
	}
	return &scale{
		k8sClient: k,
	}
}

func (s *scale) updateReplicas(replicas *int32) []k8s.Resource {
	var k8sResource []k8s.Resource
	for _, deployment := range s.k8sClient.GetResources() {
		k8sObjects := make([]runtime.Object, 0)

		for _, resource := range deployment.Objects {
			if kind := strings.ToLower(resource.GetObjectKind().GroupVersionKind().Kind); kind == "deployment" {
				req := resource.(*appsV1.Deployment)
				req.Spec.Replicas = replicas
				k8sObjects = append(k8sObjects, req.DeepCopyObject())
			}
		}
		if len(k8sObjects) > 0 {
			k8sResource = append(k8sResource, k8s.Resource{FileName: deployment.FileName, Objects: k8sObjects})
		}
	}
	return k8sResource
}

func (s *scale) scale(*kingpin.ParseContext) error {

	if s.patternName == "burst" {
		log.Printf("Auto-scale pattern: %s", s.patternName)
		log.Printf("Starting Prombench-Scaler:\n\t max: %d\n\t min: %d\n\t interval: %s", s.max, s.min, s.interval)

		maxResourceObjects := s.updateReplicas(&s.max)
		minResourceObjects := s.updateReplicas(&s.min)

		for {
			log.Printf("Scaling Deployment to %d", s.max)
			if err := s.k8sClient.ResourceApply(maxResourceObjects); err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error scaling deployment"))
			}

			time.Sleep(s.interval)

			log.Printf("Scaling Deployment to %d", s.min)
			if err := s.k8sClient.ResourceApply(minResourceObjects); err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error scaling deployment"))
			}

			time.Sleep(s.interval)
		}

	} else if s.patternName == "step" {
		log.Printf("Auto-scale pattern: %s", s.patternName)

		updateScalingFactor := false

		if s.scalingFactor >= s.max {
			log.Printf("scalingFactor (%d) >= max (%d)", s.scalingFactor, s.max)
			updateScalingFactor = true
		}
		if s.scalingFactor == 0 {
			log.Print("scalingFactor is set to 0.")
			updateScalingFactor = true
		}

		if updateScalingFactor == true {
			s.scalingFactor = int32(s.max / 10) // 10 steps
			log.Printf("Updating the scaling factor to: %d", s.scalingFactor)
		}

		log.Printf("Starting Prombench-Scaler:\n\t max: %d\n\t min: %d\n\t interval: %s\n\t scalingFactor: %d", s.max, s.min, s.interval, s.scalingFactor)

		numberOfResources := s.min
		for {
			log.Printf("Scaling Deployment to %d", numberOfResources)

			resourceObjects := s.updateReplicas(&numberOfResources)
			if err := s.k8sClient.ResourceApply(resourceObjects); err != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error scaling deployment"))
			}

			time.Sleep(s.interval)

			numberOfResources += s.scalingFactor
			if numberOfResources > s.max {
				numberOfResources = s.max
			}
		}

	} else {
		log.Printf("Invalid pattern: %s", s.patternName)
		os.Exit(2)
	}

	return nil
}

func main() {

	app := kingpin.New(filepath.Base(os.Args[0]), "The Prombench-Scaler tool")
	app.HelpFlag.Short('h')

	s := newScaler()

	k8sApp := app.Command("scale", "Scale a Kubernetes deployment object periodically up and down. \nex: ./scaler scale -v NAMESPACE:scale -f fake-webserver.yaml 20 1 15m").
		Action(s.k8sClient.DeploymentsParse).
		Action(s.scale)
	k8sApp.Flag("file", "yaml file or folder that describes the parameters for the deployment.").
		Required().
		Short('f').
		ExistingFilesOrDirsVar(&s.k8sClient.DeploymentFiles)
	k8sApp.Flag("vars", "When provided it will substitute the token holders in the yaml file. Follows the standard golang template formating - {{ .hashStable }}.").
		Short('v').
		StringMapVar(&s.k8sClient.DeploymentVars)
	k8sApp.Arg("max", "Number of Replicas to scale up.").
		Required().
		Int32Var(&s.max)
	k8sApp.Arg("min", "Number of Replicas to scale down.").
		Required().
		Int32Var(&s.min)
	k8sApp.Arg("interval", "Time to wait before changing the number of replicas.").
		Required().
		DurationVar(&s.interval)
	k8sApp.Arg("patternName", "Auto-scaling pattern. Defines the scaling function for K8s resources. Available values: burst, step.").
		Required().
		Default("burst").
		StringVar(&s.patternName)
	k8sApp.Arg("scalingFactor", "Indicates the 'step height' during step-like autoscaling.").
		Int32Var(&s.scalingFactor)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	}
}
