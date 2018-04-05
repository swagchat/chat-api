package rtm

type notuseProvider struct{}

func (np notuseProvider) PublishMessage(mi *MessagingInfo) error {
	// Do not process anything
	return nil
}
