package repositories

import (
	"api/interfaces"
	"api/models"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Define a estrutura do repositório do DynamoDB.
type DynamoDB struct {
	client interfaces.DynamoDBClient
	ttl    time.Duration
	name   string
}

// Cria uma nova instância do repositório do DynamoDB.
func NewDynamoDBRepository(client interfaces.DynamoDBClient, name string, ttl time.Duration) *DynamoDB {
	return &DynamoDB{
		client: client,
		name:   name,
		ttl:    ttl,
	}
}

// Cria a tabela DynamoDB com os índices secundários globais necessários.
func (p *DynamoDB) Create(ctx context.Context) error {
	_, err := p.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("date"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("statusCode"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("date-statusCode-index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("date"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("statusCode"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		TableName:   &p.name,
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	// deve aguardar até a tabela ser criada e estar disponível para uso
	waiter := dynamodb.NewTableExistsWaiter(p.client)
	err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{TableName: &p.name}, 5*time.Minute)
	if err != nil {
		return err
	}
	// só é possível habilitar TTL na tabela após ela ter sido criada
	_, err = p.client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: &p.name,
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String("expiration"),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		return err
	}
	return nil
}

// Salva o registro na tabela DynamoDB.
// Se já houver registro com o mesmo id, ele será substituído.
func (p *DynamoDB) Save(ctx context.Context, event *models.Event) error {
	if event.Expiration == 0 {
		event.Expiration = time.Now().Add(p.ttl).Unix()
	}
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return err
	}
	_, err = p.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &p.name,
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

// Deleta o registro da tabela DynamoDB pelo id.
func (p *DynamoDB) Delete(ctx context.Context, id string) (event *models.Event, err error) {
	out, err := p.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &p.name,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return nil, err
	}
	if out.Attributes == nil {
		return nil, nil
	}
	err = attributevalue.UnmarshalMap(out.Attributes, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Recupera o registro da tabela DynamoDB pelo id.
func (p *DynamoDB) Get(ctx context.Context, id string) (event *models.Event, err error) {
	out, err := p.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &p.name,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, nil
	}
	err = attributevalue.UnmarshalMap(out.Item, &event)
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Procura registros com a data entre o período especificado e com o status code fornecido.
func (p *DynamoDB) FindByDateAndReturnCode(ctx context.Context, from time.Time, to time.Time, statusCode int) (events []*models.Event, err error) {
	condition := &dynamodb.QueryInput{
		TableName: aws.String(p.name),
		IndexName: aws.String("date-statusCode-index"),
		KeyConditionExpression: aws.String(
			"statusCode = :statusCode AND #date BETWEEN :from AND :to",
		),
		ExpressionAttributeNames: map[string]string{
			"#date": "date", // palavra reservada no DynamoDB
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":statusCode": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", statusCode)},
			":from":       &types.AttributeValueMemberS{Value: from.Format(time.RFC3339)},
			":to":         &types.AttributeValueMemberS{Value: to.Format(time.RFC3339)},
		},
	}
	paginator := dynamodb.NewQueryPaginator(p.client, condition)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range page.Items {
			event := &models.Event{}
			err = attributevalue.UnmarshalMap(item, event)
			if err != nil {
				return nil, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}
