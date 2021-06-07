# 温湿度传感器demo （kubeedge router）

目的： 云端能下发一个指令收集温湿度数据，把温湿度传感器的数据上传云端。（主要是为了应用router和eventbus部分的代码）

[代码](https://github.com/lcw2/temperatureRouterDemo)

## 背景知识

rest：云端应用

eventbus：负责转发边缘侧应用产生的数据和设备收集的数据

Mqtt : 物联网协议

## 部署crd文件

1. 见参考资料里的router文档：部署RuleEndpoint 和 Rule文档 （RuleEndpoint和Rule的定义见官方文档）

2. 部署项目里的temperature_eventbus.yaml, temperature_rest.yaml

## 从云端到边缘(rest-eventbus)

3. 部署create_rule_rest_eventbus.yaml

   ```yaml
   apiVersion: rules.kubeedge.io/v1
   kind: Rule
   metadata:
     name: operationrule
     labels:
       description: temperation-operation
   spec:
     source: "temperature-rest" 
     sourceResource: {"path":"/operation"}  ()
     target: "temperature-eventbus"
     targetResource: {"topic":"operation"} ()
   ```

   

## 从边缘到云端(eventbus-rest)

4. 部署create_rule_eventbus_rest.yaml

```yaml
apiVersion: rules.kubeedge.io/v1
kind: Rule
metadata:
  name: temperature-rule-eventbus-rest
  labels:
    description: rule-about-collecting-temperature-and0humidity
spec:
  source: "temperature-eventbus"
  sourceResource: {"topic": "device/temperature/and/humidity（边缘应用发布的主题）","node_name": "raspberrypi(自己的节点名称)"}
  target: "temperature-rest"
  targetResource: {"resource":"http://127.0.0.1:8080/environment（云端应用的url）"}
```



## 部署云端应用和边缘应用

docker build -t .....  制作镜像

部署代码里temperature_deployment.yaml 和 httpserver_deployment.yaml

## 部署成功

云端下发一个指令：{"operation":"start"} ，温湿度传感器开始收集温湿度数据并上传

```shell
curl -H "Content-Type:application/json" -X POST -d '{"operation": "start"}' http://192.168.0.4:9443/raspberrypi/default/operation
```

Url地址为 http://{cloudcore_ip}:9443/{node_name}/default/topic

Or 

云端下发一个指令：{"operation":"stop"},温湿度传感器停止收集温湿度数据

```shell
curl -H "Content-Type:application/json" -X POST -d '{"operation": "stop"}' http://192.168.0.4:9443/raspberrypi/default/operation
```



云端查看温湿度数据：

```
curl http://127.0.0.1:8080/environment
```



**note**：（mqtt命令）

mqtt 订阅者 ：发布订阅主题

```shell
mosquitto_sub -t 'topic' -d
```

mqtt 发布者： 发布某topic内容。

```shell
mosquitto_pub -t 'topic' -d -m '{"edgemsg":"msgtocloud"}'
```



## 参考资料

[temperature-demo](https://github.com/kubeedge/examples/tree/master/temperature-demo)

[router文档](https://docs.kubeedge.io/en/docs/developer/custom_message_deliver/)