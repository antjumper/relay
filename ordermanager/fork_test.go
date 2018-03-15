/*

  Copyright 2017 Loopring Project Ltd (Loopring Foundation).

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.

*/

package ordermanager_test

import (
	"github.com/Loopring/relay/ordermanager"
	"sort"
	"testing"
)

func TestInnerForkEventList_Sort(t *testing.T) {
	var list ordermanager.InnerForkEventList
	list = append(list, ordermanager.InnerForkEvent{BlockNumber: 32, LogIndex: 3})
	list = append(list, ordermanager.InnerForkEvent{BlockNumber: 32, LogIndex: 2})
	list = append(list, ordermanager.InnerForkEvent{BlockNumber: 35, LogIndex: 3})

	sort.Sort(list)

	for _, v := range list {
		t.Logf("event blockNumber:%d, logIndex:%d", v.BlockNumber, v.LogIndex)
	}
}
