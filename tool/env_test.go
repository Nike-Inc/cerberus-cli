/*
 *  Copyright (c) 2019 Nike, Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package tool_test

import (
	"cerberus-cli/tool"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestGetEnvVariable(t *testing.T) {
	Convey("Given an environment variable is set to some value", t, func() {
		token := "TESTINGVALUE"
		os.Setenv(tool.EnvCerbToken, token)

		Convey("When the value of the variable is retrieved", func() {
			value := tool.GetEnvVariable(tool.EnvCerbToken)

			Convey("The value returned by GetEnvVariable should be the same", func() {
				So(value, ShouldEqual, token)

			})
		})
	})
	Convey("And when the variable is unset", t, func() {
		os.Unsetenv(tool.EnvCerbToken)

		Convey("The value returned by GetEnvVariable should be an empty string", func() {
			So(tool.GetEnvVariable(tool.EnvCerbToken), ShouldBeEmpty)
		})
	})
}
