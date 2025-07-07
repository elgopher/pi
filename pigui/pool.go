// Copyright 2025 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package pigui

import (
	"log"

	"github.com/elgopher/pi/pipool"
)

func getPropagateToChildrenFromThePool() *propagateToChildren {
	propagate := propagateToChildrenPool.Get()
	propagateToChildrenToken++
	propagate.value = true
	propagate.token = propagateToChildrenToken
	return propagate
}

var propagateToChildrenToken = 0

// the object is valid only during element's Draw/Update.
type propagateToChildren struct {
	value bool
	token int // token is verified when setting the value.
}

func (p *propagateToChildren) set(v bool, token int) {
	if p.token == token { // when object is reused the token is changed
		p.value = v
	} else {
		log.Println("code was trying to stop event propagation after pigui element was updated/drawn")
	}
}

var propagateToChildrenPool pipool.Pool[propagateToChildren]
