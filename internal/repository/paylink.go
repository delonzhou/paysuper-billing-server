package repository

import (
	"context"
	"fmt"
	"github.com/paysuper/paysuper-billing-server/internal/database"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-proto/go/billingpb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	mongodb "gopkg.in/paysuper/paysuper-database-mongo.v2"
)

const (
	collectionPaylinks = "paylinks"

	cacheKeyPaylink         = "paylink:id:%s"
	cacheKeyPaylinkMerchant = "paylink:id:%s:merchant_id:%s"
)

type paylinkRepository repository

// NewPaylinkRepository create and return an object for working with the paylink repository.
// The returned object implements the PaylinkRepositoryInterface interface.
func NewPaylinkRepository(db mongodb.SourceInterface, cache database.CacheInterface) PaylinkRepositoryInterface {
	s := &paylinkRepository{db: db, cache: cache}
	return s
}

func (r *paylinkRepository) Insert(ctx context.Context, pl *billingpb.Paylink) error {
	oid, err := primitive.ObjectIDFromHex(pl.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldQuery, pl.Id),
		)
	}

	filter := bson.M{"_id": oid}
	opts := options.FindOneAndReplace().SetUpsert(true)
	err = r.db.Collection(collectionPaylinks).FindOneAndReplace(ctx, filter, pl, opts).Err()

	if err != nil {
		zap.S().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationInsert),
			zap.Any(pkg.ErrorDatabaseFieldDocument, pl),
		)
	}

	_ = r.updateCaches(pl)

	return nil
}

func (r *paylinkRepository) Update(ctx context.Context, pl *billingpb.Paylink) error {
	pl.IsExpired = pl.IsPaylinkExpired()

	plMgo, err := pl.GetMgoPaylink()

	if err != nil {
		zap.S().Error(
			pkg.MethodFinishedWithError,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldDocument, pl),
		)
		return err
	}

	query := bson.M{"_id": plMgo.Id}

	set := bson.M{"$set": bson.M{
		"expires_at":     plMgo.ExpiresAt,
		"updated_at":     plMgo.UpdatedAt,
		"products":       plMgo.Products,
		"name":           plMgo.Name,
		"no_expiry_date": plMgo.NoExpiryDate,
		"products_type":  plMgo.ProductsType,
		"is_expired":     plMgo.IsExpired,
	}}
	_, err = r.db.Collection(collectionPaylinks).UpdateOne(ctx, query, set)

	if err != nil {
		zap.S().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationUpdate),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSet, set),
		)
		return err
	}

	_ = r.updateCaches(pl)

	return nil
}

func (r *paylinkRepository) Delete(ctx context.Context, obj *billingpb.Paylink) error {
	obj.Deleted = true

	oid, err := primitive.ObjectIDFromHex(obj.Id)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldQuery, obj.Id),
		)
	}

	filter := bson.M{"_id": oid}
	_, err = r.db.Collection(collectionPaylinks).ReplaceOne(ctx, filter, obj)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationDelete),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, filter),
		)
		return err
	}

	_ = r.updateCaches(obj)

	return nil
}

func (r *paylinkRepository) UpdateTotalStat(ctx context.Context, pl *billingpb.Paylink) error {
	plMgo, err := pl.GetMgoPaylink()

	if err != nil {
		zap.S().Error(
			pkg.MethodFinishedWithError,
			zap.Error(err),
			zap.Any(pkg.ErrorDatabaseFieldDocument, pl),
		)
		return err
	}

	query := bson.M{"_id": plMgo.Id}

	set := bson.M{"$set": bson.M{
		"visits":                plMgo.Visits,
		"conversion":            plMgo.Conversion,
		"total_transactions":    plMgo.TotalTransactions,
		"sales_count":           plMgo.SalesCount,
		"returns_count":         plMgo.ReturnsCount,
		"gross_sales_amount":    plMgo.GrossSalesAmount,
		"gross_returns_amount":  plMgo.GrossReturnsAmount,
		"gross_total_amount":    plMgo.GrossTotalAmount,
		"transactions_currency": plMgo.TransactionsCurrency,
		"is_expired":            plMgo.IsExpired,
	}}
	_, err = r.db.Collection(collectionPaylinks).UpdateOne(ctx, query, set)

	if err != nil {
		zap.S().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
			zap.Any(pkg.ErrorDatabaseFieldSet, set),
		)
		return err
	}

	_ = r.updateCaches(pl)

	return nil
}

func (r *paylinkRepository) FindCount(ctx context.Context, merchantId string, projectId string) (n int64, err error) {
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return int64(0), err
	}

	query := bson.M{
		"deleted":     false,
		"merchant_id": merchantOid,
	}

	if projectId != "" {
		query["project_id"], err = primitive.ObjectIDFromHex(projectId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
				zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
			)
			return int64(0), err
		}
	}

	n, err = r.db.Collection(collectionPaylinks).CountDocuments(ctx, query)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldOperation, pkg.ErrorDatabaseFieldOperationCount),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return int64(0), nil
	}

	return n, nil
}

func (r *paylinkRepository) Find(ctx context.Context, merchantId string, projectId string, limit, offset int64) (result []*billingpb.Paylink, err error) {
	merchantOid, err := primitive.ObjectIDFromHex(merchantId)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseInvalidObjectId,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.String(pkg.ErrorDatabaseFieldQuery, merchantId),
		)
		return nil, err
	}

	query := bson.M{
		"deleted":     false,
		"merchant_id": merchantOid,
	}

	if projectId != "" {
		query["project_id"], err = primitive.ObjectIDFromHex(projectId)

		if err != nil {
			zap.L().Error(
				pkg.ErrorDatabaseInvalidObjectId,
				zap.Error(err),
				zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
				zap.String(pkg.ErrorDatabaseFieldQuery, projectId),
			)
			return nil, err
		}
	}

	if limit <= 0 {
		limit = pkg.DatabaseRequestDefaultLimit
	}

	if offset <= 0 {
		offset = 0
	}

	opts := options.Find().
		SetSort(bson.M{"_id": 1}).
		SetLimit(limit).
		SetSkip(offset)
	cursor, err := r.db.Collection(collectionPaylinks).Find(ctx, query, opts)

	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return
	}

	err = cursor.All(ctx, &result)

	if err != nil {
		zap.L().Error(
			pkg.ErrorQueryCursorExecutionFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, query),
		)
		return
	}

	return
}

func (r *paylinkRepository) GetById(ctx context.Context, id string) (pl *billingpb.Paylink, err error) {
	key := fmt.Sprintf(cacheKeyPaylink, id)
	oid, _ := primitive.ObjectIDFromHex(id)
	dbQuery := bson.M{"_id": oid, "deleted": false}
	return r.getBy(ctx, key, dbQuery)
}

func (r *paylinkRepository) GetByIdAndMerchant(ctx context.Context, id, merchantId string) (pl *billingpb.Paylink, err error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	merchantOid, _ := primitive.ObjectIDFromHex(merchantId)
	key := fmt.Sprintf(cacheKeyPaylinkMerchant, id, merchantId)
	dbQuery := bson.M{"_id": oid, "merchant_id": merchantOid, "deleted": false}
	return r.getBy(ctx, key, dbQuery)
}

func (r *paylinkRepository) getBy(ctx context.Context, key string, dbQuery bson.M) (pl *billingpb.Paylink, err error) {
	if err = r.cache.Get(key, &pl); err == nil {
		pl.IsExpired = pl.GetIsExpired()
		return pl, nil
	}

	err = r.db.Collection(collectionPaylinks).FindOne(ctx, dbQuery).Decode(&pl)
	if err != nil {
		zap.L().Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldCollection, collectionPaylinks),
			zap.Any(pkg.ErrorDatabaseFieldQuery, dbQuery),
		)
		return
	}

	pl.IsExpired = pl.GetIsExpired()

	if pl.Deleted == false {
		err = r.updateCaches(pl)
	}
	return
}

func (r *paylinkRepository) updateCaches(pl *billingpb.Paylink) (err error) {
	key1 := fmt.Sprintf(cacheKeyPaylink, pl.Id)
	key2 := fmt.Sprintf(cacheKeyPaylinkMerchant, pl.Id, pl.MerchantId)

	if pl.Deleted {
		err = r.cache.Delete(key1)
		if err != nil {
			return
		}

		err = r.cache.Delete(key2)
		return
	}

	err = r.cache.Set(key1, pl, 0)
	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key1),
			zap.Any(pkg.ErrorCacheFieldData, pl),
		)
		return
	}

	err = r.cache.Set(key2, pl, 0)
	if err != nil {
		zap.L().Error(
			pkg.ErrorCacheQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorCacheFieldCmd, "SET"),
			zap.String(pkg.ErrorCacheFieldKey, key2),
			zap.Any(pkg.ErrorCacheFieldData, pl),
		)
	}
	return
}
