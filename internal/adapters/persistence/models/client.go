package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/MAD-py/pandora-core/internal/adapters/persistence/models/utils"
	"github.com/MAD-py/pandora-core/internal/domain/entities"
	"github.com/MAD-py/pandora-core/internal/domain/enums"
	"github.com/jackc/pgx/v5/pgtype"
)

var clientType = []string{"organization", "developer"}

type Client struct {
	ID pgtype.Int4

	Type  pgtype.Text
	Name  pgtype.Text
	Email pgtype.Text

	CreatedAt pgtype.Timestamptz
}

func (c *Client) ValidateModel() error {
	return c.validateType()
}

func (c *Client) validateType() error {
	if t, _ := c.Type.Value(); t != nil {
		if slices.Contains(clientType, t.(string)) {
			return nil
		}
	}

	return fmt.Errorf(
		"invalid type: must be %s", strings.Join(clientType, ", "),
	)
}

func (c *Client) ToEntity() (*entities.Client, error) {
	clientType, err := enums.ParseClientType(utils.PgtypeTextToString(c.Type))
	if err != nil {
		return nil, err
	}

	return &entities.Client{
		ID:        utils.PgtypeInt4ToInt(c.ID),
		Type:      clientType,
		Name:      utils.PgtypeTextToString(c.Name),
		Email:     utils.PgtypeTextToString(c.Email),
		CreatedAt: utils.PgtypeTimestamptzToTime(c.CreatedAt),
	}, nil
}

func ClientsToEntity(array []*Client) ([]*entities.Client, error) {
	result := make([]*entities.Client, len(array))
	for i, v := range array {
		vv, err := v.ToEntity()
		if err != nil {
			return nil, err
		}
		result[i] = vv
	}
	return result, nil
}

func ClientFromEntity(client *entities.Client) *Client {
	return &Client{
		ID:        utils.IntToPgtypeInt4(client.ID),
		Type:      utils.StringToPgtypeText(client.Type.String()),
		Name:      utils.StringToPgtypeText(client.Name),
		Email:     utils.StringToPgtypeText(client.Email),
		CreatedAt: utils.TimeToPgtypeTimestamptz(client.CreatedAt),
	}
}
