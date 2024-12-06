package homeboxclient

import "time"

type PaginationResult[T any] struct {
	Items    []T `json:"items"`
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type PaginatedItems struct {
	Items    []Item `json:"items"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	Total    int    `json:"total"`
}

type Item struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Attachments   []Attachment `json:"attachments"`
	ImageID       string       `json:"imageId"`
	AssetID       string       `json:"assetId"`
	Archived      bool         `json:"archived"`
	Insured       bool         `json:"insured"`
	Quantity      int          `json:"quantity"`
	PurchasePrice float64      `json:"purchasePrice"`
	PurchaseFrom  string       `json:"purchaseFrom"`
	PurchaseTime  time.Time    `json:"purchaseTime"`
	CreatedAt     time.Time    `json:"createdAt"`
	UpdatedAt     time.Time    `json:"updatedAt"`
	Notes         string       `json:"notes"`
	Labels        []Label      `json:"labels"`
	Fields        []ItemField  `json:"fields"`
}

type Attachment struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Document  DocumentOut `json:"document"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	Primary   bool        `json:"primary"`
}

type DocumentOut struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Title string `json:"title"`
}

type ItemField struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	TextValue    string `json:"textValue"`
	NumberValue  int    `json:"numberValue"`
	BooleanValue bool   `json:"booleanValue"`
}

type AttachmentToken struct {
	Token string `json:"token"`
}

type ItemCreate struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	LabelIDs    []string `json:"labelIds"`
	LocationID  string   `json:"locationId"`
	ParentID    string   `json:"parentId,omitempty"`
}

type ItemUpdate struct {
	Name            string      `json:"name"`
	Description     string      `json:"description"`
	LabelIDs        []string    `json:"labelIds"`
	LocationID      string      `json:"locationId"`
	ParentID        string      `json:"parentId,omitempty"`
	Fields          []ItemField `json:"fields"`
	Archived        bool        `json:"archived"`
	AssetID         string      `json:"assetId"`
	Insured         bool        `json:"insured"`
	Manufacturer    string      `json:"manufacturer"`
	ModelNumber     string      `json:"modelNumber"`
	SerialNumber    string      `json:"serialNumber"`
	PurchaseFrom    string      `json:"purchaseFrom"`
	PurchasePrice   float64     `json:"purchasePrice,omitempty"`
	PurchaseTime    time.Time   `json:"purchaseTime"`
	Quantity        int         `json:"quantity"`
	Notes           string      `json:"notes"`
	WarrantyDetails string      `json:"warrantyDetails"`
	WarrantyExpires time.Time   `json:"warrantyExpires"`
}
