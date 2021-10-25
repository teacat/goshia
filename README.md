# Goshia

Goshia 是採用 [go-gorm/gorm](https://github.com/go-gorm/gorm/) 與 [teacat/rushia](https://github.com/teacat/rushia/) 的私有套件，這能夠以 Gorm（資料庫連線）作為 Rushia（SQL 語法建置函式庫）的執行基底。

## 使用方式

打開終端機並且透過 `go get` 安裝此套件即可，這個套件的版本通常會與 Rushia 相同。

```bash
$ go get github.com/teacat/goshia/v3
```

### 初始化

使用 Goshia 之前，需要先初始化一個 Gorm 資料庫連線並最後將其帶入至 `goshia.New(db)`。

```go
// 初始化 Gorm 至資料庫的連線。
db, err := gorm.Open(mysql.Open(dsn))
if err != nil {
    log.Fatalf(err)
}
// 將 Gorm 帶入給 Goshia 並初始化一個輔助函式。
goshia := goshia.New(db)
```

### 查詢（Query）

這是最常使用的函式，查詢（如：Select 語法）可以讓 Goshia 將資料直接映射到指定的指針。

```go
var user User
err := goshia.Query(rushia.NewQuery("Users").Where("user_id = ?", ColumnUserID, 10).Select(), &user)
// 等效於：SELECT * FROM Users WHERE `user_id` = ?
```

### 查詢與計數（QueryCount）

這與 `Query` 的用法相同，但會私底下額外執行一個相同但沒有筆數限制的 SQL 去取得計數筆數。

這個使用時機通常是：希望進行分頁並且需要知道總筆數。由於執行 `LIMIT 1, 10` 諸如此類篩選限制 SQL 會導致沒辦法取得真正的總筆數，

因此需要透過 `QueryCount` 自動額外取得不受限制時的總筆數資料。

```go
var users []User
count, err := goshia.Query(rushia.NewQuery("Users").Limit(10, 20).Select(), &users)
// 等效於：SELECT * FROM Users LIMIT 10, 20
// 等效於：SELECT COUNT(*) FROM Users
```

### 執行（Exec）

執行 `DELETE` 或 `UPDATE`…等語法時，可以使用 `Exec` 執行並且取得影響的筆數為何。

```go
var user User
affectedRows, err := goshia.Query(rushia.NewQuery("Users").Where("user_id = ?", 30).Update(user))
// 等效於：UPDATE Users SET ... WHERE `user_id` = ?
```

### 執行與編號（ExecID）

如果欄位帶有自動遞增且希望在 `INSERT` 插入時取得這個新的遞增值，就可以使用 `ExecID` 來取得這個新的編號。

```go
var user User
id, affectedRows, err := goshia.Query(rushia.NewQuery("Users").Insert(user))
// 等效於：INSERT INTO Users SET ...
```

### 交易與回溯（Transaction）

若希望初始化一個事務交易（Transaction），正如其名可以直接呼叫 `Transaction` 並在裡面執行 SQL 語法。

如果該處理函式回傳任何 `error` 則會導致 Goshia 自動回溯所有行為。

如果一個 Goshia 不是由 `Transaction` 被建立的而卻呼叫其事務交易函式會直接觸發 `panic`。

除此之外，Goshia 也提供使用 Gorm 的 `Rollback`、`RollbackTo`、`SavePoint`…等函式。

```go
goshia.Transaction(func (tx *goshia.Goshia) error {
    _, err := tx.Exec(rushia.NewQuery("Users").Where("user_id = ?", 30).Update(user))
    // 等效於：UPDATE Users SET ... WHERE `user_id` = ?
    if err != nil {
        return err
    }
    // 蹦蹦！記得要 Commit！
    tx.Commit()
    return nil
})
```

### 輔助資料型態

透過 `goshia.Int`、`goshia.Bool` 等函式，可以將一個值轉換成指針作為填補 SQL 的 Nullable 欄位。這個用法正如 [AWS SDK](https://docs.aws.amazon.com/sdk-for-go/api/aws/) 一樣。
