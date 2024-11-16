package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
	"log"
	"strings"
)

type Config struct {
	db              string
	table           string
	packageJavaBean string
	packageService  string
	packageDao      string
	packageMapper   string
}

var (
	DATA_SOURCE_NAME = "root:root@tcp(127.0.0.1:3306)/mountain?charset=utf8"
	CONFIG_FILE      = "mysql-config.ini"
	TEMPLATE_PATH    = "./mysql/template"
	OUT_PATH         = "./mysql/out/"
	DB               = "mountain"
	TABLE            = "sys_role"
	PACKAGE_JAVABEAN = "com.site.mountain.entity"
	PACKAGE_SERVICE  = "com.site.mountain.service"
	PACKAGE_DAO      = "com.site.mountain.dao.test2"
	MAPPER_PATH      = "com.site.mountain.dao.mapper"
)

var relationType = make(map[string]string)

func initRelationType() {
	relationType["bigint"] = "java.math.BigInteger"
	relationType["bit"] = "java.lang.Boolean"
	relationType["blob"] = "java.lang.byte[]"
	relationType["char"] = "java.lang.String"
	relationType["date"] = "java.sql.Date"
	relationType["datetime"] = "java.sql.Timestamp"
	relationType["decimal"] = "java.math.BigDecimal"
	relationType["double"] = "java.lang.Double"
	relationType["double precision"] = "java.lang.Double"
	relationType["enum"] = "java.lang.String"
	relationType["float"] = "java.lang.Float"
	relationType["int"] = "java.lang.Integer"
	relationType["integer"] = "java.lang.Long"
	relationType["longblob"] = "java.lang.byte[]"
	relationType["longtext"] = "java.lang.String"
	relationType["mediumblob"] = "java.lang.byte[]"
	relationType["mediumint"] = "java.lang.Integer"
	relationType["mediumtext"] = "java.lang.String"
	relationType["set"] = "java.lang.String"
	relationType["smallint"] = "java.lang.Integer"
	relationType["text"] = "java.lang.String"
	relationType["time"] = "java.sql.Time"
	relationType["timestamp"] = "java.sql.Timestamp"
	relationType["tinyblob"] = "java.lang.byte[]"
	relationType["tinyint"] = "java.lang.Integer"
	relationType["tinytext"] = "java.lang.String"
	relationType["varchar"] = "java.lang.String"
	relationType["year"] = "java.sql.Date"
	//扩展
	relationType["list"] = "java.util.List"
}

func initConfig() {
	c, err := config.ReadDefault(CONFIG_FILE)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}
	TEMPLATE_PATH, _ = c.String("template", "TEMPLATE_PATH")
	OUT_PATH, _ = c.String("template", "OUT_PATH")
	DATA_SOURCE_NAME, _ = c.String("mysql", "DATA_SOURCE_NAME")
	DB, _ = c.String("mysql", "DB")
	TABLE, _ = c.String("mysql", "TABLE")
	PACKAGE_JAVABEAN, _ = c.String("package", "PACKAGE_JAVABEAN")
	PACKAGE_DAO, _ = c.String("package", "PACKAGE_DAO")
	MAPPER_PATH, _ = c.String("package", "MAPPER_PATH")

}

func main() {
	initConfig()
	initRelationType()

	db, err := sql.Open("mysql", DATA_SOURCE_NAME)
	if err != nil {
		panic(err)
	}

	fmt.Println("******************")
	fmt.Println("DATA_SOURCE_NAME=", DATA_SOURCE_NAME)
	fmt.Println("TEMPLATE_PATH=", TEMPLATE_PATH)
	fmt.Println("PACKAGE_JAVABEAN=", PACKAGE_JAVABEAN)
	fmt.Println("PACKAGE_DAO=", PACKAGE_DAO)
	fmt.Println("MAPPER_PATH=", MAPPER_PATH)
	fmt.Println("TABLE=", TABLE)
	fmt.Println("DB=", DB)
	fmt.Println("******************")
	var tables []string

	tables, err = getAllTables(db)
	for _, table := range tables {
		goMapperTools(db, table)
	}

	defer db.Close()
}

func getAllTables(db *sql.DB) ([]string, error) {
	var tables []string
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		tables = append(tables, tableName)
	}

	// 检查是否有错误
	if err = rows.Err(); err != nil {
		return tables, err
	}
	return tables, nil
}

func goMapperTools(db *sql.DB, table string) {
	rows, err := db.Query("select column_name,column_comment,data_type " +
		"from information_schema.columns " +
		"where table_name='" + table + "' and table_schema='" + DB + "'")
	if err != nil {
		panic(err)
	}
	var result []map[string]string
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		// 创建一个map来存储当前行的数据
		row := make(map[string]string)

		// 创建一个切片来存储列的值
		valuePtrs := make([]any, len(columns))
		for i, _ := range valuePtrs {
			valuePtrs[i] = new(string)
		}

		// 扫描当前行的值
		if err := rows.Scan(valuePtrs...); err != nil {
			panic(err)
		}

		// 将值存入map
		for i, col := range columns {
			row[strings.ToLower(col)] = *(valuePtrs[i].(*string))
		}

		// 将map添加到结果切片中
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
	//
	GetJavaBean(result, table)
	GetDaoFile(table)
	GetMapperFile(result, table)

	defer rows.Close()
}
