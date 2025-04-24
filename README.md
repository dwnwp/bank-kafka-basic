# bank-kafka-basic
event เปิด/ปิด บัญชี ฝาก/ถอน ลง kafka consumer บันทึกลง db<br /> <br />

include:
``` cmd
go get github.com/IBM/sarama
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
go get go.mongodb.org/mongo-driver/v2/mongo
```
<br />
สร้าง topics ใน kafka "OpenAccountEvent", "CloseAccountEvent", "DepositFundEvent", "WithdrawFundEvent"
