"use strict";
function getCookie(name){
	var arr,reg=new RegExp("(^| )"+name+"=([^;]*)(;|$)");
	if(arr=document.cookie.match(reg))
		return unescape(arr[2]);
	else
		return null;
}
function setCookie(name,value,expiredays){
	var exp = new Date();
	exp.setDate(exp.getDate()+expiredays)
	document.cookie=name+"="+escape(value)+((expiredays==null) ? "" : ";expires="+exp.toGMTString())
}
function setLocal(name,value){
	setCookie(name,value,2592000000)//30*24*60*60*1000
}
function delCookie(name){
	var exp = new Date();
	exp.setTime(exp.getTime() - 1);
	document.cookie= name + "=;expires="+exp.toGMTString();
}
function init(){
	if(window.localStorage){
		window.addEventListener("storage",function(event){
			if(event.key =="getsessionStorage"){
				localStorage.setItem("sessionStorage",JSON.stringify(sessionStorage));
				localStorage.removeItem("sessionStorage");
			}else if(event.key=="sessionStorage"&& !sessionStorage.length){
				var data=JSON.parse(event.newValue);
				for(var key in data)
					sessionStorage.setItem(key,data[key]);
			}
		});
	}else{
		localStorage.setItem=setLocal;
		localStorage.getItem=getCookie;
		localStorage.removeItem=delCookie;
		sessionStorage.setItem=setCookie;
		sessionStorage.getItem=getCookie;
		sessionStorage.removeItem=delCookie;
	}
}
init()