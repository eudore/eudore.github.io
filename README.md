# test

主页：https://www.wejass.com



环境配置：centos7+nginx1.13+goalng1.8+MaribDB5.5


需求：

	进行技术测试。
  
实施：

	前端：
  
		html+css+js的常规布局，文档树目录直接从h5本地存储读取数据，没有数据就从服务器缓存一定数量数据到本地。
		
		文档类容处理方法一样，但由于本地空间限制可能会限制文档数量。
		
		使用回话存储保存登录分配session id和token，session id用于每次访问时获得对应权限，认证失败时使用请求临时token再次认证，非法则退出登录。
    
		参考一些其他的开源的富文本编辑器，开发自己的富文本编辑器，通过登录权限来给予是否编辑的功能。
		
    
	web服务器：
  
		使用nginx1.13.0作为web服务器，同时开发nginx过滤模块来进行session控制，调用memcache存储过程进行token认证，将全站非法请求全部过滤,对回话进行等级划分。
    
	后台：
  
		使用golang开发后台，处理nginx接收的请求，主要是与数据交互操作。
		在数据库建立session表保存当前登录登录信息。
   
	数据库：
  
		使用MariaDB5.5.52数据库，建立相关表、触发、存储过程。
		
    
	其他：
  
		使用RSA、ECC双证书。
		研究阿里oss作为CND使用。
		
web

	oauth：nginx模块过滤请求
	
	editor：
	
	file：
	
