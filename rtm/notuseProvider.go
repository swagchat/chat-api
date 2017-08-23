package rtm

type NotUseProvider struct{}

func (provider NotUseProvider) Init() error {
	return nil
}

func (provider NotUseProvider) PublishMessage(mi *MessagingInfo) error {
	return nil
}
