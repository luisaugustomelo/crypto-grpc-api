package databases

import (
	"errors"
	"klever/grpc/databases/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"golang.org/x/net/context"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	Find(context.Context, interface{}) (CursorResultHelper, error)
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filter interface{}) (int64, error)
	FindOneAndUpdate(context.Context, interface{}, interface{}) SingleResultHelper
}

type SingleResultHelper interface {
	Decode(v interface{}) error
	Err() error
}

type CursorResultHelper interface {
	Close(context.Context) error
	Next(context.Context) bool
	Decode(v interface{}) error
	Err() error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect(context.Context) error
	StartSession() (mongo.Session, error)
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursorResult struct {
	cr *mongo.Cursor
}

type mongoSession struct {
	mongo.Session
}

func NewSecureClient(cnf *config.Config) (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().SetAuth(
		options.Credential{
			Username:   cnf.Username,
			Password:   cnf.Password,
			AuthSource: cnf.DatabaseName,
		}).ApplyURI(cnf.URL))

	return &mongoClient{cl: c}, err
}

func NewClient(cnf *config.Config) (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(cnf.URL))

	return &mongoClient{cl: c}, err
}

func Database(cnf *config.Config, client ClientHelper) DatabaseHelper {
	return client.Database(cnf.DatabaseName)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (CursorResultHelper, error) {
	cursorResult, err := mc.coll.Find(ctx, filter)
	return &mongoCursorResult{cr: cursorResult}, err
}

func (mc *mongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOneAndUpdate(ctx, filter, update)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	id, err := mc.coll.InsertOne(ctx, document)
	return id.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (int64, error) {
	count, err := mc.coll.DeleteOne(ctx, filter)
	return count.DeletedCount, err
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (cr *mongoCursorResult) Decode(v interface{}) error {
	return cr.cr.Decode(v)
}

func (cr *mongoCursorResult) Close(ctx context.Context) error {
	return cr.cr.Close(ctx)
}

func (cr *mongoCursorResult) Next(ctx context.Context) bool {
	return cr.cr.Next(ctx)
}

func (cr *mongoCursorResult) Err() error {
	return cr.cr.Err()
}

func (sr *mongoSingleResult) Err() error {
	return sr.sr.Err()
}

var ErrNoDocuments = errors.New("mongo: no documents in result")
