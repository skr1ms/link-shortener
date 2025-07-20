package link

import (
	"linkshortener/pkg/db"

	"gorm.io/gorm/clause"
)

type LinkRepository struct {
	db *db.Db
}

func NewLinkRepository(db *db.Db) *LinkRepository {
	return &LinkRepository{db: db}
}

func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.db.DB.Table("links").First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}

func (repo *LinkRepository) Create(link *Link) (*Link, error) {
	link.Hash = CheckUniqueAndGenerateHash(repo.db)

	result := repo.db.DB.Table("links").Create(link)
	if result.Error != nil {
		return nil, result.Error
	}

	return link, nil
}

func (repo *LinkRepository) Update(link *Link) (*Link, error) {
	result := repo.db.DB.Clauses(clause.Returning{}).Updates(link)
	if result.Error != nil {
		return nil, result.Error
	}
	return link, nil
}

func (repo *LinkRepository) FindById(id uint) error {
	result := repo.db.DB.Table("links").Find(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *LinkRepository) Delete(id uint) error {
	if err := repo.FindById(id); err != nil {
		return err
	}

	result := repo.db.DB.Table("links").Delete(&Link{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *LinkRepository) GetLinks(limit, offset uint) ([]Link, error) {
	var links []Link

	result := repo.db.DB.
		Table("links").
		Where("deleted_at IS NULL").
		Order("id DESC").
		Limit(int(limit)).
		Offset(int(offset)).
		Scan(&links)

	if result.Error != nil {
		return nil, result.Error
	}
	return links, nil
}

func (repo *LinkRepository) GetLinksCount() (int64, error) {
	var count int64
	result := repo.db.DB.
		Table("links").
		Where("deleted_at IS NULL").
		Count(&count)
		
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}
