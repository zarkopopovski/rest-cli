package services;

import (
  //"fmt"
  //"log"
  
	//"time"
	
	"github.com/zarkopopovski/rest-cli/models"
	"github.com/zarkopopovski/rest-cli/db"
)

type DBService struct {
  DBManager *db.DBManager
}

func (userService *DBService) CreateDBConnection() {
	dbManager := &db.DBManager{}
	dbManager.OpenConnection()
	
	userService.DBManager = dbManager
}

func (userService *DBService) CreateNewCollection(collectionName string) (err error) {
	query := "INSERT INTO collection(name, date_created, date_modified) VALUES($1, datetime('now'), datetime('now'));"

	_, err = userService.DBManager.DB.Exec(query, collectionName)

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) UpdateCollectionName(collectionID int64, collectionName string) (err error) {
	query := "UPDATE collection SET name='$1', date_modified=datetime('now') WHERE id=$2;"

	_, err = userService.DBManager.DB.Exec(query, collectionName, collectionID)

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) DeleteCollection(collectionID int64) (err error) {
	query := "DELETE FROM collection WHERE id=$1;"

	_, err = userService.DBManager.DB.Exec(query, collectionID)

	if err == nil {
		return nil
	}

	err = userService.DeleteAllCollectionRequests(collectionID)

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) ListAllCollection() (collections []*models.Collection, err error) {
	query := "SELECT id, name FROM collection ORDER BY date_created DESC;"

	rows, err := userService.DBManager.DB.Query(query)
	
	if err == nil {
		collections := make([]*models.Collection, 0)

		for rows.Next() {
			newCollection := new(models.Collection)

			err = rows.Scan(&newCollection.Id, &newCollection.Name)
			
			if err != nil {
					return nil, err
			}
			
			collections = append(collections, newCollection)
		}

		return collections, nil
	}

	return nil, err
}

func (userService *DBService) FindLastCollection() (collection *models.Collection, err error) {
	query := "SELECT id, name FROM collection ORDER BY date_created DESC LIMIT 1;"

	row := userService.DBManager.DB.QueryRow(query)
	
	if row != nil {
		newCollection := new(models.Collection)

		err = row.Scan(&newCollection.Id, &newCollection.Name)
		
		if err != nil {
			return nil, err
		}

		return newCollection, nil
	}

	return nil, err
}

func (userService *DBService) ListAllRequestsForCollection(collectionID int64) (urlRequests []*models.UrlRequest, err error) {
	query := "SELECT id, collection_id, url, method, params_data, header_data, cookie_data, body_data FROM url_request WHERE collection_id=$1 ORDER BY date_created DESC"

	rows, err := userService.DBManager.DB.Query(query, collectionID)
	
	if err == nil {
		urlRequestModels := make([]*models.UrlRequest, 0)

		for rows.Next() {
			newUrlRequest := new(models.UrlRequest)

			err = rows.Scan(&newUrlRequest.Id, &newUrlRequest.CollectionID, &newUrlRequest.Url, &newUrlRequest.Method, &newUrlRequest.ParamsData, &newUrlRequest.HeaderData, &newUrlRequest.CookieData, &newUrlRequest.BodyData)
			
			if err != nil {
					return nil, err
			}
			
			urlRequestModels = append(urlRequestModels, newUrlRequest)
		}

		return urlRequestModels, nil
	}

	return nil, err
}

func (userService *DBService) CreateNewCollectionRequest(collectionID int64, requestData map[string]string) (err error) {
	query := "INSERT INTO url_request(collection_id, name, url, method, params_data, header_data, cookie_data, body_data, date_created, date_modified) VALUES($1, $2, $3, $4, $5, $6, $7, $8, datetime('now'), datetime('now'));"

	_, err = userService.DBManager.DB.Exec(query, collectionID, requestData["name"], requestData["url"], requestData["method"], requestData["params_data"], requestData["header_data"], requestData["cookie_data"], requestData["body_data"])

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) UpdateCollectionRequest(collRequestID int64, requestData map[string]string) (err error) {
	query := "UPDATE url_request SET name='$1', url='$2', method='$3', params_data='$4', header_data='$5', cookie_data='$6', body_data='$7', date_modified=datetime('now') WHERE id=$8;"

	_, err = userService.DBManager.DB.Exec(query, requestData["name"], requestData["url"], requestData["method"], requestData["params_data"], requestData["header_data"], requestData["cookie_data"], requestData["body_data"], collRequestID)

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) DeleteCollectionRequest(collRequestID int64) (err error) {
	query := "DELETE FROM url_request WHERE id=$1;"

	_, err = userService.DBManager.DB.Exec(query, collRequestID)

	if err == nil {
		return nil
	}

	return err
}

func (userService *DBService) DeleteAllCollectionRequests(collectionID int64) (err error) {
	query := "DELETE FROM url_request WHERE collection_id=$1;"

	_, err = userService.DBManager.DB.Exec(query, collectionID)

	if err == nil {
		return nil
	}

	return err
}
