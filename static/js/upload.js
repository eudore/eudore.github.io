(function (global, factory) {
    typeof exports === 'object' && typeof module !== 'undefined' ? module.exports = factory() :
    typeof define === 'function' && define.amd ? define(factory) :
    (global.Upload = factory());
}(this, (function () { 'use strict';

var Ready,QUEUED,WAIT,UPLOADING,DONE,FAILED

function Upjob(m,obj){
    this.m=m
	this.file=obj
}
Upjob.prototype = {
	target: "oss",
    expire: 0,
	upfile: function(t){
        var up=eval(t+"up()")
	},
	localup: function(){
		var data = new FormData();
        if(this.file.length==0)
            return
        for(var i of this.file){
            data.append('file',i)
        }
        fetch(location.pathname,{ method: 'POST',credentials: 'include', body: data}).then(function(response) {
            if (!response.ok) throw new Error(response.statusText)
            document.getElementById('myframe').contentWindow.location.reload(true);
            return response.json()
        }).then(function(data){
            //console.log(data)
        }).catch(function(err) {
            console.log(err)
        })
	},
    osspolicy: function(call){
        var now = Date.parse(new Date())/1000
        if(this.expire < now+3){  
            fetch("http://47.52.173.119:8081/file/policy",{method: 'GET'}).then(
                function(response) {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            }).then(function(obj){
                console.log(obj)
                this.host = obj['host']
                this.policyBase64 = obj['policy']
                this.accessid = obj['accessid']
                this.signature = obj['signature']
                this.expire = parseInt(obj['expire'])
                this.callbackbody = obj['callback'] 
                this.key = obj['dir']
                call()
            }).catch(function(err) {
                console.log(err)
            })
        }else{
            call()
        }
    },
    ossup: function(){
        this.osspolicy(function(){
            if(this.file.size < 512<<20){
                this.ossoneup()
            }else{
                this.ossmultiup()
            }
        })
    },
	ossoneup: function(){
        var data = new FormData();
        data.append("name","s")
        data.append("key",this.key+this.file.name)
        data.append("policy",this.policyBase64)
        data.append("OSSAccessKeyId",this.accessid)
        data.append("success_action_status",200)
        data.append("callback",this.callbackbody)
        data.append("signature",this.signature)
        data.append("file",this.file)
        console.log(data)
        fetch(this.host,{ method: 'POST', body: data}).then(function(response) {
            if (!response.ok) throw new Error(response.statusText)
            return response.text()
        }).then(function(data){
            console.log(data)
        }).catch(function(err) {
            console.log(err)
        })
	},
    ossmultiup: function(){
        
    }
}

function Upload() {
}
Upload.prototype = {
	queue: new Array(),
	thread: 3,
	start: function(){
		if(this.thread<3){
			this.thread++
		}
	},
	add: function(obj){
		if(obj instanceof FileList || obj instanceof File){
			queue.append(new Upjob(this,obj))
		}else if(obj instanceof HTMLInputElement && obj.type=="file"){
			queue.append(new Upjob(this,obj.files))
		}
	},
    osspolicy: function(call){
        var now = Date.parse(new Date())/1000
        if(this.expire < now+3){  
            fetch("http://47.52.173.119:8081/file/policy",{method: 'GET'}).then(
                function(response) {
                if (!response.ok) throw new Error(response.statusText)
                return response.json()
            }).then(function(obj){
                console.log(obj)
                this.host = obj['host']
                this.policyBase64 = obj['policy']
                this.accessid = obj['accessid']
                this.signature = obj['signature']
                this.expire = parseInt(obj['expire'])
                this.callbackbody = obj['callback'] 
                this.key = obj['dir']
                call()
            }).catch(function(err) {
                console.log(err)
            })
        }else{
            call()
        }
    }
}






function timer(t){
  return new Promise(resolve=>setTimeout(resolve, t))
  .then(function(res) {
    console.log('timeout');
  });
}

function consume(reader) {
  var total = 0
  return new Promise((resolve, reject) => {
    function pump() {
      reader.read().then(({done, value}) => {
        if (done) {
          resolve();
          return;
        }
        total += value.byteLength;
        console.log(`received ${value.byteLength} bytes (${total} bytes in total)`);
        pump();
      }).catch(reject)
    }
    pump();
  });
}


return window.Upload || Upload
})));