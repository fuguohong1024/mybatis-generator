package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// 结构体中字段首字母大写外部才能访问
type JavaBean struct {
	Packages     string
	Imports      string
	JavaBeanName string
	Fields       string
	Gets         string
	Sets         string
}

type DaoBean struct {
	Packages     string
	Imports      string
	JavaBeanName string
	Methods      map[string]string
}

type MapperBean struct {
	DaoPackage  string
	PojoBean    string
	Results     string
	Columns     string
	InsertSQL   string
	InsertMap   map[string]string
	UpdateMap   map[string]string
	FindListMap map[string]string
	DeleteMap   map[string]string
}

var myTemplate *template.Template
var err error

func GetMapperFile(objs []map[string]string, table string) {
	// 读取模板
	myTemplate, err = template.ParseFiles(TEMPLATE_PATH + "/mapper.tpl")
	if err != nil {
		fmt.Println(TEMPLATE_PATH+"/mapper.tpl"+"--parse file err:", err)
		return
	}

	javaName := ToJavaName(table)

	_dao := PACKAGE_DAO + "." + javaName + "Dao"
	_pojo := PACKAGE_JAVABEAN + "." + javaName

	// 读取字段
	_results := ""
	_columns := ""
	_columnsPOJO := ""
	_insertSQL := ""
	_javaField := ""

	_ifcolumn := ""
	_ifpojo := ""
	_updateSet := ""
	_updateWhere := ""
	_where := ""
	for _, obj := range objs {
		_javaField = strings.ToLower(obj["column_name"])
		_javaField = ToHumpField(_javaField)
		_results = _results + `		<result column="` + obj["column_name"] + `" property="` + _javaField + `"/>` + "\n"
		_columns = _columns + "		`" + obj["column_name"] + "`,\n"
		_columnsPOJO = _columnsPOJO + "		#{" + _javaField + "},\n"

		_ifcolumn = _ifcolumn + "			<if test=\"" + _javaField + "!=null\">`" + obj["column_name"] + "`,</if> \n"
		_ifpojo = _ifpojo + "			<if test=\"" + _javaField + "!=null\">#{" + _javaField + "},</if> \n"
		_where = _where + "			<if test=\"" + _javaField + " != null\">AND " + obj["column_name"] + " = #{" + _javaField + "}</if>\n"

		// update
		_updateSet = _updateSet + "			<if test=\"" + _javaField + "!=null\">`" + obj["column_name"] + "`= #{" + _javaField + "},</if> \n"
	}

	// 去除最后一个,
	lastIndex := strings.LastIndex(_columns, ",")
	_columns = _columns[0:lastIndex] + strings.Replace(_columns[lastIndex:len(_columns)], ",", "", -1)

	lastIndexPOJO := strings.LastIndex(_columnsPOJO, ",")
	_columnsPOJO = _columnsPOJO[0:lastIndexPOJO] + strings.Replace(_columnsPOJO[lastIndexPOJO:len(_columnsPOJO)], ",", "", -1)

	_insertSQL = "	INSERT INTO " + table + " (\n" +
		_columns + "\n" +
		"	) VALUES ( \n" +
		_columnsPOJO +
		"	) \n"

	insertMap := make(map[string]string)
	insertMap["table"] = table
	insertMap["ifcolumn"] = strings.Trim(_ifcolumn, "		")
	insertMap["ifpojo"] = strings.Trim(_ifpojo, "		")

	updateMap := make(map[string]string)
	updateMap["table"] = table
	updateMap["updateSet"] = strings.Trim(_updateSet, "		")
	updateMap["updateWhere"] = strings.Trim(_updateWhere, "		")

	findListMap := make(map[string]string)
	findListMap["table"] = table
	findListMap["where"] = strings.Trim(_where, "		")

	deleteMap := make(map[string]string)
	deleteMap["table"] = table
	deleteMap["where"] = strings.Trim(_where, "		")

	mapperBean := MapperBean{
		DaoPackage:  strings.Trim(_dao, "		"),
		PojoBean:    strings.Trim(_pojo, "		"),
		Results:     strings.Trim(_results, "		"),
		Columns:     strings.Trim(_columns, "		"),
		InsertSQL:   strings.Trim(_insertSQL, "		"),
		InsertMap:   insertMap,
		UpdateMap:   updateMap,
		FindListMap: findListMap,
		DeleteMap:   deleteMap,
	}

	//myTemplate.Execute(os.Stdout, mapperBean) // 控制台输出
	_path := OUT_PATH + strings.Replace(MAPPER_PATH, ".", "/", -1)
	checkPath(_path)
	file, err := os.OpenFile(_path+"/"+javaName+"Mapper.xml", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}
	myTemplate.Execute(file, mapperBean)

	// --读取字段
	//tpl = strings.Replace(tpl, "#dao#", _dao, -1)
	//tpl = strings.Replace(tpl, "#pojo#", _pojo, -1)
	//tpl = strings.Replace(tpl, "#columns#", _columns, -1)
	//tpl = strings.Replace(tpl, "#results#", _results, -1)
	//tpl = strings.Replace(tpl, "#insertSQL#", _insertSQL, -1)
	//
	//tpl = strings.Replace(tpl, "#table#", table, -1)
	//tpl = strings.Replace(tpl, "#ifcolumn#", _ifcolumn, -1)
	//tpl = strings.Replace(tpl, "#ifpojo#", _ifpojo, -1)
	//tpl = strings.Replace(tpl, "#where#", _where, -1)
	//
	//// 写文件
	//_path := OUT_PATH + strings.Replace(MAPPER_PATH, ".", "/", -1)
	//checkPath(_path)
	//d1 := []byte(tpl)
	//ioutil.WriteFile(_path+"/"+javaName+"Mapper.xml", d1, 0644)

}

func GetDaoFile(table string) {
	// 读取模板
	myTemplate, err = template.ParseFiles(TEMPLATE_PATH + "/dao.tpl")
	if err != nil {
		fmt.Println(TEMPLATE_PATH+"/dao.tpl"+"--parse file err:", err)
		return
	}

	javaName := ToJavaName(table)

	//获取方法
	_imports := ""
	var _methods = make(map[string]string)

	_imports = "import java.util.List;\n"
	_imports = _imports + "import " + PACKAGE_JAVABEAN + "." + javaName + ";\n"
	_imports = _imports + "import org.apache.ibatis.annotations.Mapper;\n"

	_insert := "int insert(" + javaName + " pojo);"
	_insertList := "int insertSelective(List<" + javaName + "> pojo);"
	_findList := "List<" + javaName + "> findList(" + javaName + " pojo);"
	_delete := "int delete(" + javaName + " pojo);"
	_update := "int updateOne(" + javaName + " pojo);"
	//_methods = _insert + "\n" + _insertList + "\n" + _findList + "\n" + _delete + "\n"
	_methods["_insert"] = _insert
	_methods["_insertList"] = _insertList
	_methods["_findList"] = _findList
	_methods["_delete"] = _delete
	_methods["_update"] = _update

	daoBean := DaoBean{
		Packages:     PACKAGE_DAO,
		Imports:      _imports,
		JavaBeanName: javaName + "Dao",
		Methods:      _methods,
	}

	//myTemplate.Execute(os.Stdout, daoBean) // 控制台输出
	_path := OUT_PATH + strings.Replace(PACKAGE_DAO, ".", "/", -1)
	checkPath(_path)
	file, err := os.OpenFile(_path+"/"+javaName+"Dao.java", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}
	myTemplate.Execute(file, daoBean)

	//tpl = strings.Replace(tpl, "#package#", PACKAGE_DAO, 1)
	//tpl = strings.Replace(tpl, "#imports#", _imports, 1)
	//tpl = strings.Replace(tpl, "#javaBeanName#", javaName+"Dao", 1)
	//tpl = strings.Replace(tpl, "#method#", _methods, 1)
	//--获取方法

	// 写JavaBean文件
	//_path := OUT_PATH + strings.Replace(PACKAGE_DAO, ".", "/", -1)
	//checkPath(_path)
	//d1 := []byte(tpl)
	//ioutil.WriteFile(_path+"/"+javaName+"Dao.java", d1, 0644)

}

/**
* 获取JavaBean文件
 */
func GetJavaBean(objs []map[string]string, table string) {

	// 读取模板
	myTemplate, err = template.ParseFiles(TEMPLATE_PATH + "/javabean.tpl")
	if err != nil {
		fmt.Println(TEMPLATE_PATH+"/javabean.tpl"+"parse file err:", err)
		return
	}

	javaName := ToJavaName(table)

	// 获取字段、get、set
	_fileds := ""
	_gets := ""
	_sets := ""
	_imports := ""
	for _, obj := range objs {
		_fileds = _fileds + ToJavaBeanField(obj["column_name"], obj["data_type"]) + "\n"
		_gets = _gets + ToFiledGetMethod(obj["column_name"], obj["data_type"]) + "\n"
		_sets = _sets + ToFiledSetMethod(obj["column_name"], obj["data_type"]) + "\n"

		if !strings.Contains(_imports, relationType[obj["data_type"]]) {
			_imports = _imports + "import " + relationType[obj["data_type"]] + ";\n"
		}
	}
	//--获取字段、get、set

	javaBean := JavaBean{
		Packages:     PACKAGE_JAVABEAN,
		Imports:      _imports,
		JavaBeanName: javaName,
		Fields:       strings.Trim(_fileds, "    "),
		Gets:         strings.Trim(_gets, "    "),
		Sets:         strings.Trim(_sets, "    "),
	}

	//javabeanTPL = strings.Replace(javabeanTPL, "#package#", PACKAGE_JAVABEAN, 1)
	//javabeanTPL = strings.Replace(javabeanTPL, "#javaBeanName#", javaName, 1)
	//javabeanTPL = strings.Replace(javabeanTPL, "#fields#", _fileds, 1)
	//javabeanTPL = strings.Replace(javabeanTPL, "#gets#", _gets, 1)
	//javabeanTPL = strings.Replace(javabeanTPL, "#sets#", _sets, 1)
	//javabeanTPL = strings.Replace(javabeanTPL, "#imports#", _imports, 1)

	// 写JavaBean文件
	//_path := OUT_PATH + strings.Replace(PACKAGE_JAVABEAN, ".", "/", -1)
	//checkPath(_path)
	//d1 := []byte(javabeanTPL)
	//ioutil.WriteFile(_path+"/"+javaName+".java", d1, 0644)

	//myTemplate.Execute(os.Stdout, javaBean) // 控制台输出
	_path := OUT_PATH + strings.Replace(PACKAGE_JAVABEAN, ".", "/", -1)
	checkPath(_path)
	file, err := os.OpenFile(_path+"/"+javaName+".java", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("open failed err:", err)
		return
	}
	myTemplate.Execute(file, javaBean)

}

func ToFiledGetMethod(field string, fieldType string) string {
	_fieldType := relationType[fieldType]
	_fieldType = GetTypeName(_fieldType)
	_javaName := ToJavaName(field)
	_field := ToHumpField(field)
	tpl := `    public ` + _fieldType + ` get` + _javaName + `() {
        return ` + _field + `;
    }`
	return tpl
}

func ToFiledSetMethod(field string, fieldType string) string {
	_fieldType := relationType[fieldType]
	_fieldType = GetTypeName(_fieldType)
	_javaName := ToJavaName(field)
	_field := ToHumpField(field)
	tpl := `    public void set` + _javaName + `(` + _fieldType + " " + _field + `) {
        this.` + _field + ` = ` + _field + `;
    }`
	return tpl
}

/**
 * 转化JavaBean字段
 */
func ToJavaBeanField(field string, fieldType string) string {
	_fieldType := relationType[fieldType]
	_fieldType = GetTypeName(_fieldType)
	_field := ToHumpField(field)
	return "    private " + _fieldType + " " + _field + ";"
}

func GetTypeName(str string) string {
	arr := strings.Split(str, ".")
	lens := len(arr)
	result := arr[lens-1]
	return result
}

/*
*
转成TxxxTxx
*/
func ToJavaName(s string) string {
	arr := strings.Split(s, "_")
	var result string = ""
	for _, str := range arr {
		slen := len(str)
		if slen == 0 {
			continue
		}
		result = result + strings.ToUpper(str[0:1]) + str[1:slen]
	}
	return result
}

/*
*
返回驼峰结构 txxxTxxx
*/
func ToHumpField(field string) string {
	arr := strings.Split(field, "_")
	var result string = ""
	for i, str := range arr {
		if i != 0 {
			slen := len(str)
			result = result + strings.ToUpper(str[0:1]) + str[1:slen]
		} else {
			result = result + str
		}
	}
	return result
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
* 校验路径是否存在，不存在创建
 */
func checkPath(path string) {
	exist, err := PathExists(path)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}
	if exist {
		fmt.Printf("has dir![%v]\n", path)
	} else {
		fmt.Printf("no dir![%v]\n", path)
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir [%v] success!\n", path)
		}
	}
}
