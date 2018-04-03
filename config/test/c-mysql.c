#include <stdio.h>  
#include <stdlib.h>  
#include <mysql.h>  
 
#define DB_SERVER "localhost"  
#define DB_NAME "Jass"  
#define DB_USER "root"  
#define DB_PWD  ""  
 
static MYSQL *db_handel,mysql;  
static MYSQL_ROW row;  
static int query_error;  
 
MYSQL_RES *query_test(char *sql);  
int query_show(MYSQL_RES *result);  
int main(int argc,char *argv[])  
{  
        MYSQL_RES * results;  
        results=query_test("select TYPE from tb_Session WHERE ID='123';");//获取记录  
        query_show(results);//显示记录  
        return 0;  
}  
 
//查询记录  
MYSQL_RES *query_test(char *sql)  
{  
        static MYSQL_RES *query_result;  
        printf("%s\n",sql);  
        mysql_init(&mysql);  
        db_handel=mysql_real_connect(&mysql,DB_SERVER,DB_USER,DB_PWD,DB_NAME,0,0,0);//打开数据库连接  
        if(db_handel==NULL)//错误处理  
        {  
                printf(mysql_error(&mysql));  
                return NULL;  
        }  
 
        query_error=mysql_query(db_handel,"set names utf8");//查询  
        query_error=mysql_query(db_handel,sql);//查询  
        if(query_error!=0)//错误处理  
        {  
                printf(mysql_error(db_handel));  
                return NULL;  
        }  
        query_result=mysql_store_result(db_handel);//获取记录  
        mysql_close(db_handel);//关闭数据库  
        return query_result;//返回记录  
}  
//显示记录  
int query_show(MYSQL_RES *result)  
{  
        unsigned int i,num_fields;  
        MYSQL_FIELD *fileds;  
        num_fields=mysql_num_fields(result);//获取字段数  
        fileds=mysql_fetch_fields(result);//获取字段数组  
        while((row=mysql_fetch_row(result))!=NULL)//循环显示  
        {  
                for(i=0;i<num_fields;i++)  
                {  
                    //正式操作这里通过fwrite写入到文件  
            printf("%s \t",row[i]);  
                }  
          
            printf("\n");  
    }  
        return 0;  
}  