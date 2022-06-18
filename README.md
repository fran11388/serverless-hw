##Build lambda

Compile your executable.
```bash
GOOS=linux go build main.go
```
Create a deployment package by packaging the executable in a .zip file.
```bash
zip function.zip main
```
###Reference
[Deploy .zip file archives](https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html)

##Lambda function
####Create function
```bash
aws lambda create-function --function-name hello-dynamo --zip-file fileb://function.zip --handler main --runtime go1.x --role arn:aws:iam::990204874157:role/lambda-ex
```

####Configure the event source
```bash
aws lambda create-event-source-mapping --function-name sqs-error-consumer  --batch-size 10  --event-source-arn arn:aws:sqs:ap-northeast-1:990204874157:consumer-error-queue
```

####Update function
```bash
 aws lambda update-function-code --function-name  kinesis-consumer  --zip-file fileb://function.zip
```


###lambda role
"Arn": "arn:aws:iam::990204874157:role/lambda-ex",

###建立執行角色
```bash
aws iam create-role --role-name lambda-ex --assume-role-policy-document '{"Version": "2012-10-17","Statement": [{ "Effect": "Allow", "Principal": {"Service": "lambda.amazonaws.com"}, "Action": "sts:AssumeRole"}]}'
```

###使用 attach-policy-to-role 命令將許可新增至角色。透過新增 AWSLambdaBasicExecutionRole 受管政策開始。
```bash
aws iam attach-role-policy --role-name lambda-ex --policy-arn arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole
```

