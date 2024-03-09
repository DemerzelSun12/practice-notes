[TOC]

# 1.如何创建pod或者创建deployment流程？(必问)


## 一、如何创建pod流程：

1. 客户端提交创建请求，可以通过API Server的Restful API，也可以使用如何命令行工具。支持的数据类型包括JSON和YAML。
2. API Server处理用户请求，存储Pod数据到etcd。
3. 调度器通过API Server查看未绑定的Pod。尝试为Pod分配主机。
4. 过滤主机 (调度预选)：调度器用一组规则过滤掉不符合要求的主机。比如Pod指定了所需要的资源量，那么可用资源比Pod需要的资源量少的主机会被过滤掉。
5. 主机打分(调度优选)：对第一步筛选出的符合要求的主机进行打分，在主机打分阶段，调度器会考虑一些整体优化策略，比如把容一个Replication Controller的副本分布到不同的主机上，使用最低负载的主机等。
6. 选择主机：选择打分最高的主机，进行binding操作，结果存储到etcd中。
7. kubelet根据调度结果执行Pod创建操作：绑定成功后，scheduler会调用APIServer的API在etcd中创建一个boundpod对象，描述在一个工作节点上绑定运行的所有pod信息。运行在每个工作节点上的kubelet也会定期与etcd同步boundpod信息，一旦发现应该在该工作节点上运行的boundpod对象没有更新，则调用Docker API创建并启动pod内的容器。


## 二、如何创建deployment流程？（tips：这个要回答出创建rs的过程）

1. 准备好一个包含应用程序的Deployment的yml文件，然后通过kubectl客户端工具发送给ApiServer。
2. ApiServer接收到客户端的请求并将资源内容存储到数据库(etcd)中。
3. Controller组件(包括scheduler、replication、endpoint)监控资源变化并作出反应。
4. ReplicaSet检查数据库变化，创建期望数量的pod实例。
5. Scheduler再次检查数据库变化，发现尚未被分配到具体执行节点(node)的Pod，然后根据一组相关规则将pod分配到可以运行它们的节点上，并更新数据库，记录pod分配情况。
6. Kubelete监控数据库变化，管理后续pod的生命周期，发现被分配到它所在的节点上运行的那些pod。如果找到新pod，则会在该节点上运行这个新pod。
7. kuberproxy运行在集群各个主机上，管理网络通信，如服务发现、负载均衡。例如当有数据发送到主机时，将其路由到正确的pod或容器。对于从主机上发出的数据，它可以基于请求地址发现远程服务器，并将数据正确路由，在某些情况下会使用轮训调度算法(Round-robin)将请求发送到集群中的多个实例。
8. kubectl提交一个请求，来创建RC，此时Controller Manager通过API server里的接口监听到这个RC事件，分析之后，发现当前集群中还没有它对应的Pod实例，于是根据RC里的Pod模板定义Pod对象；接下来，此事件被Scheduler发现，它立即执行一个复杂的调度流程，为这个新Pod选定一个落户的Node，这个过程可称为绑定；随后模板Node上运行的Kubelet进程通过API Server监测到这个“新生的”Pod并按照它的定义，启动Pod并负责后期的管理；随后我们通过Kubectl提交一个映射到该Pod的Server的创建请求，Controller Manager会通过Label标签查询到相关联的Pod实例，然后生成Service的Endpoints信息；接下来，所有Node上运行的Proxy进程通过API Server查询并监听Service对象及其对应的Endpoints信息，建立一个负载均衡器来实现Service访问到后端Pod的流量转发功能；



# 2. k8s中pod之间是怎么通信的？

1. 同一个Pod内的多个容器之间：lo（通过localhost回环地址通信）。
2. 各Pod之间的通讯：overlay Network （覆盖网络）。
3. Pod与 Service之间的通讯：（使用各节点的 Iptables规则）。

## 3. k8s中的pod中OOM机制有了解过吗，原理是什么？

（tips：蚂蚁金服的oceanbase数据库产品问了这个问题，当时被问还是有点点蒙的。这个涉及到pod的QOS的相关知识。阿里面试官首先介绍了系统中进程的OOM机制，会给每个进程进行打分，然后kill掉得分最高的那一个进程，然后问pod的OOM是怎么样的？）

答：当系统 OOM上时，对于处理不同 OOMScore 的进程表现不同，OOMScore 是针对 memory 的，当宿主上 memory 不足时系统会优先 kill 掉 OOMScore 值高的进程，可以使用如下指令：
```bash
  $ cat /proc/$PID/oom_score
```
查看进程的 OOMScore。OOMScore 的取值范围为 [-1000, 1000]，Guaranteed pod 的默认值为 -998，Burstable pod 的值为 2~999，BestEffort pod 的值为 1000，也就是说当系统 OOM 时，首先会 kill 掉 BestEffort pod 的进程，若系统依然处于 OOM 状态，然后才会 kill 掉  Burstable pod，最后是 Guaranteed pod；

## 4. k8s的基本组件及其作用？
一个kubernetes集群主要是由控制节点(master)、工作节点(node)构成，每个节点上都会安装不同的组件。

- master：集群的控制平面，负责集群的决策 ( 管理 )。

  - ApiServer :资源操作的唯一入口，接收用户输入的命令，提供认证、授权、API注册和发现等机制。

  - Scheduler : 负责集群资源调度，按照预定的调度策略将Pod调度到相应的node节点上。

  - ControllerManager :负责维护集群的状态，比如程序部署安排、故障检测、自动扩展、滚动更新等。

  - Etcd ：负责存储集群中各种资源对象的信息，k/v方式存储，所有的 k8s 集群数据存放在此。

  - Kuberctl: 命令行配置工具。

- node：集群的数据平面，负责为容器提供运行环境 ( 干活 )。

  - Kubelet :负责维护容器的生命周期，即通过控制docker，来创建、更新、销毁容器，会按固定频率检查节点健康状态并上报给 APIServer，该状态会记录在 Node 对象的 status 中。

  - KubeProxy :负责提供集群内部的服务发现和负载均衡。

  - Docker :负责节点上容器的各种操作。



## 5.kubevirt的基本组件及其作用？
- virt-api：kubevirt是以CRD形式去管理VM Pod，virt-api就是所有虚拟化操作的入口，这里面包括常规的CDR更新验证、以及console、vm start、stop等操作。

- virt-controller：virt-controller会根据vmi CRD，生成对应的virt-launcher Pod，并且维护CRD的状态。与kubernetes api-server通讯监控VMI资源的创建删除等状态。

- virt-handler：virt-handler会以deamonset形式部署在每一个节点上，负责监控节点上的每个虚拟机实例状态变化，一旦检测到状态的变化，会进行响应并且确保相应的操作能够达到所需（理想）的状态。virt-handler还会保持集群级别VMI Spec与相应libvirt域之间的同步；报告libvirt域状态和集群Spec的变化；调用以节点为中心的插件以满足VMI Spec定义的网络和存储要求。

- virt-launcher：每个virt-launcher，pod对应着一个VMI，kubelet只负责virt-launcher pod运行状态，不会去关心VMI创建情况。virt-handler会根据CRD参数配置去通知virt-launcher去使用本地的libvirtd实例来启动VMI，随着Pod的生命周期结束，virt-lanuncher也会去通知VMI去执行终止操作；其次在每个virt-launcher pod中还对应着一个libvirtd，virt-launcher通过libvirtd去管理VM的生命周期，这样做到去中心化，不再是以前的虚拟机那套做法，一个libvirtd去管理多个VM。

- virtctl：virtctl是kubevirt自带类似kubectl的命令行工具，它是越过virt-launcher pod这一层去直接管理VM虚拟机，可以控制VM的start、stop、restart。



## 6.Etcd的写数据流程是什么？

- 总体上的请求流程从上至下依次为客户端 → API 接口层 → etcd Server → etcd raft 算法库。）
- 读请求：客户端通过负载均衡选择一个 etcd 节点发出读请求，API 接口层提供了 Range RPC 方法，etcd 服务端拦截到 gRPC 读请求后，调用相应的处理器处理请求。
- 写请求：客户端通过负载均衡选择一个 etcd 节点发起写请求，etcd 服务端拦截到 gRPC 写请求，涉及一些校验和监控，之后 KVServer 向 raft 模块发起提案，内容即为写入数据的命令。经过网络转发，当集群中的多数节点达成一致并持久化数据后，状态变更且 MVCC 模块执行提案内容。



## 7.在异构架构下，比如有些主机是x86，有些主机是arm架构，k8s怎么保证异构架构下这些主机拉起不同的架构的镜像？

docker每一个镜像包含了一个文件，这个文件包含了有关于镜像信息，如层、大小和摘要。docker manifest命令还向用户提供附加信息，比如构建镜像的操作系统和体系结构。而manifest list是一个镜像清单列表，用于存放多个不同os/arch的镜像信息。我们主要用到manifest的目的，其实还是多用于存放不同的os/arch信息，也就是方便我们在不同的CPU架构（arm或者x86）或者操作系统中，通过一个镜像名称拉取对应架构或者操作系统的镜像，这个尤其是在K8S中，对于异构CPU的服务器中的镜像显得尤为有效。


## 9.Prometheus相关问题

pro的四种数据构：Counter，Gauge，Histogram，SummaryPrometheus
组成及架构 Prometheus 生态圈中包含了多个组件，其中许多组件是可选择的：
- 1.Prometheus Server: ------服务端 ---处理，储存数据
负责收集和存储时间，Prometheus 组成及架构序列数据（time series data），并且提供查询接口。
- 2.Jobs/Exporters: ------客户端 ---采集数据

客户端，用于暴露已有的第三方服务的 metrics 给 Prometheus。
监控并采集指标，对外暴露HTTP服务（/metrics）；目前已经有很多的软件原生就支持Prometheus，提供/metrics，可以直接使用；对于像操作系统已经不提供/metrics的应用，可以使用现有的exporters或者开发自己的exporters来提供/metrics服务；

- 3.Push Gateway: -----相当于代理---转发数据

针对push系统设计，Short-lived jobs定时将指标push到Pushgateway，再由Prometheus Server从Pushgateway上pull；主要用于短期的 jobs。由于这类 jobs 存在时间较短，可能在 Prometheus 来 pull 之前就消失了。为此，这次 jobs 可以直接向 Prometheus server 端推送它们的 metrics。这种方式主要用于服务层面的 metrics，对于机器层面的 metrices，需要使用 node exporter。

- 4.Alertmanager: ----告警方式---实现告警

报警组件，从 Prometheus server 端接收到 alerts 后，会进行去除重复数据，分组，并路由到对的接受方式，发出报警。常见的接收方式有：电子邮件，pagerduty，OpsGenie, webhook 等。

- 5.Web UI：

Prometheus内置一个简单的Web控制台，可以查询指标，查看配置信息或者Service Discovery等，实际工作中，查看指标或者创建仪表盘通常使用Grafana，Prometheus作为Grafana的数据源；

