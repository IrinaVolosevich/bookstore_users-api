package services

var (
	ItemsService ItemServiceInterface = &itemsService{}
)

type ItemServiceInterface interface {
	GetItem()
	SaveItem()
}

type itemsService struct {
}

func (s *itemsService) GetItem() {

}

func (s *itemsService) SaveItem() {

}
