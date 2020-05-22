#!/bin/bash

echo "测试开始..."

echo "注册教师teacher..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/register \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=teacher \
  --data pwd=123 \
  --data role=Teacher

echo -e "\n注册审查者reviewer..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/register \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=reviewer \
  --data pwd=123 \
  --data role=Reviewer

echo -e "\n教师登录..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/login \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=teacher \
  --data pwd=123 \
  --data role=Teacher

echo -e "\n教师发布试题..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/put_question \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=teacher \
  --data id=1 \
  --data question="1+1=?" \
  --data answer=2 \

sleep 2s # 生成区块需要2s

echo -e "\n审查者登录..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/login \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=reviewer \
  --data pwd=123 \
  --data role=Reviewer

echo -e "\n审查者获取待审核事件..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/get_cache \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=reviewer


echo -e "\n审查者批准待审核试题..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/approve \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=reviewer  \
  --data op=put \
  --data id=cache_put_1

echo -e "\n注册学生student..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/register \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=student \
  --data pwd=123 \
  --data role=Student

echo -e "\n学生登录..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/login \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=student \
  --data pwd=123 \
  --data role=Student

echo -e "\n学生获取试题..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/get_question \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=student  \
  --data id=1 

echo -e "\n学生获取所有试题..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/get_all_questions \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=student

echo -e "\n注册管理员admin..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/register \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=admin2 \
  --data pwd=123 \
  --data role=Admin

echo -e "\n管理员登录..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/login \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=admin2 \
  --data pwd=123 \
  --data role=Admin

echo -e "\n管理员获取日志..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/get_logs \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=admin2 

echo -e "\n教师获取所有试题..."
curl --request POST \
  --url http://127.0.0.1:8080/api/v1/get_all_questions \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data name=teacher

echo -e "\n"

