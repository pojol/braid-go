


### v1.2.8
* 添加mailbox中的并行接口 PubAsync()
 > 原有的接口目前作为串行发布接口，并可以获取到错误返回，新的并行接口将用于(go - chan - go)模型上，避免消费端阻塞调度器。