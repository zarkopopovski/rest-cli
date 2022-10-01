package models;

type UrlRequest struct {
  Id            int64
  CollectionID  int64
  Url           string
  Method        string
  ParamsData    string
  HeaderData    string
  CookieData    string
  BodyData      string
}
