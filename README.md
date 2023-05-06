# k8s-platform
K8s管理系统后端: go+gin
项目背景，整体设计，Client-go，框架搭建

一、项目背景
随着容器技术的广泛应用，kubernetes逐渐成为业内的核心技术，是容器编排技术的首选工具。而k8s管理平台在日常的容器维护中也发挥着举足轻重的作用，但随着k8s的定制化功能越来越多，dashboard已经无法满足日常的维护需求，且dashboard的源码学习成本较高，抽象程度较高，二次发成本也就比较高。
本项目使用当下较主流的前端vue+element plus及后端go+gin框架，打造与dashboard对标的k8s管理功能，且可定制化程度高，可根据自身需求，进行灵活定制开发。

二、前后端分离概述
前后端分离已成为互联网项目开发的业界标准使用方式，通过nginx+tomcat的方式( 也可以中间加一个node.js)有效的进行解耦，并且前后端分离会为以后的大型分布式架构、弹性计算架构、微服务架构、多端化服务( 多种客户端，例如:浏览器，车载终端，安卓，IOS等等)打下坚实的基础。这个步要是系统架构从进化成人的必经之路。
 前后分离的优势 :
1.可以实现真正的前后端解耦，前端服务器使用nginx。
2.发现bug，可以快速定位是谁的问题，不会出现互相踢皮球的现象

3.在大并发情况下，可以同时水平扩展前后端服务器
4.增加代码的维护性&易读性(前后端耦在一起的代码读起来相当费劲 )

5.提升开发效率，因为可以前后端并行开发，而不是像以前的强依赖

三、功能设计
 四、client-go介绍
1、简介
client-go是kubernetes官方提供的go语言的客户端库，go应用使用该库可以访问kubernetes的API Server，这样我们就能通过编程来对kubermetes资源进行增删改查操作:
除了提供丰富的API用于操kubernetes资源，client-go还为controller和operator提供了重要支持client-go的informer机制可以将controller关注的资源变化及时带给此controller，使controller能够及时响应变化。
通过client-go提供的客户端对象与kubernetes的API Server进行交豆，而client-go提供了以下四种客户端对象
(1)RESTClient:这是最基础的客户端对象，仅对HTTPRequest进行了封装，实现RESTFul风格API，这个对象的使用并不方便，因为很多参数都要使用者来设置，于是client-go基于RESTClient又实现了三种新的客户端对象;
(2)ClientSet;把Resoure和Version也封装成方法了，用起亲更简单直接，一个资源是一个客户端，多个资源就对应了多个客户端所以ClientSet就是多个客户端的集合了，这样就理解了，不过ClientSet只能访问内置资源，访问不了自定义资源
(3)DynamiClient;可以访问内置资源和自定义资源，个人感觉有点像Java的集合操作，拿出的内容是Object类型，按实际情况自己去做强制转换，当然了也会有强转失败的风险:
(4)DiscoveryClient: 用于发现kubernetes的API Server支持的Group、Version、Resources等信息;

// kubectl api-resources

