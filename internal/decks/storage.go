package decks

import (
	"log"
	"net/http"
	"time"

	"tts_deck_build/internal/collections"
	"tts_deck_build/internal/config"
	"tts_deck_build/internal/errors"
	"tts_deck_build/internal/fs"
	"tts_deck_build/internal/images"
	"tts_deck_build/internal/network"
	"tts_deck_build/internal/utils"
)

type DeckStorage struct {
	Config            *config.Config
	CollectionService *collections.CollectionService
}

func NewDeckStorage(config *config.Config, collectionService *collections.CollectionService) *DeckStorage {
	return &DeckStorage{
		Config:            config,
		CollectionService: collectionService,
	}
}

func (s *DeckStorage) Create(gameID, collectionID string, deck *DeckInfo) (*DeckInfo, error) {
	// Check ID
	if deck.ID == "" {
		return nil, errors.BadName.AddMessage(deck.Type)
	}

	// Check if such an object already exists
	if val, _ := s.GetByID(gameID, collectionID, deck.ID); val != nil {
		return nil, errors.DeckExist
	}

	// Writing info to file
	if err := fs.CreateAndProcess(deck.Path(gameID, collectionID), Deck{Deck: deck}, fs.JsonToWriter[Deck]); err != nil {
		return nil, err
	}

	if len(deck.BacksideImage) > 0 {
		// Download image
		if err := s.CreateImage(gameID, collectionID, deck.ID, deck.BacksideImage); err != nil {
			return nil, err
		}
	}

	return deck, nil
}
func (s *DeckStorage) GetByID(gameID, collectionID, deckID string) (*DeckInfo, error) {
	// Check if the collection exists
	_, err := s.CollectionService.Item(gameID, collectionID)
	if err != nil {
		return nil, err
	}

	deck := DeckInfo{ID: deckID}

	// Check if such an object exists
	isExist, err := fs.IsFileExist(deck.Path(gameID, collectionID))
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.DeckNotExists
	}

	// Read info from file
	readDeck, err := fs.OpenAndProcess(deck.Path(gameID, collectionID), fs.JsonFromReader[Deck])
	if err != nil {
		return nil, err
	}
	return readDeck.Deck, nil
}
func (s *DeckStorage) GetAll(gameID, collectionID string) ([]*DeckInfo, error) {
	decks := make([]*DeckInfo, 0)

	// Check if the collection exists
	collection, err := s.CollectionService.Item(gameID, collectionID)
	if err != nil {
		return decks, err
	}

	// Get list of objects
	folders, err := fs.ListOfFiles(collection.Path(gameID))
	if err != nil {
		return decks, err
	}

	// Get each deck
	for _, deckFileName := range folders {
		deckID := fs.GetFilenameWithoutExt(deckFileName)
		deck, err := s.GetByID(gameID, collectionID, deckID)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		if deck == nil {
			log.Println("Invalid deck file:", deckFileName)
			continue
		}
		decks = append(decks, deck)
	}

	return decks, nil
}
func (s *DeckStorage) Update(gameID, collectionID, deckID string, dto *UpdateDeckDTO) (*DeckInfo, error) {
	// Get old object
	oldDeck, err := s.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, err
	}

	// Create deck object
	deck := NewDeckInfo(dto.Type, dto.BacksideImage)
	deck.CreatedAt = oldDeck.CreatedAt
	if deck.ID == "" {
		return nil, errors.BadName.AddMessage(dto.Type)
	}

	// If the id has been changed, rename the object
	if deck.ID != oldDeck.ID {
		// Check if such an object already exists
		if val, _ := s.GetByID(gameID, collectionID, deck.ID); val != nil {
			return nil, errors.DeckExist
		}

		// If image exist, rename
		if data, _, _ := s.GetImage(gameID, collectionID, oldDeck.ID); data != nil {
			err = fs.MoveFolder(oldDeck.ImagePath(gameID, collectionID), deck.ImagePath(gameID, collectionID))
			if err != nil {
				return nil, err
			}
		}

		// Rename object
		err = fs.MoveFolder(oldDeck.Path(gameID, collectionID), deck.Path(gameID, collectionID))
		if err != nil {
			return nil, err
		}
	}

	// If the object has been changed, update the object file
	if !oldDeck.Compare(deck) {
		deck.UpdatedAt = utils.Allocate(time.Now())
		// Writing info to file
		if err := fs.CreateAndProcess(deck.Path(gameID, collectionID), Deck{Deck: deck}, fs.JsonToWriter[Deck]); err != nil {
			return nil, err
		}
	}

	// If the image has been changed
	if deck.BacksideImage != oldDeck.BacksideImage {
		// If image exist, delete
		if data, _, _ := s.GetImage(gameID, collectionID, deck.ID); data != nil {
			err = fs.RemoveFile(deck.ImagePath(gameID, collectionID))
			if err != nil {
				return nil, err
			}
		}

		if len(deck.BacksideImage) > 0 {
			// Download image
			if err = s.CreateImage(gameID, collectionID, deck.ID, deck.BacksideImage); err != nil {
				return nil, err
			}
		}
	}

	return deck, nil
}
func (s *DeckStorage) DeleteByID(gameID, collectionID, deckID string) error {
	deck := DeckInfo{ID: deckID}

	// Check if such an object exists
	if val, _ := s.GetByID(gameID, collectionID, deckID); val == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
	}

	// Remove object
	if err := fs.RemoveFile(deck.Path(gameID, collectionID)); err != nil {
		return err
	}

	// Remove image
	return fs.RemoveFile(deck.ImagePath(gameID, collectionID))
}
func (s *DeckStorage) GetImage(gameID, collectionID, deckID string) ([]byte, string, error) {
	// Check if such an object exists
	deck, err := s.GetByID(gameID, collectionID, deckID)
	if err != nil {
		return nil, "", err
	}

	// Check if an image exists
	isExist, err := fs.IsFileExist(deck.ImagePath(gameID, collectionID))
	if err != nil {
		return nil, "", err
	}
	if !isExist {
		return nil, "", errors.DeckImageNotExists
	}

	// Read an image from a file
	data, err := fs.OpenAndProcess(deck.ImagePath(gameID, collectionID), fs.BinFromReader)
	if err != nil {
		return nil, "", err
	}

	imgType, err := images.ValidateImage(data)
	if err != nil {
		return nil, "", err
	}

	return data, imgType, nil
}
func (s *DeckStorage) CreateImage(gameID, collectionID, deckID, imageURL string) error {
	// Check if such an object exists
	deck, _ := s.GetByID(gameID, collectionID, deckID)
	if deck == nil {
		return errors.DeckNotExists.HTTP(http.StatusBadRequest)
	}

	// Download image
	imageBytes, err := network.DownloadBytes(imageURL)
	if err != nil {
		return err
	}

	// Validate image
	_, err = images.ValidateImage(imageBytes)
	if err != nil {
		return err
	}

	// Write image to file
	return fs.CreateAndProcess(deck.ImagePath(gameID, collectionID), imageBytes, fs.BinToWriter)
}
func (s *DeckStorage) GetAllDecksInGame(gameID string) ([]*DeckInfo, error) {
	// Get all collections in selected game
	listCollections, err := s.CollectionService.List(gameID, "")
	if err != nil {
		return make([]*DeckInfo, 0), err
	}

	// Mark unique deck types
	uniqueDecks := make(map[string]struct{})

	// Go through all collections and find unique types of decks
	decks := make([]*DeckInfo, 0)
	for _, collection := range listCollections {
		// Get all decks in selected collection
		collectionDecks, err := s.GetAll(gameID, collection.ID)
		if err != nil {
			return make([]*DeckInfo, 0), err
		}

		// Go through all decks and keep only unique decks
		for _, deck := range collectionDecks {
			if _, ok := uniqueDecks[deck.Type+deck.BacksideImage]; ok {
				// If we have already seen such a deck, we skip it
				continue
			}
			// If deck unique, put mark in map
			uniqueDecks[deck.Type+deck.BacksideImage] = struct{}{}
			decks = append(decks, deck)
		}
	}
	return decks, nil
}
