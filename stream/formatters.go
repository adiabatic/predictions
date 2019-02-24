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

package stream

func deduplicateTags(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	j := 0
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ss[j] = v
		j++
	}
	return ss[:j]
}

// TagsUsed returns a list of all tags used in the given slice of stream.Stream.
func TagsUsed(ss []Stream) []string {
	ret := make([]string, 0)
	for _, s := range ss {
		for _, p := range s.Predictions {
			for _, t := range p.Tags {
				ret = append(ret, t)
			}
		}
	}
	return deduplicateTags(ret)
}
