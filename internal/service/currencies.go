package service

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"net/http"

	"golang.org/x/net/html/charset"
)

type ValItem struct {
	Id          string `xml:"ID,attr"`
	Name        string `xml:"Name"`
	EngName     string `xml:"EngName"`
	Nominal     int    `xml:"Nominal"`
	ParentCode  string `xml:"ParentCode"`
	ISONumCode  int    `xml:"ISO_Num_Code"`
	ISOCharCode string `xml:"ISO_Char_Code"`
}

type Valuta struct {
	Name string    `xml:"name,attr"`
	Item []ValItem `xml:"Item"`
}

func (i ValItem) String() string {
	return fmt.Sprintf("ID: %s, Name: %s, EngName: %s,Nominal: %d, ParentCode: %s,ISO_Num_Code: %d, ISO_Char_Code: %s", i.Id, i.Name, i.EngName, i.Nominal, i.ParentCode, i.ISONumCode, i.ISOCharCode)
}
func (val Valuta) String() string {
	return fmt.Sprintf("Name: %s, Item: %s", val.Name, val.Item)
}

func (c *Client) GetCbrCurrencies(ctx context.Context, req *http.Request) error {
	resp, err := c.sendRequest(ctx, req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Не удалось получить справочник валют: %d", resp.StatusCode)
	}
	decoder := xml.NewDecoder(bytes.NewReader(resp.Body))
	decoder.CharsetReader = charset.NewReaderLabel
	var val Valuta
	if err := decoder.Decode(&val); err != nil {
		return fmt.Errorf("Не удалось получить справочник валют: %w", err)
	}
	c.log.Infof("%s", val)

	return nil
}
