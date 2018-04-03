"use strict";
void function(){
    if(window.localStorage){
		if(localStorage.getItem("IndexData")==null){
    		localStorage.setItem("IndexData",JSON.stringify({}))
		}
		if(localStorage.getItem("ContentData")==null){
    		localStorage.setItem("ContentData",JSON.stringify({}))
		}
        //监听本地回话同步
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
    }
}();

