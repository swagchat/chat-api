package rtm

type notuseProvider struct{}

func (np notuseProvider) Publish(rtmEvent *RTMEvent) error {
	// Do not process anything
	return nil
}
