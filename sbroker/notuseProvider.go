package sbroker

type notuseProvider struct{}

func (np notuseProvider) SubscribeMessage() error {
	// Do not process anything
	return nil
}

func (np notuseProvider) UnsubscribeMessage() error {
	// Do not process anything
	return nil
}
