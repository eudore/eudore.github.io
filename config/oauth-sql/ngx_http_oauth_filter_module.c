/*ngx_http_oauth_module.c*/
#include <stdio.h>
#include <mysql.h>

#include <ngx_config.h>
#include <ngx_core.h>
#include <ngx_http.h>

typedef struct
{
    ngx_flag_t type;
}ngx_http_oauth_conf_t;

typedef struct
{
    ngx_int_t oauth;
}ngx_http_oauth_ctx_t;

/*用static修饰只在本文件生效，因此允许所有的过滤模块都有自己的这两个指针*/
static ngx_http_output_header_filter_pt ngx_http_next_header_filter;
static ngx_http_output_body_filter_pt    ngx_http_next_body_filter;
static ngx_int_t ngx_http_oauth_init(ngx_conf_t *cf);
static ngx_int_t ngx_http_oauth_header_filter(ngx_http_request_t *r);
static ngx_int_t ngx_http_oauth_body_filter(ngx_http_request_t *r, ngx_chain_t *in);
static void* ngx_http_oauth_create_conf(ngx_conf_t *cf);
static char* ngx_http_oauth_merge_conf(ngx_conf_t *cf,void*parent,void*child);
static ngx_flag_t ngx_http_oauth_DB(ngx_str_t *token,ngx_int_t type);


/*处理感兴趣的配置项*/
static ngx_command_t ngx_http_oauth_commands[]=
{
    {
        ngx_string("oauth"), //配置项名称
        NGX_HTTP_MAIN_CONF | NGX_HTTP_SRV_CONF | NGX_HTTP_LOC_CONF | NGX_HTTP_LMT_CONF | NGX_CONF_FLAG,//配置项只能携带一个参数并且是on或者off
        ngx_conf_set_num_slot,//使用nginx自带方法,参数on/off
        NGX_HTTP_LOC_CONF_OFFSET,//使用create_loc_conf方法产生的结构体来存储
        //解析出来的配置项参数
        offsetof(ngx_http_oauth_conf_t, type),//on/off
        NULL
    },
    ngx_null_command //
};


/*模块上下文*/
static ngx_http_module_t ngx_http_oauth_module_ctx=
{
    NULL,                       /* preconfiguration方法  */
    ngx_http_oauth_init,        /* postconfiguration方法 */

    NULL,                       /*create_main_conf 方法 */
    NULL,                       /* init_main_conf方法 */

    NULL,                       /* create_srv_conf方法 */
    NULL,                       /* merge_srv_conf方法 */

    ngx_http_oauth_create_conf, /* create_loc_conf方法 */
    ngx_http_oauth_merge_conf   /*merge_loc_conf方法*/
};

/*定义过滤模块,ngx_module_t结构体实例化*/
ngx_module_t ngx_http_oauth_module =
{
    NGX_MODULE_V1,                 /*Macro*/
    &ngx_http_oauth_module_ctx,         /*module context*/
    ngx_http_oauth_commands,            /*module directives*/
    NGX_HTTP_MODULE,                       /* module type */
    NULL,                                  /* init master */
    NULL,                                  /* init module */
    NULL,                                  /* init process */
    NULL,                                  /* init thread */
    NULL,                                  /* exit thread */
    NULL,                                  /* exit process */
    NULL,                                  /* exit master */
    NGX_MODULE_V1_PADDING                  /*Macro*/
};

static void* ngx_http_oauth_create_conf(ngx_conf_t *cf)
{
    ngx_http_oauth_conf_t  *mycf;

    //创建存储配置项的结构体
    mycf = (ngx_http_oauth_conf_t  *)ngx_pcalloc(cf->pool, sizeof(ngx_http_oauth_conf_t));
    if (mycf == NULL)
    {
        return NULL;
    }

    //ngx_flat_t类型的变量，如果使用预设函数ngx_conf_set_flag_slot
    //解析配置项参数，必须初始化为NGX_CONF_UNSET
    mycf->type = NGX_CONF_UNSET;
    return mycf;
}

static char* ngx_http_oauth_merge_conf(ngx_conf_t *cf,void*parent,void*child)
{
    ngx_http_oauth_conf_t *prev = (ngx_http_oauth_conf_t *)parent;
    ngx_http_oauth_conf_t *conf = (ngx_http_oauth_conf_t *)child;

    //合并ngx_flat_t类型的配置项type
    ngx_conf_merge_value(conf->type, prev->type, 0);

    return NGX_CONF_OK;

}

/*初始化方法*/
static ngx_int_t ngx_http_oauth_init(ngx_conf_t*cf)
{
    //插入到头部处理方法链表的首部
    ngx_http_next_header_filter=ngx_http_top_header_filter;
    ngx_http_top_header_filter=ngx_http_oauth_header_filter;
    ngx_http_next_body_filter=ngx_http_top_body_filter;
    ngx_http_top_body_filter=ngx_http_oauth_body_filter;
    return NGX_OK;
}

/*头部处理方法*/
static ngx_int_t ngx_http_oauth_header_filter(ngx_http_request_t *r)
{
    ngx_http_oauth_ctx_t *ctx;
    ngx_http_oauth_conf_t *conf;
    //处理响应码非200的情形
    if (r->headers_out.status != NGX_HTTP_OK){
        return ngx_http_next_header_filter(r);
    }



    /*获取http上下文*/
    ctx = ngx_http_get_module_ctx(r, ngx_http_oauth_module);
    if(ctx){
        //该请求的上下文已经存在，这说明
        // ngx_http_oauth_header_filter已经被调用过1次，
        //直接交由下一个过滤模块处理
        return ngx_http_next_header_filter(r);
    }


    //获取存储配置项参数的结构体
    conf = ngx_http_get_module_loc_conf(r, ngx_http_oauth_module);
    //如果type成员为0，也就是配置文件中没有配置oauth配置项，
    //或者oauth配置项的参数值是off，这时直接交由下一个过滤模块处理
    if (conf->type == 0){
        return ngx_http_next_header_filter(r);
    }  


    ngx_str_t name = ngx_string("token");
    ngx_str_t token = ngx_null_string;
#ifdef NGX_DEBUG
    ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0, "type: \"%i\"", conf->type); 
#endif
    if(NGX_OK == ngx_http_parse_multi_header_lines(&r->headers_in.cookies, &name, &token)){
#ifdef NGX_DEBUG
        ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0, "cook %V: \"%V\"",&name, &token); 
#endif
        if(ngx_http_oauth_DB(&token,conf->type)){
            return ngx_http_next_header_filter(r);
        }
    }
	if(NGX_OK == ngx_http_arg(r,name.data,name.len,&token)) {
#ifdef NGX_DEBUG
        ngx_log_error(NGX_LOG_DEBUG, r->connection->log, 0, "args %V: \"%V\"",&name, &token); 
#endif
        if(ngx_http_oauth_DB(&token,conf->type)){
			return ngx_http_next_header_filter(r);
        }

    }else{
        return NGX_HTTP_CLOSE;
    }

    return NGX_HTTP_FORBIDDEN;
}

static ngx_int_t ngx_http_oauth_body_filter(ngx_http_request_t *r, ngx_chain_t *in){
    ngx_http_oauth_ctx_t *ctx;
    ctx = ngx_http_get_module_ctx(r, ngx_http_oauth_module);
    if(ctx){
        return ngx_http_next_body_filter(r, in);
    }
    return ngx_http_next_body_filter(r, in);
}

//查询记录  
static ngx_flag_t ngx_http_oauth_DB(ngx_str_t *token,ngx_int_t type){
    static MYSQL *db_handel,mysql;
    static int query_error;
    static MYSQL_RES *query_result;
    static MYSQL_ROW row;
    u_char sql[80]={'\0'};
    ngx_snprintf(sql,sizeof(sql),"call pro_Oauth('%V',%i);",token,type);
    mysql_init(&mysql);
    db_handel=mysql_real_connect(&mysql,"localhost","root","","Jass",0,0,0);//打开数据库连接  
    if(db_handel==NULL){
        return 0;
    }

    query_error=mysql_query(db_handel,"set names utf8");//查询  
    query_error=mysql_query(db_handel,(char *)sql);//查询  
    if(query_error!=0){
        return 0;
    }
    query_result=mysql_store_result(db_handel);//获取记录  
    mysql_close(db_handel);//关闭数据库  
    row=mysql_fetch_row(query_result);
    return (ngx_flag_t)row[0];
}