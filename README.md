# beego
1.安装包

    go get github.com/astaxie/beego

    go get gopkg.in/mgo.v2

    go get github.com/garyburd/redigo/redis

    go get github.com/dgrijalva/jwt-go

2.启动
    bee run -gendoc=true -downdoc=true
    
3.测试
    http://127.0.0.1:8080/v1/auth?uid=1
    http://127.0.0.1:8080/v1/user/2?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjF9.S4Qy_CDXYfWVvN9IHowtVHyemPaff3yjnqNfTVe-BVw
