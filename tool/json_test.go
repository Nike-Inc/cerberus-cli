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
	"testing"
)

func TestToJSON(t *testing.T) {
	Convey("Given a valid map", t, func() {
		v := map[string]interface{}{"abc": "123",
			"xyz": "987",
		}

		Convey("When parsed to JSON", func() {
			jsoned, err := tool.ToJSON(v)

			Convey("The output is formatted correctly", func() {
				output := "{\n\t\"abc\": \"123\",\n\t\"xyz\": \"987\"\n}"
				So(jsoned, ShouldEqual, output)
				So(err, ShouldBeNil)
			})
		})
	})

	//var x string
	Convey("Given an invalid map", t, func() {
		v := map[interface{}]interface{}{nil: nil}

		Convey("When parsed to JSON", func() {
			jsoned, err := tool.ToJSON(v)

			Convey("An error occurs", func() {
				So(jsoned, ShouldBeEmpty)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
