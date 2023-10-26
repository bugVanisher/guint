package advanced

import (
	"github.com/bugVanisher/gunit"
	"math/rand"
	"time"
)

type cleanup func()

type FancyFixture struct {
	Base
	mysqlClient *DB
	CleanUps    []cleanup
	streamId    uint64
}

type Base struct {
	*gunit.Fixture
}

type DB struct {
	Base
}

func InitDB(b Base) *DB {
	return &DB{
		b,
	}
}

func (d *DB) Close() bool {
	d.GetLogger().Info().Msg("mysql close")
	return true
}

func (d *DB) Exec(sql string) (any interface{}) {
	d.GetLogger().Info().Msgf("exec sql:%s", sql)
	return any
}

func (f *FancyFixture) GetMysqlClient() *DB {
	if f.mysqlClient == nil {
		f.mysqlClient = InitDB(f.Base)
	}
	return f.mysqlClient
}

//StartLiving means go live-streaming
func (f *FancyFixture) StartLiving() {
	f.streamId = rand.Uint64()
	f.GetLogger().Info().Msgf("start living now:%d", f.streamId)

	// register cleanup method and then do it together
	f.CleanUps = append(f.CleanUps, f.EndLiving)
}

func (f *FancyFixture) EndLiving() {
	f.GetLogger().Info().Msgf("end living now:%d", f.streamId)
}

//Release close resource or cleanup
func (f *FancyFixture) Release() {
	for _, c := range f.CleanUps {
		c()
	}
	f.mysqlClient.Close()
}

func (f *FancyFixture) Streaming() {
	f.GetLogger().Info().Msg("start streaming")

	//launch a blocking task with a goroutine
	go func() {
		time.Sleep(200 * time.Millisecond)
		//task exits with unknown reason
		f.GetLogger().Error().Msg("streaming stop...")
		f.FatalStop()
	}()
}
