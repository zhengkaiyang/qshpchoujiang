package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net/http"
	"time"
)





var db *sql.DB
var redisclient *redis.Client
func init() {
	db, _ = sql.Open("mysql", "root:123456@/mysql?charset=utf8")
	db.SetMaxOpenConns(2000)
	db.SetMaxIdleConns(1000)
	db.Ping()

	redisclient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "123456", // no password set
		DB:       0,  // use default DB
	})
	pong, err := redisclient.Ping().Result()
	fmt.Println(pong, err)
}

// WXLogin 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况
//func WXLogin(code string) (*WXLoginResp, error) {
//	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
//
//	// 合成url, 这里的appId和secret是在微信公众平台上获取的
//	url = fmt.Sprintf(url, appId, secret, code)
//
//	// 创建http get请求
//	resp,err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer resp.Body.Close()
//
//	// 解析http请求中body 数据到我们定义的结构体中
//	wxResp := WXLoginResp{}
//	decoder := json.NewDecoder(resp.Body)
//	if err := decoder.Decode(&wxResp); err != nil {
//		return nil, err
//	}
//
//	// 判断微信接口返回的是否是一个异常情况
//	if wxResp.ErrCode != 0 {
//		return nil, errors.New(fmt.Sprintf("ErrCode:%s  ErrMsg:%s", wxResp.ErrCode, wxResp.ErrMsg))
//	}
//
//	return &wxResp, nil
//}

// GetRandomString 随机生成字符串
func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

// WeightedRandomIndex 从weights中按照加权抽取一个数,返回抽到的
func WeightedRandomIndex(weights []int) int {
	if len(weights) == 1 {
		return 0
	}
	var sum float32 = 0.0
	for _, w := range weights {
		sum += float32(w)
	}
	r := rand.Float32() * sum
	var t float32 = 0.0
	for i, w := range weights {
		t += float32(w)
		if t > r {
			return i
		}
	}
	return len(weights) - 1
}

func check(err error) bool{
	if err != nil {
		fmt.Println("有些错误")
		return false
	}
	return true
}

// WXlogin1 新添加一个用户在数据表中
func WXlogin1(context *gin.Context){
	uid := context.PostForm("userid")
	//stm2, err := db.Prepare("INSERT drawuser (openid) VALUES ('?')")
	//check(err)
	//last,err := stm2.Exec(uid)
	//check(err)
	//fmt.Println(last.LastInsertId())
	stm,err :=db.Exec("INSERT drawusers (openid) values (?)",uid)
	if check(err){
		fmt.Println(stm.LastInsertId())
	}
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
	}}
	context.SetCookie("choujiang_name", uid, 3600, "/", "localhost", false, true)
	context.JSON(http.StatusOK,m)
}

// GetTheFra 抽取碎片
func GetTheFra(context *gin.Context){
	uid, _ := context.Cookie("choujiang_name")
	//defer db.Close()
	fmt.Println(uid)
	sql := fmt.Sprintf("SELECT chou from drawusers where openid = '%s'",uid)
	rows := db.QueryRow(sql)
	var times int
	err := rows.Scan(&times)
	check(err)
	fmt.Println(times)
	if times<=0{
		m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
			"fragment":[]int{0},
		}}
		context.JSON(http.StatusOK,m)
	}else {
		sql = fmt.Sprintf("update drawusers set chou = chou-1 where openid = '%s'",uid)
		_,err = db.Exec(sql)
		check(err)
		fmt.Println(sql)
		if times==14{
			sui1,sui2 := getwhatsuipian(uid),getwhatsuipian(uid)
			m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
				"fragment":[]int{sui1,sui2},
			}}
			context.JSON(http.StatusOK,m)

		}else {
			sui1 := getwhatsuipian(uid)
			m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
				"fragment":[]int{sui1},
			}}
			context.JSON(http.StatusOK,m)
		}
	}
}

// GetMyFra 查看自己的碎片
func GetMyFra(context *gin.Context){
	uid, _ := context.Cookie("choujiang_name")
	fmt.Println(uid)
	sql := fmt.Sprintf("SELECT sui1,sui2,sui3,sui4,sui5,sui6,sui7,sui8,sui9 from drawusers where openid = '%s'",uid)
	fmt.Println(sql)
	rows := db.QueryRow(sql)
	var fra  = make([]int,9)
	err := rows.Scan(&fra[0],&fra[1],&fra[2],&fra[3],&fra[4],&fra[5],&fra[6],&fra[7],&fra[8])
	check(err)
	fmt.Println(fra)
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
		"sui":fra,
	}}
	context.JSON(http.StatusOK,m)

}

// GiftTime 获得剩余的开奖时间,返回一个标准的时间的格式
func GiftTime(context *gin.Context){
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
		"stilltime":"2021-12-20-11:00:00",
	}}
	context.JSON(http.StatusOK,m)

}

//返回自己中几等奖
func how(context *gin.Context)int{
	uid, _ := context.Cookie("choujiang_name")
	fmt.Println(uid)
	sql := fmt.Sprintf("SELECT thing FROM `luckydogs` where SID = '%s'",uid)
	fmt.Println(sql)
	row := db.QueryRow(sql)
	var giftt int
	row.Scan(&giftt)
	if giftt==0{
		fmt.Println("您没有抽到礼物")
	}else{
		fmt.Println("你的礼物是",giftt)
	}
	return giftt
}

// 查询用户(自己)是否抽到奖，告诉抽到几等奖。还有抽到123等奖的信息
func gift(context *gin.Context){
	//uid, _ := context.Cookie("choujiang_name")
	//fmt.Println(uid)
	//sql := fmt.Sprintf("SELECT thing FROM `luckydogs` where SID = '%s'",uid)
	//fmt.Println(sql)
	//row := db.QueryRow(sql)
	//var giftt int
	//row.Scan(&giftt)
	//if giftt==0{
	//
	//	fmt.Println("您没有抽到礼物")
	//}else{
	//	fmt.Println("你的礼物是",giftt)
	//}
	giftt := how(context)

	var id  string //获奖用户的信息，可能需要用户的头像，还需要修改
	var firstprice []map[string]string // 一等奖的信息
	var secondprice []map[string]string // 二等奖的信息
	var thirdprice []map[string]string // 三等奖的信息

	sql := fmt.Sprintf("SELECT SID FROM `luckydogs` where thing = 1")
	rows,_ := db.Query(sql)
	for rows.Next(){
		rows.Scan(&id)
		temp := make(map[string]string)
		temp["wxNickName"] = id
		temp["picture"] = "https://"
		firstprice = append(firstprice,temp)
	}
	fmt.Println(firstprice)

	sql = fmt.Sprintf("SELECT SID FROM `luckydogs` where thing = 2")
	rows,_ = db.Query(sql)
	for rows.Next(){
		rows.Scan(&id)
		temp := make(map[string]string)
		temp["wxNickName"] = id
		temp["picture"] = "https://"
		secondprice = append(secondprice,temp)
	}
	fmt.Println(secondprice)

	sql = fmt.Sprintf("SELECT SID FROM `luckydogs` where thing = 3")
	rows,_ = db.Query(sql)
	for rows.Next(){
		rows.Scan(&id)
		temp := make(map[string]string)
		temp["wxNickName"] = id
		temp["picture"] = "https://"
		thirdprice = append(thirdprice,temp)
	}
	fmt.Println(thirdprice)


	m := map[string]interface{}{
		"errcode":0,
		"data": map[string]interface{}{
		"myAward":giftt,
		"awardList":map[string]interface{}{
			"level1":firstprice,
			"level2":secondprice,
			"level3":thirdprice,
		},
			},
	}
	context.JSON(200,m )
}

// GiftInform 将个人信息上传到获奖人的数据库(我是中奖人)
func GiftInform(context *gin.Context){
	giftt := how(context)
	realname := context.PostForm("realname")
	xuehao := context.PostForm("xuehao")
	tele := context.PostForm("tele")
	uid, _ := context.Cookie("choujiang_name")
	//首先要判断一个人是否已经存在，存在则更新（不用）
	//sql := fmt.Sprintf("select")
	sql := fmt.Sprintf("INSERT into  dogsinfor (SID,thing,realname,xuehao,tele) VALUES('%s',%d,'%s','%s','%s')",uid,giftt,realname,xuehao,tele)
	_ , err := db.Exec(sql)
	check(err)
	fmt.Println(sql)
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{

	}}
	context.JSON(200,m)
}

// GiftRe 获奖人重新更新到数据库(我是中奖人)
func GiftRe(context *gin.Context){
	realname := context.PostForm("realname")
	xuehao := context.PostForm("xuehao")
	tele := context.PostForm("tele")
	uid, _ := context.Cookie("choujiang_name")
	//首先要判断一个人是否已经存在，存在则更新（不用）
	//sql := fmt.Sprintf("select")
	sql := fmt.Sprintf("update dogsinfor set realname='%s' , xuehao='%s',tele='%s' where SID='%s'",realname,xuehao,tele,uid)
	_ , err := db.Exec(sql)
	check(err)
	fmt.Println(sql)
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{

	}}
	context.JSON(200,m)
}

// GiftQRcode 返回一个二维码所包含的信息,用户扫描这个二维码来获得抽奖机会
func GiftQRcode(context *gin.Context)  {
	var xiaogift int
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(100)<50{xiaogift =1}
	token := GetRandomString(32)
	fmt.Println(token)
	setKey := "choujiangtokens"
	//SADD
	num, err := redisclient.SAdd(setKey, token).Result()
	if err != nil {
		fmt.Printf("Sadd faild,err:%v\n", err)
	}
	fmt.Printf("Sadd %d succ.\n", num)
	fmt.Println(xiaogift)

	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
		"token" : token,
		"gift" : xiaogift,
	}}
	context.JSON(200,m)
}

// GetTheTimes 前端通过调用接口来让用户获得一次的抽奖机会
func GetTheTimes(context *gin.Context){
	token := context.PostForm("token")
	setKey := "choujiangtokens"
	num, err := redisclient.SRem(setKey, token).Result()
	var can int // 是否允许抽奖
	if err == redis.Nil{
		fmt.Println("没有数据")
	}
	if err != nil {
		fmt.Printf("Srem faild,err:%v\n", err)
	}
	if num>=1{
		can = 1
		uid, _ := context.Cookie("choujiang_name")
		stm := fmt.Sprintf("update drawusers set chou=chou+1 WHERE openid = '%s'",uid)
		_,err = db.Exec(stm)
		check(err)
	}
	m := map[string]interface{}{"errcode":0,"data": map[string]interface{}{
		"can":can,
	}}
	context.JSON(200,m)
}

//返回切片最大的值
func mmax(a []int)int{
	var b = a[0]
	for _,v := range a{
		if v>b{
			b = v
		}
	}
	return b
}

// 每次抽碎片调用这个函数,看一个人抽到什么碎片,返回抽到的碎片
func getwhatsuipian(openid string)int{
	stm := db.QueryRow("SELECT sui1,sui2,sui3,sui4,sui5,sui6,sui7,sui8,sui9 FROM `drawusers` where openid=?",openid)
	var jq = make([]int,9) //碎片量
	err := stm.Scan(&jq[0],&jq[1],&jq[2],&jq[3],&jq[4],&jq[5],&jq[6],&jq[7],&jq[8])
	basesui := mmax(jq)
	var jiaquan = make([]int,9) //加权
	for i,_ := range jiaquan{
		jiaquan[i] = basesui-jq[i]
	}

	var jieguo = WeightedRandomIndex(jiaquan)
	fmt.Println(jieguo)
	stm2 := fmt.Sprintf("update drawusers set sui%d=sui%d+1 WHERE openid = '%s'",jieguo+1,jieguo+1,openid)
	_,err = db.Exec(stm2)
	fmt.Println(jiaquan)
	check(err)
	return jieguo+1
}


//返回剩余时间的标准格式(废案)
//func timestamp()int64{
//	//m,_ := time.ParseDuration("+1m")
//	//result :=time.Now().Add(m)
//	//fmt.Println(result)
//	//fmt.Println(result.Format("2006-01-02-15:04:05"))
//	//return result.Unix()
//	m := 312500
//	day := m/86400
//	hour := m/3600%24
//	min := m%3600/60
//	second := m%3600%60
//	stilltime := fmt.Sprintf("%d-%d:%d:%d",day,hour,min,second)
//
//}



func main(){
	route := gin.Default()
	route.GET("/test", func(context *gin.Context) {
		context.SetCookie("choujiang_name", "你好", 3600, "/", "localhost", false, true)
		context.String(200,"asdasd")
	})
	route.POST("/",WXlogin1)
	route.GET("/GetTheFra",GetTheFra)
	route.GET("/GetMyFra",GetMyFra)
	route.GET("/GiftTime",GiftTime)
	route.GET("/gift",gift)
	route.POST("/GiftInform",GiftInform)
	route.POST("/GiftRe",GiftRe)
	route.GET("/GiftQRcode",GiftQRcode)
	route.POST("/GetTheTimes",GetTheTimes)
	route.Run(":8080")


	//getwhatsuipian("what")
}

