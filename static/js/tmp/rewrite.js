//	nginx autoindex rewrite
//	sub_filter '</head>' '<script>window.onload=function(){var domain=document.domain;var cdn="xyedit.oss-cn-hongkong.aliyuncs.com";for(var a of document.getElementsByTagName("a")){if(!a.href.endsWith("/")){a.href=a.href.replace(domain,cdn)}}};</script></head>';
//	sub_filter_once     on;
window.onload = function(){
	var domain = document.domain;
	var cdn = "xyedit.oss-cn-hongkong.aliyuncs.com";
	for(var a of document.getElementsByTagName('a')){
		if(!a.href.endsWith('/')){
			a.href=a.href.replace(domain,cdn);
		}
	}
};