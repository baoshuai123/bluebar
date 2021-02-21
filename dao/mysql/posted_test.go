package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"testing"
)

func init() {
	dbcfg := settings.MysqlConfig{
		Host:         "localhost",
		Port:         3306,
		User:         "root",
		Password:     "root",
		Dbname:       "bluebell",
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbcfg)
	if err != nil {
		panic(err)
	}
}
func TestCreatePost(t *testing.T) {
	post := models.Post{
		ID:          10,
		AuthorID:    121,
		CommunityID: 1,
		Title:       "test",
		Content:     "test test test",
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
	//type args struct {
	//	p *models.Post
	//}
	//tests := []struct {
	//	name    string
	//	args    args
	//	wantErr bool
	//}{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if err := CreatePost(tt.args.p); (err != nil) != tt.wantErr {
	//			t.Errorf("CreatePost() error = %v, wantErr %v", err, tt.wantErr)
	//		}
	//	})
	//}
}
