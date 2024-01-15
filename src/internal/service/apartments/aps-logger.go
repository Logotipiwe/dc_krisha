package apartments

import (
	"fmt"
	"krisha/src/internal/domain"
	"krisha/src/internal/service/tg"
	"log"
	"strconv"
)

type ApsLoggerService struct{}

func NewApsLoggerService() *ApsLoggerService {
	return &ApsLoggerService{}
}

func (s *ApsLoggerService) LogMissingAp(m *domain.Ap) {
	log.Println(fmt.Sprintf("Missing ap %s", strconv.FormatInt(m.ID, 10)))
	log.Println(m)
}

func (s *ApsLoggerService) LogNewAp(data *domain.Ap) {
	log.Println("=======================================================================")
	log.Println("NEW AP FOUND")
	log.Println("ID:")
	log.Println(data.ID)

	log.Println(data.GetLink())
	log.Println("=======================================================================")
}

type ApsTgSenderService struct {
	tgService *tg.TgService
}

func NewApsTgSenderService(tgService *tg.TgService) *ApsTgSenderService {
	return &ApsTgSenderService{tgService}
}

func (s *ApsTgSenderService) LogInTg(data *domain.Ap) error {
	var imagesUrls = make([]string, 0)
	photos := data.Photos
	message := fmt.Sprintf("Link: %s\r\n", data.GetLink())
	if len(photos) > 0 {
		for _, photo := range photos {
			imagesUrls = append(imagesUrls, photo.Src)
		}
		if len(imagesUrls) > 10 {
			imagesUrls = imagesUrls[:10]
		}
		return s.tgService.SendImgMessageToOwner(message, imagesUrls)
	} else {
		return s.tgService.SendMessageToOwner(message)
	}
}
