package datastore

type BotStore interface {
	CreateBotStore()

	//	InsertBot(bot *models.Bot) StoreResult
	SelectBot(userId string) StoreResult
	//	SelectBotByUserId(userId string) StoreResult
	//	SelectBots() StoreResult
	//	UpdateBot(bot *models.Bot) StoreResult
	//	UpdateBotDeleted(botId string) StoreResult
}
