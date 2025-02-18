package queue

import domain_primitives "chatgraph/domain/primitives"

type IQueue interface {
	GetMessages() (<-chan domain_primitives.UserCall, error)
}
