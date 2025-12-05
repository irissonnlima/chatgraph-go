package d_context

import (
	"encoding/json"
)

func (c *ChatContext[Obs]) GetObservation() Obs {
	return c.UserState.Observation
}

func (c *ChatContext[Obs]) SetObservation(observation Obs) error {
	if c.Context.Err() != nil {
		return c.Context.Err()
	}

	obsString, err := json.Marshal(observation)
	if err != nil {
		return err
	}
	return c.router.SetObservation(c.UserState.ChatID, string(obsString))
}
