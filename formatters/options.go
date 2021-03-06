// © 2019 Nathan Galt
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package formatters

// Option is the type used for public-facing formatting options.
type Option func(o *formattingOptions)

// ForPublic is an option that says whether to sanitize the output for public consumption.
func ForPublic(b bool) Option {
	return func(o *formattingOptions) {
		o.forPublic = b
	}
}

type formattingOptions struct {
	forPublic bool
}
