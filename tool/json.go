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

package tool

import (
	"bytes"
	"encoding/json"
)

func ToJSON(i interface{}) (string, error) {
	// Json the result
	jsonOutput, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	jsonerr := json.Indent(&out, jsonOutput, "", "\t")
	return out.String(), jsonerr
}
