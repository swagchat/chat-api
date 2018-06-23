package pbroker

type notuseProvider struct{}

func (np notuseProvider) PublishMessage(rtmEvent *RTMEvent) error {
	// Do not process anything
	return nil
}
