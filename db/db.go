package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GptProductDatabase interface {
	Init() error
	CreateProductSpec(
		sku string,
		query string,
		url string,
		datasheet string,
		catalogue string,
		output string,
		upperTx *gorm.DB,
	) error
}
type GptProductDatabaseImpl struct {
	mainDb *gorm.DB
}

func NewGptProductDatabase() GptProductDatabase {
	db := &GptProductDatabaseImpl{}
	err := db.Init()
	if err != nil {
		panic(err)
	}
	return db
}

type ProductSpec struct {
	Sku       string `gorm:"primaryKey;autoIncrement:false;unique"`
	Query     string
	Url       string
	Datasheet string
	Catalogue string
	Output    string
}

func (g *GptProductDatabaseImpl) Init() error {
	db, err := gorm.Open(sqlite.Open("gpt_products.db"), &gorm.Config{})
	if err != nil {
		return err
	}
	// Migrate the schema
	err = db.AutoMigrate(&ProductSpec{})
	if err != nil {
		return err
	}
	g.mainDb = db
	return nil
}

func (g *GptProductDatabaseImpl) CreateProductSpec(
	sku string,
	query string,
	url string,
	datasheet string,
	catalogue string,
	output string,
	upperTx *gorm.DB,
) error {
	createPrompt := func(tx *gorm.DB) error {
		p := &ProductSpec{
			Sku:       sku,
			Query:     query,
			Url:       url,
			Datasheet: datasheet,
			Catalogue: catalogue,
			Output:    output,
		}
		txErr := tx.Create(p)
		if txErr != nil {
			if txErr.Error != nil {
				if txErr.Error != gorm.ErrRecordNotFound {
					return txErr.Error
				}
			}
		}
		return nil
	}
	var err error
	if upperTx != nil {
		err = createPrompt(upperTx)
	} else {
		err = g.mainDb.Transaction(func(tx *gorm.DB) error {
			err = createPrompt(tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			return nil
		})
	}
	if err != nil {
		return err
	}
	return nil
}
