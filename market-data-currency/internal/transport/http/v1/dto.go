package v1

import (
	"github.com/boldlogic/PortfolioLens/pkg/models"
)

type currencyDTO struct {
	ISOCode  int16  `json:"isoCode"`
	CharCode string `json:"code"`
	NameRu   string `json:"nameRu,omitempty"`
	NameEn   string `json:"nameEn,omitempty"`
}

func currencyToDTO(cur models.Currency) currencyDTO {
	var out currencyDTO
	out.ISOCode = cur.ISOCode
	out.CharCode = cur.ISOCharCode
	out.NameEn = cur.LatName

	if cur.Name != nil {
		out.NameRu = *cur.Name
	}

	return out
}
