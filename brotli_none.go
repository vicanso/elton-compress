// Copyright 2018 tree xie
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

// +build !brotli

package compress

import (
	"errors"

	"github.com/vicanso/cod"
)

type (
	brCompressor struct{}
)

// Accept just not accept all
func (b *brCompressor) Accept(_ *cod.Context) (acceptable bool, encoding string) {
	return
}

// Compress just return not support error
func (b *brCompressor) Compress(buf []byte, level int) ([]byte, error) {
	return nil, errors.New("not support brotli")
}
