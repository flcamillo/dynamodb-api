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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Define a configuração do repositório do DynamoDB.
type DynamoDBConfig struct {
	// log de aplicação
	Log interfaces.Log
	// cliente do DynamoDB
	Client interfaces.DynamoDBClient
	// nome da tabela
	Table string
	// tempo de expiração dos registros
	TTL time.Duration
}

// Define a estrutura do repositório do DynamoDB.
type DynamoDB struct {
	// cliente do DynamoDB
	config *DynamoDBConfig
	// configura o tracer
	tracer trace.Tracer
}

// Cria uma nova instância do repositório do DynamoDB.
func NewDynamoDBRepository(config *DynamoDBConfig) *DynamoDB {
	return &DynamoDB{
		config: config,
		tracer: otel.Tracer("dynamodb.repository"),
	}
}

// Cria a tabela DynamoDB com os índices secundários globais necessários.
func (p *DynamoDB) Create(ctx context.Context) error {
	ctx, span := p.tracer.Start(ctx, "create", trace.WithAttributes(
		attribute.String("tableName", p.config.Table)),
	)
	defer span.End()
	_, err := p.config.Client.CreateTable(ctx, &dynamodb.CreateTableInput{
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
				AttributeType: types.ScalarAttributeTypeN,
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
						AttributeName: aws.String("statusCode"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("date"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
			},
		},
		TableName:   &p.config.Table,
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		span.RecordError(err)
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return err
	}
	// deve aguardar até a tabela ser criada e estar disponível para uso
	span.AddEvent("waiting for table to be ready")
	waiter := dynamodb.NewTableExistsWaiter(p.config.Client)
	err = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{TableName: &p.config.Table}, 5*time.Minute)
	if err != nil {
		span.RecordError(err)
		return err
	}
	span.AddEvent("table ready")
	// só é possível habilitar TTL na tabela após ela ter sido criada
	_, err = p.config.Client.UpdateTimeToLive(ctx, &dynamodb.UpdateTimeToLiveInput{
		TableName: &p.config.Table,
		TimeToLiveSpecification: &types.TimeToLiveSpecification{
			AttributeName: aws.String("expiration"),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

// Salva o registro na tabela DynamoDB.
// Se já houver registro com o mesmo id, ele será substituído.
func (p *DynamoDB) Save(ctx context.Context, event *models.Event) error {
	ctx, span := p.tracer.Start(ctx, "save", trace.WithAttributes(
		attribute.String("id", event.Id)),
	)
	defer span.End()
	if event.Expiration == 0 {
		event.Expiration = time.Now().Add(p.config.TTL).Unix()
	}
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		span.RecordError(err)
		return err
	}
	_, err = p.config.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &p.config.Table,
		Item:      item,
	})
	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}

// Deleta o registro da tabela DynamoDB pelo id.
func (p *DynamoDB) Delete(ctx context.Context, id string) (event *models.Event, err error) {
	ctx, span := p.tracer.Start(ctx, "delete", trace.WithAttributes(
		attribute.String("id", id)),
	)
	defer span.End()
	out, err := p.config.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &p.config.Table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if out.Attributes == nil {
		span.AddEvent("record not found")
		return nil, nil
	}
	err = attributevalue.UnmarshalMap(out.Attributes, &event)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return event, nil
}

// Recupera o registro da tabela DynamoDB pelo id.
func (p *DynamoDB) Get(ctx context.Context, id string) (event *models.Event, err error) {
	ctx, span := p.tracer.Start(ctx, "get", trace.WithAttributes(
		attribute.String("id", id)),
	)
	defer span.End()
	out, err := p.config.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &p.config.Table,
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if out.Item == nil {
		span.AddEvent("record not found")
		return nil, nil
	}
	err = attributevalue.UnmarshalMap(out.Item, &event)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return event, nil
}

// Procura registros com a data entre o período especificado e com o status code fornecido.
func (p *DynamoDB) FindByDateAndReturnCode(ctx context.Context, from time.Time, to time.Time, statusCode int) (events []*models.Event, err error) {
	ctx, span := p.tracer.Start(ctx, "find-by-date-and-return-code", trace.WithAttributes(
		attribute.String("from", from.Format(time.RFC3339)),
		attribute.String("to", from.Format(time.RFC3339)),
		attribute.Int("statusCode", statusCode),
	))
	defer span.End()
	condition := &dynamodb.QueryInput{
		TableName: aws.String(p.config.Table),
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
	paginator := dynamodb.NewQueryPaginator(p.config.Client, condition)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			span.RecordError(err)
			return nil, err
		}
		for _, item := range page.Items {
			event := &models.Event{}
			err = attributevalue.UnmarshalMap(item, event)
			if err != nil {
				span.RecordError(err)
				return nil, err
			}
			events = append(events, event)
		}
	}
	return events, nil
}
