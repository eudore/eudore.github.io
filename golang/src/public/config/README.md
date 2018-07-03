# config


解析5种配置
* default
* file
* mode
* env
* flag

依次加载默认配置、文件配置、模式配置、环境变量、输入配置，每次加载覆盖之前配置。
输入的文件路径配置和模式配置会提前解析。


文件配置可以是本地文件or远程读取，使用json格式

file://...

http://...

config.Reload(cs ...string) error 

按参数重新加载指定配置，默认全部