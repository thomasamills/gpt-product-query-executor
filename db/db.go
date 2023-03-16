package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GptProductDatabase interface {
	Init(conf MySQLConnectionConfig) error
	CreateProductSpec(
		sku string,
		query string,
		url string,
		datasheet string,
		catalogue string,
		output string,
		upperTx *gorm.DB,
	) error
	DoesProductSpecExist(
		sku string,
		upperTx *gorm.DB,
	) (bool, error)
}
type GptProductDatabaseImpl struct {
	mainDb *gorm.DB
}

func NewGptProductDatabase(config MySQLConnectionConfig) GptProductDatabase {
	db := &GptProductDatabaseImpl{}
	err := db.Init(config)
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

type MySQLConnectionConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DbName   string
}

func (g *GptProductDatabaseImpl) Init(conf MySQLConnectionConfig) error {
	db, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.DbName),
	), &gorm.Config{},
	)
	if err != nil {
		return err
	}
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
	createProductSpec := func(tx *gorm.DB) error {
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
		err = createProductSpec(upperTx)
	} else {
		err = g.mainDb.Transaction(func(tx *gorm.DB) error {
			err = createProductSpec(tx)
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

func (g *GptProductDatabaseImpl) DoesProductSpecExist(
	sku string,
	upperTx *gorm.DB,
) (bool, error) {
	result := &ProductSpec{}
	getPrompt := func(tx *gorm.DB) error {
		err := tx.
			Where("sku = ?", sku).
			First(&result).
			Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return err
			}
		}
		return nil
	}
	var err error
	if upperTx != nil {
		err = getPrompt(upperTx)
	} else {
		err = g.mainDb.Transaction(func(tx *gorm.DB) error {
			err = getPrompt(tx)
			if err != nil {
				tx.Rollback()
				return nil
			}

			return nil
		})
	}

	if err != nil {
		return false, err
	}

	return result != nil, nil
}
