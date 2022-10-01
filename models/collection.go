package models;

type Collection struct {
  Id        int64
  Name      string
  Requests  []*UrlRequest
}
