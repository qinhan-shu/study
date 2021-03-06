package s

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

// Config : 全局配置
type Config struct {
	PvpServerAddrs   []string `json:"pvp_server_addrs" env:"PVP_SERVER_ADDRS"`
	MatchServerAddrs []string `json:"match_server_addrs" env:"MATCH_SERVER_ADDRS"`
	ConsulAddr       string   `json:"consul_addr" env:"CONSUL_ADDRESS"`

	GM struct {
		URL      string `json:"url" env:"GM_REPO_URL" required:"game,pvp"`
		UserName string `json:"username" env:"GM_REPO_USERNAME" required:"game,pvp"`
		TOKEN    string `json:"token" env:"GM_REPO_TOKEN" required:"game,pvp"`
	} `json:"gm_repo"`

	Redis struct {
		Addrs    []string `json:"addrs" env:"REDIS_ADDR" required:"all"`
		Password string   `json:"password" env:"REDIS_PASS" required:"all" allowempty:"true"`
		MaxConn  int      `json:"maxconn" env:"REDIS_MAXCONN" envDefault:"100" required:"all"`
	} `json:"redis"`

	MySQL struct {
		Addr                   string `json:"addr" env:"MYSQL_ADDR" required:"game,pvp"`
		Username               string `json:"username" env:"MYSQL_USER" required:"game,pvp"`
		Password               string `json:"password" env:"MYSQL_PASS" required:"game,pvp"`
		Database               string `json:"database" env:"MYSQL_DATABASE" required:"game,pvp"`
		ShardingSize           int    `json:"sharding_size" env:"MYSQL_SHARDING_SIZE" required:"game,pvp"`
		MaxConn                int    `json:"maxconn" env:"MYSQL_MAXCONN" envDefault:"100" required:"game,pvp"`
		LoadAllPlayerOnStartup bool   `json:"load_all_player_on_startup" env:"MYSQL_LOAD_ALL_PLAYER_ON_STARTUP" required:"game,pvp"`
	} `json:"mysql"`

	MongoDB struct {
		Addrs         []string `json:"addrs" env:"MONGODB_ADDR"`
		Username      string   `json:"username" env:"MONGODB_USERNAME"`
		Password      string   `json:"password" env:"MONGODB_PASSWORD" `
		Database      string   `json:"database" env:"MONGODB_DATABASE"`
		MaxConn       uint16   `json:"maxconn" env:"MONGODB_MAXCONN" envDefault:"100"`
		ReplicaSet    string   `json:"replica_set" env:"MONGODB_REPLICASET"`
		AuthMechanism string   `json:"auth_mechanism" env:"MONGODB_AUTHMECHANISM" envDefault:"SCRAM-SHA-1"`
	} `json:"mongodb"`

	SMTP struct {
		Enable          bool     `json:"enable" env:"SMTP_ENABLE"`
		Addr            string   `json:"addr" env:"SMTP_ADDR"`
		Host            string   `json:"host" env:"SMTP_HOST"`
		Username        string   `json:"username" env:"SMTP_USERNAME"`
		Password        string   `json:"password" env:"SMTP_PASSWORD"`
		MailTo          []string `json:"to" env:"SMTP_MAILTO"`
		MailFrom        string   `json:"from" env:"SMTP_MAILFROM" envDefault:"sausageshooter@xindong.com"`
		ContentType     string   `json:"content_type" env:"SMTP_CONTENT_TYPE" envDefault:"Content-Type: text/plain; charset=UTF-8"`
		MessageTemplate string   `json:"message_template" env:"SMTP_MSG_TEMPLATE" envDefault:"To:%s\r\nFrom:<%s>\r\nSubject:Sausage Panic Log\r\n%s\r\n\r\n%s"`
	} `json:"smtp"`
}

// CheckStruct : 通过反射，然后拿到一些结构体中自定义的tag（ allowempty、required 等）
func (c Config) CheckStruct() {
	doCheck(reflect.ValueOf(c), "")
}

func doCheck(refValue reflect.Value, parentName string) {
	refType := refValue.Type()
	for i := 0; i < refType.NumField(); i++ {
		// 拿到具体field的value & type
		refField := refValue.Field(i)
		refFieldType := refType.Field(i)

		if refField.Kind() == reflect.Ptr && !refField.IsNil() {
			// 因为上面的Config中不会有指针，因此这儿直接continue即可
			continue
		}
		if refField.Kind() == reflect.Struct {
			// 如果是struct，递归check
			doCheck(refField, refFieldType.Name)
		}

		// 拿到自定义的字段，这儿以 allowempty 、 required 为例子
		allowEmpty := refFieldType.Tag.Get("allowempty") == "true"
		required := strings.Split(refFieldType.Tag.Get("required"), ",")

		fieldName := refFieldType.Name
		if parentName != "" {
			fmt.Printf("%s.%s", parentName, refFieldType.Name)
		}
		fmt.Printf("allowEmpty[%v], required[%v] ---> %s \n", allowEmpty, required, fieldName)

		err := fmt.Errorf("missing config `%s`", fieldName)
		// 在这根据tag的字段来做配置的检查
		switch refField.Kind() {
		case reflect.Slice:
			if refField.IsNil() {
				log.Fatal(err)
			}
			for j := 0; j < refField.Len(); j++ {
				item := refField.Index(j)
				if item.Kind() != reflect.String {
					break
				}
				if item.CanSet() {
					item.SetString(strings.TrimSpace(item.String()))
				}
			}
		case reflect.String:
			if !allowEmpty && refField.String() == "" {
				log.Fatal(err)
			}
		case reflect.Int, reflect.Int64:
			if !allowEmpty && refField.Int() == 0 {
				log.Fatal(err)
			}
		}
	}
}

// Show : 展示
func (c Config) Show() {
	fmt.Println("PvpServerAddrs   		: ", c.PvpServerAddrs)
	fmt.Println("MatchServerAddrs 		: ", c.MatchServerAddrs)
	fmt.Println("ConsulAddr       		: ", c.ConsulAddr)
	// gm
	fmt.Println("GM_URL      			: ", c.GM.URL)
	fmt.Println("GM_UserName 			: ", c.GM.UserName)
	fmt.Println("GM_TOKEN    			: ", c.GM.TOKEN)
	// redis
	fmt.Println("Redis_Addrs    			: ", c.Redis.Addrs)
	fmt.Println("Redis_Password 			: ", c.Redis.Password)
	fmt.Println("Redis_MaxConn  			: ", c.Redis.MaxConn)
	// mysql
	fmt.Println("MySQL_Addr                   	: ", c.MySQL.Addr)
	fmt.Println("MySQL_Username               	: ", c.MySQL.Username)
	fmt.Println("MySQL_Password               	: ", c.MySQL.Password)
	fmt.Println("MySQL_Database               	: ", c.MySQL.Database)
	fmt.Println("MySQL_ShardingSize           	: ", c.MySQL.ShardingSize)
	fmt.Println("MySQL_MaxConn                	: ", c.MySQL.MaxConn)
	fmt.Println("MySQL_LoadAllPlayerOnStartup 	: ", c.MySQL.LoadAllPlayerOnStartup)
	// mongo
	fmt.Println("MongoDB_Addrs          		: ", c.MongoDB.Addrs)
	fmt.Println("MongoDB_Username       		: ", c.MongoDB.Username)
	fmt.Println("MongoDB_Password       		: ", c.MongoDB.Password)
	fmt.Println("MongoDB_Database       		: ", c.MongoDB.Database)
	fmt.Println("MongoDB_MaxConn        		: ", c.MongoDB.MaxConn)
	fmt.Println("MongoDB_ReplicaSet     		: ", c.MongoDB.ReplicaSet)
	fmt.Println("MongoDB_AuthMechanism  		: ", c.MongoDB.AuthMechanism)
	// smtp
	fmt.Println("SMTP_Enable            		: ", c.SMTP.Enable)
	fmt.Println("SMTP_Addr              		: ", c.SMTP.Addr)
	fmt.Println("SMTP_Host              		: ", c.SMTP.Host)
	fmt.Println("SMTP_Username          		: ", c.SMTP.Username)
	fmt.Println("SMTP_Password          		: ", c.SMTP.Password)
	fmt.Println("SMTP_MailTo            		: ", c.SMTP.MailTo)
	fmt.Println("SMTP_MailFrom          		: ", c.SMTP.MailFrom)
	fmt.Println("SMTP_ContentType       		: ", c.SMTP.ContentType)
	fmt.Println("SMTP_MessageTemplate   		: ", c.SMTP.MessageTemplate)
}
