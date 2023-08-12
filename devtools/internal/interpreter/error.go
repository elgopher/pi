// (c) 2023 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package interpreter

type ErrInvalidIdentifier struct {
	message string
}

func (e ErrInvalidIdentifier) Error() string {
	return e.message
}
