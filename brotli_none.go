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

	"github.com/vicanso/elton"
)

type (
	BrCompressor struct{}
)

// Accept just not accept all
func (b *BrCompressor) Accept(_ *elton.Context) (acceptable bool, encoding string) {
	return false, ""
}

// Compress just return not support error
func (b *BrCompressor) Compress(buf []byte, level int) ([]byte, error) {
	return nil, errors.New("not support brotli")
}

// Pipe brotli pipe
func (b *BrCompressor) Pipe(c *elton.Context, level int) (err error) {
	return errors.New("not support brotli")
}
