package link_test

import (
	"errors"
	"linkshortener/internal/link"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type MockDB struct {
	links     map[string]*link.Link
	linksById map[uint]*link.Link
	nextID    uint
}

func NewMockDB() *MockDB {
	return &MockDB{
		links:     make(map[string]*link.Link),
		linksById: make(map[uint]*link.Link),
		nextID:    1,
	}
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	result := &gorm.DB{}

	if linkPtr, ok := dest.(*link.Link); ok {
		if len(conds) >= 2 {
			if conds[0] == "hash = ?" {
				hash := conds[1].(string)
				if foundLink, exists := m.links[hash]; exists {
					*linkPtr = *foundLink
				} else {
					result.Error = gorm.ErrRecordNotFound
				}
			}
		} else if len(conds) == 1 {
			// Поиск по ID
			id := conds[0].(uint)
			if foundLink, exists := m.linksById[id]; exists {
				*linkPtr = *foundLink
			} else {
				result.Error = gorm.ErrRecordNotFound
			}
		}
	}

	return result
}

type MockLinkRepository struct {
	db *MockDB
}

func NewMockLinkRepository() *MockLinkRepository {
	return &MockLinkRepository{
		db: NewMockDB(),
	}
}

func (repo *MockLinkRepository) GetByHash(hash string) (*link.Link, error) {
	if linkItem, exists := repo.db.links[hash]; exists {
		return linkItem, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (repo *MockLinkRepository) Create(linkItem *link.Link) (*link.Link, error) {
	linkItem.Hash = generateUniqueHash(repo.db)
	linkItem.ID = repo.db.nextID
	repo.db.nextID++

	repo.db.links[linkItem.Hash] = linkItem
	repo.db.linksById[linkItem.ID] = linkItem

	return linkItem, nil
}

func (repo *MockLinkRepository) Update(linkItem *link.Link) (*link.Link, error) {
	if existingLink, exists := repo.db.linksById[linkItem.ID]; exists {
		existingLink.OriginalURL = linkItem.OriginalURL
		existingLink.Hash = linkItem.Hash
		repo.db.links[linkItem.Hash] = existingLink
		return existingLink, nil
	}
	return nil, errors.New("link not found")
}

func (repo *MockLinkRepository) FindById(id uint) error {
	if _, exists := repo.db.linksById[id]; exists {
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (repo *MockLinkRepository) Delete(id uint) error {
	if err := repo.FindById(id); err != nil {
		return err
	}

	linkItem := repo.db.linksById[id]
	delete(repo.db.linksById, id)
	delete(repo.db.links, linkItem.Hash)

	return nil
}

func (repo *MockLinkRepository) GetLinks(limit, offset uint) ([]link.Link, error) {
	links := make([]link.Link, 0)
	count := uint(0)

	for _, linkItem := range repo.db.linksById {
		if count >= offset {
			if uint(len(links)) >= limit {
				break
			}
			links = append(links, *linkItem)
		}
		count++
	}

	return links, nil
}

func (repo *MockLinkRepository) GetLinksCount() (int64, error) {
	return int64(len(repo.db.linksById)), nil
}

func generateUniqueHash(mockDB *MockDB) string {
	for {
		hash := link.GenerateHash()
		if _, exists := mockDB.links[hash]; !exists {
			return hash
		}
	}
}

func TestLinkRepositoryGetByHashSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	testLink := &link.Link{
		OriginalURL: "https://example.com",
		Hash:        "test123",
	}
	testLink.ID = 1
	repo.db.links["test123"] = testLink
	repo.db.linksById[1] = testLink

	linkItem, err := repo.GetByHash("test123")

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if linkItem == nil {
		t.Fatal("Expected link to be found")
	}

	if linkItem.OriginalURL != "https://example.com" {
		t.Fatalf("Expected URL https://example.com, got %s", linkItem.OriginalURL)
	}
}

func TestLinkRepositoryGetByHashNotFound(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	_, err := repo.GetByHash("nonexistent")

	if err == nil {
		t.Fatal("Expected error for non-existent hash")
	}

	if err != gorm.ErrRecordNotFound {
		t.Fatalf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestLinkRepositoryCreateSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	newLink := link.NewLink("https://google.com")

	createdLink, err := repo.Create(newLink)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdLink == nil {
		t.Fatal("Expected created link to be returned")
	}

	if createdLink.OriginalURL != "https://google.com" {
		t.Fatalf("Expected URL https://google.com, got %s", createdLink.OriginalURL)
	}

	if createdLink.Hash == "" {
		t.Fatal("Expected hash to be generated")
	}

	if createdLink.ID == 0 {
		t.Fatal("Expected ID to be assigned")
	}
}

func TestLinkRepositoryUpdateSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	testLink := &link.Link{
		OriginalURL: "https://example.com",
		Hash:        "test123",
	}
	testLink.ID = 1
	repo.db.links["test123"] = testLink
	repo.db.linksById[1] = testLink

	updateLink := &link.Link{
		OriginalURL: "https://updated.com",
		Hash:        "updated123",
	}
	updateLink.ID = 1

	updatedLink, err := repo.Update(updateLink)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedLink.OriginalURL != "https://updated.com" {
		t.Fatalf("Expected URL https://updated.com, got %s", updatedLink.OriginalURL)
	}

	if updatedLink.Hash != "updated123" {
		t.Fatalf("Expected hash updated123, got %s", updatedLink.Hash)
	}
}

func TestLinkRepositoryDeleteSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	testLink := &link.Link{
		OriginalURL: "https://example.com",
		Hash:        "test123",
	}
	testLink.ID = 1
	repo.db.links["test123"] = testLink
	repo.db.linksById[1] = testLink

	err := repo.Delete(1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = repo.GetByHash("test123")
	if err == nil {
		t.Fatal("Expected link to be deleted")
	}
}

func TestLinkRepositoryDeleteNotFound(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	err := repo.Delete(999)

	if err == nil {
		t.Fatal("Expected error for non-existent link")
	}
}

func TestLinkRepositoryGetLinksSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	for i := 1; i <= 5; i++ {
		linkItem := &link.Link{
			OriginalURL: "https://example.com",
			Hash:        "test" + string(rune(i)),
		}
		linkItem.ID = uint(i)
		repo.db.links[linkItem.Hash] = linkItem
		repo.db.linksById[uint(i)] = linkItem
	}

	links, err := repo.GetLinks(3, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(links) != 3 {
		t.Fatalf("Expected 3 links, got %d", len(links))
	}
}

func TestLinkRepositoryGetLinksCountSuccess(t *testing.T) {
	godotenv.Load()

	repo := NewMockLinkRepository()

	for i := 1; i <= 7; i++ {
		linkItem := &link.Link{
			OriginalURL: "https://example.com",
			Hash:        "test" + string(rune(i)),
		}
		linkItem.ID = uint(i)
		repo.db.links[linkItem.Hash] = linkItem
		repo.db.linksById[uint(i)] = linkItem
	}

	count, err := repo.GetLinksCount()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if count != 7 {
		t.Fatalf("Expected count 7, got %d", count)
	}
}
