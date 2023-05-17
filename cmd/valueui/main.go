// Copyright 2017 The hchart Authors, All rights reserved.
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

package main

import (
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hooto/hlog4g/hlog"
	"github.com/hooto/htoml4g/htoml"
	"github.com/hooto/httpsrv"

	"github.com/valuedig/valueui/internal/status"
	"github.com/valuedig/valueui/internal/websrv"
	"github.com/valuedig/valueui/pkg/uiapi"
)

var (
	projectPath = "./project/"
	tplFileRx   = regexp.MustCompile(`(.*)(\.html)$`)
	dslFileRx   = regexp.MustCompile(`(.*)(\.toml)$`)
	err         error
)

func templateRefresh() error {

	if projectPath, err = filepath.Abs(projectPath); err != nil {
		return err
	}

	load := func(path string) error {

		if tplFileRx.MatchString(path) {
			if b, err := ioutil.ReadFile(path); err != nil {
				return err
			} else {
				status.Layouts.Sync(path[len(projectPath):], &uiapi.Template{
					Kind: "Template",
					Html: string(b),
				})
			}
		} else if dslFileRx.MatchString(path) {
			var item uiapi.LayoutContainer
			if err := htoml.DecodeFromFile(path, &item); err == nil {
				hlog.Printf("info", "layout %s, kind %s, name %v",
					path[len(projectPath):], item.Kind, item.Name)
				status.Layouts.Sync(path[len(projectPath):], &item)
			}
		}
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err = filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			hlog.Printf("info", "watch %s", path)
			return watcher.Add(path)
		}
		return load(path)
	}); err != nil {
		return err
	}

	var (
		updates = map[string]int64{}
	)

	for {
		select {

		case event, ok := <-watcher.Events:

			// hlog.Printf("info", "fsnotify event hit %v, file %v", event.Op, event.Name)

			if !ok || (!dslFileRx.MatchString(event.Name) &&
				!tplFileRx.MatchString(event.Name)) {
				continue
			}

			if (event.Op&fsnotify.Create) == fsnotify.Create ||
				(event.Op&fsnotify.Write) == fsnotify.Write ||
				// (event.Op&fsnotify.Rename) == fsnotify.Rename ||
				(event.Op&fsnotify.Remove) == fsnotify.Remove {

				tn := time.Now().UnixNano() / 1e6
				if (tn - updates[event.Name]) < 1e3 {
					continue
				}
				updates[event.Name] = tn

				time.Sleep(100e6)
				hlog.Printf("info", "fsnotify event %v, file %v", event.Op, event.Name)

				load(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				hlog.Printf("info", "fsnotify err %s", err.Error())
			}
		}
	}
}

func main() {

	httpsrv.GlobalService.Config.UrlBasePath = "/demo"
	httpsrv.GlobalService.Config.HttpPort = 8002

	httpsrv.GlobalService.ModuleRegister("/", websrv.NewModule())
	httpsrv.GlobalService.ModuleRegister("/api/v2", websrv.NewApiModule())

	hlog.Print("Running")
	go httpsrv.GlobalService.Start()

	if err := templateRefresh(); err != nil {
		panic(err)
	}
}
