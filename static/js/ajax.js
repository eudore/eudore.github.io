function InitAjax(){
	var ajax = false;
	//开始初始化XMLHttpRequest对象
	if(window.XMLHttpRequest) { //Mozilla 浏览器
	    ajax = new XMLHttpRequest();
	    if (ajax.overrideMimeType){ //设置MiME类别
	        ajax.overrideMimeType("text/xml");
		}
	}
	else if (window.ActiveXObject) 
	{     // IE浏览器
	    try {
	        ajax = new ActiveXObject("Msxml2.XMLHTTP");
	    }
	    catch (e) {
			ajax = new ActiveXObject("Microsoft.XMLHTTP");
	    }
	}
	if (!ajax) 
	{     // 异常，创建对象实例失败
	    window.alert("不能创建XMLHttpRequest对象实例.");
	    return null;
	}
	return ajax
}

function runAjax(method,url,sync,data,func){
	var ajax = InitAjax();
	if (ajax != null){
		ajax.open(method, url, sync);
		ajax.setRequestHeader("Content-Type","application/x-www-form-urlencoded");
		ajax.send(data);
		ajax.onreadystatechange = function(){
			if (ajax.readyState == 4 && ajax.status == 200){
				func(ajax.responseText);
			}
		}
	}	
}
