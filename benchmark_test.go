package mdbm_test

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	mdbm "github.com/torden/go-mdbm"
)

var gRandomNo = rand.New(rand.NewSource(time.Now().UnixNano()))

func getBenchMarkRandom() int {

	var key int

	for {
		key = gRandomNo.Intn(1000)
		if key >= 0 {
			break
		}
	}

	return key
}

func Benchmark_boltdb_Store(b *testing.B) {

	os.Remove(pathTestBoltDBBenchmark1)

	db, err := bolt.Open(pathTestBoltDBBenchmark1, 0644, nil)
	if err != nil {
		b.Fatalf("failured, can't open the boltdb, path=%s, err=%v", pathTestBoltDBBenchmark1, err)
	}
	defer db.Close()

	bucketName := []byte("MyBucket")

	i := 0

	for i = 0; i < b.N; i++ {
		key := []byte(strconv.Itoa(i))
		value := []byte(strconv.Itoa(i))

		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucketName)
			if err != nil {
				b.Fatalf("failed, can't data(kv=%d) add to the boltdb(=%s), err=%v", i, pathTestBoltDBBenchmark1, err)
				return err
			}

			err = bucket.Put(key, value)
			if err != nil {
				b.Fatalf("failed, can't data(kv=%d) add to the boltdb(=%s), err=%v", i, pathTestBoltDBBenchmark1, err)
				return err
			}
			return nil
		})

		if err != nil {
			b.Fatalf("exception : boltdb(=%s), err=%v", pathTestBoltDBBenchmark1, err)
		}
	}
}

func Benchmark_mdbm_Store(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMBenchmark1, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects|mdbm.Trunc|mdbm.AnyLocks, 0644, 0, 0)
	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		rv, err := dbm.Store(i, time.Now().UnixNano(), mdbm.Replace)
		if err != nil {
			b.Fatalf("failed, can't data(=%d) add to the mdbm file(=%s), rv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
		}
	}
}

func Benchmark_mdbm_StoreWithLock(b *testing.B) {

	dbm := mdbm.NewMDBM()
	err := dbm.Open(pathTestDBMBenchmark2, mdbm.Create|mdbm.Rdrw|mdbm.LargeObjects|mdbm.Trunc|mdbm.AnyLocks, 0644, 0, 0)
	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		rv, err := dbm.StoreWithLock(i, time.Now().UnixNano(), mdbm.Replace)
		if err != nil {
			b.Fatalf("failed, can't data(=%d) add to the mdbm file(=%s), rv=%d, err=%v", i, dbm.GetDBMFile(), rv, err)
		}
	}
}

func Benchmark_boltdb_Fetch(b *testing.B) {

	db, err := bolt.Open(pathTestBoltDBBenchmark1, 0644, nil)
	if err != nil {
		b.Fatalf("failured, can't open the boltdb, path=%s, err=%v", pathTestBoltDBBenchmark1, err)
	}
	defer db.Close()

	bucketName := []byte("MyBucket")

	for i := 0; i < b.N; i++ {
		key := []byte(strconv.Itoa(getBenchMarkRandom()))

		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bucketName)
			if bucket == nil {
				return fmt.Errorf("Bucket %q not found!", bucketName)
			}

			val := bucket.Get(key)
			if len(val) < 0 || err != nil {
				b.Fatalf("failured, not exists key(=%d), path=%s, err=%v", i, pathTestBoltDBBenchmark1, err)
			}

			return nil
		})
	}
}

func Benchmark_mdbm_Fetch(b *testing.B) {

	var err error
	var val string

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMBenchmark1, 0644)

	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		val, err = dbm.Fetch(getBenchMarkRandom())
		if len(val) < 0 || err != nil {
			b.Fatalf("failured, not exists key(=%d), path=%s, err=%v", i, dbm.GetDBMFile(), err)
		}
	}
}

func Benchmark_mdbm_FetchWithLock(b *testing.B) {

	var err error
	var val string

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	for i := 0; i < b.N; i++ {
		val, err = dbm.FetchWithLock(getBenchMarkRandom())
		if len(val) < 0 || err != nil {
			b.Fatalf("failured, not exists key(=%d), path=%s, err=%v", i, dbm.GetDBMFile(), err)
		}
	}
}

func Benchmark_mdbm_PreLoad_Fetch(b *testing.B) {

	var err error
	var val string

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.PreLoad()
	if err != nil {
		b.Fatalf("failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	for i := 0; i < b.N; i++ {
		val, err = dbm.Fetch(getBenchMarkRandom())
		if len(val) < 0 || err != nil {
			b.Fatalf("failured, not exists key(=%d), path=%s, err=%v", i, dbm.GetDBMFile(), err)
		}
	}
}

func Benchmark_mdbm_PreLoad_FetchWithLock(b *testing.B) {

	var err error
	var val string

	dbm := mdbm.NewMDBM()
	err = dbm.EasyOpen(pathTestDBMBenchmark1, 0644)
	defer dbm.EasyClose()
	if err != nil {
		b.Fatalf("failured, can't open the mdbm, path=%s, err=%v", dbm.GetDBMFile(), err)
	}

	rv, err := dbm.PreLoad()
	if err != nil {
		b.Fatalf("failured, can't pre-load the mdbm, path=%s, rv=%d, err=%v", dbm.GetDBMFile(), rv, err)
	}

	for i := 0; i < b.N; i++ {
		val, err = dbm.FetchWithLock(getBenchMarkRandom())
		if len(val) < 0 || err != nil {
			b.Fatalf("failured, not exists key(=%d), path=%s, err=%v", i, dbm.GetDBMFile(), err)
		}
	}
}
