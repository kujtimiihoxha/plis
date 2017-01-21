// Copyright © 2016 Kujtim Hoxha <kujtimii.h@gmail.com>
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

package fs

import (
	"github.com/kujtimiihoxha/plis/helpers"
	"github.com/spf13/afero"
)

var (
	workingDirFs           afero.Fs
	generatorTemplateDirFs afero.Fs
)

func Init() {
	workingDirFs = &afero.OsFs{}
}
func WorkingDirFs() afero.Fs {
	return workingDirFs
}
func TemplatesDirFs() afero.Fs {
	return generatorTemplateDirFs
}

func InitTemplatesDirFs(generator string) {
	generatorTemplateDirFs = afero.NewBasePathFs(workingDirFs, helpers.GeneratorTemplatePath(generator))
}

func SetToMemSystem() {
	workingDirFs = afero.NewMemMapFs()
}
