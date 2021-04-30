# Conf

配置特性：

   * 使用`.toml`格式
   * 不支持复杂类型，只支持K-V格式
   * 不区分大小写
   * 默认选择项目根目录下`${project}.toml`文件，可以通过环境变量`CONF_PATH`指定
   * 支持动态加载


原因

   * 和环境变量兼容，所有的配置都可以通过环境变量覆写
   * 一个`repo`下可能有多个`app`，如果使用`yaml`格式，每个`app`都需要定义自己的`config`，然后解析，业务无法配置解耦
   * 配置跟环境变量一样，不区分大小写字母
   
   
使用
```golang
import "${project}/util/conf"

v := conf.Get("THIS_IS_CONFIG")
```
