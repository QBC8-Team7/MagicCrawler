package notification

import (
	"github.com/QBC8-Team7/MagicCrawler/config"
	"github.com/QBC8-Team7/MagicCrawler/internal/repositories"
	"github.com/QBC8-Team7/MagicCrawler/pkg/db/sqlc"
	"github.com/QBC8-Team7/MagicCrawler/pkg/logger"
)

type AdminNotifier struct {
	repo  *repositories.AdminRepository
	notif *Service
}

func NewAdminNotifier(conf *config.Config, queries *sqlc.Queries) (*AdminNotifier, error) {
	notificationLogger := logger.NewAppLogger(conf)
	notificationLogger.InitCustomLogger(conf.Logger.Path, conf.Logger.SysPath)
	notificationService, err := GetService(conf.Bot.Token, notificationLogger)
	if err != nil {
		return &AdminNotifier{}, err
	}

	return &AdminNotifier{
		repo:  repositories.NewAdminRepository(queries),
		notif: notificationService,
	}, nil
}

func (n AdminNotifier) Send(message string) error {
	admin, err := n.repo.GetNextAdmin()
	if err != nil {
		return err
	}

	err = n.notif.SendMessage(admin.TgID, message)
	return err
}
