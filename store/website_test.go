// Copyright 2015 CoreStore Authors
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

package store_test

//func TestWebsite(t *testing.T) {
//	wsInvalid, err := storeManager.Website().Get(312)
//	assert.Nil(t, wsInvalid)
//	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())
//
//	ws, err := storeManager.Website().Get(1)
//	assert.NoError(t, err)
//	assert.Equal(t, `Main Website`, ws.Name.String)
//	assert.Equal(t, "base", ws.Code.String)
//	assert.True(t, ws.Code.Valid)
//
//	wsInvalid, err = storeManager.Website().Get(-1, "oxid")
//	assert.Nil(t, wsInvalid)
//	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())
//
//	ws, err = storeManager.Website().Get(-1, "base")
//	assert.NoError(t, err)
//	assert.Equal(t, `Main Website`, ws.Name.String)
//	assert.Equal(t, "base", ws.Code.String)
//	assert.True(t, ws.Code.Valid)
//
//	wc := storeManager.Website().Collection()
//	assert.NotNil(t, wc)
//	assert.Len(t, wc, 3)
//}
//
//func TestWebsiteGroup(t *testing.T) {
//	gInvalid, err := storeManager.Website().Group(321)
//	assert.Nil(t, gInvalid)
//	assert.EqualError(t, store.ErrWebsiteNotFound, err.Error())
//
//	group, err := storeManager.Website().Group(1)
//	assert.NoError(t, err)
//	t.Logf("\n%#v\n", group)
//}
