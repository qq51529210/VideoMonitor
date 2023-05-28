// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/records": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "云端录像"
                ],
                "summary": "列表",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "大于创建时间戳，比较",
                        "name": "afterCreateTime",
                        "in": "query"
                    },
                    {
                        "maxLength": 64,
                        "type": "string",
                        "description": "app ，精确",
                        "name": "app",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "小于删除时间戳 ，比较",
                        "name": "beforeDeleteTime",
                        "in": "query"
                    },
                    {
                        "minimum": 1,
                        "type": "integer",
                        "description": "条数，小于 1 不匹配",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "创建时间戳 ，精确",
                        "name": "createTime",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "删除时间戳 ，精确",
                        "name": "deleteTime",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "number",
                        "description": "时长，单位秒 ，精确",
                        "name": "duration",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "description": "软删除 ，精确",
                        "name": "isDeleted",
                        "in": "query"
                    },
                    {
                        "enum": [
                            0,
                            1
                        ],
                        "type": "integer",
                        "description": "是否在录像时间内的文件 ，精确",
                        "name": "isRecording",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "偏移，小于 0 不匹配",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "排序，\"column [desc]\"",
                        "name": "order",
                        "in": "query"
                    },
                    {
                        "minimum": 0,
                        "type": "integer",
                        "description": "大小，字节 ，精确",
                        "name": "size",
                        "in": "query"
                    },
                    {
                        "maxLength": 64,
                        "type": "string",
                        "description": "stream ，精确",
                        "name": "stream",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/util.GORMList-db_Record"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "云端录像"
                ],
                "summary": "添加",
                "parameters": [
                    {
                        "description": "数据",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/records.postReq"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/internal.IDResult-string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "云端录像"
                ],
                "summary": "批量删除",
                "parameters": [
                    {
                        "description": "条件",
                        "name": "ids",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        },
        "/records/{id}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "云端录像"
                ],
                "summary": "详情",
                "parameters": [
                    {
                        "type": "string",
                        "description": "主键",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/db.Record"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            },
            "delete": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "云端录像"
                ],
                "summary": "删除",
                "parameters": [
                    {
                        "type": "string",
                        "description": "主键",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/internal.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "db.Record": {
            "type": "object",
            "properties": {
                "app": {
                    "description": "app",
                    "type": "string"
                },
                "createTime": {
                    "description": "创建时间戳",
                    "type": "integer"
                },
                "deleteTime": {
                    "description": "删除时间戳",
                    "type": "integer"
                },
                "duration": {
                    "description": "时长，单位秒",
                    "type": "number"
                },
                "id": {
                    "description": "名称",
                    "type": "string"
                },
                "isDeleted": {
                    "description": "软删除",
                    "type": "integer"
                },
                "isRecording": {
                    "description": "是否在录像时间内的文件",
                    "type": "integer"
                },
                "size": {
                    "description": "大小，字节",
                    "type": "integer"
                },
                "stream": {
                    "description": "stream",
                    "type": "string"
                }
            }
        },
        "internal.Error": {
            "type": "object",
            "properties": {
                "detail": {
                    "description": "详细信息"
                },
                "phrase": {
                    "description": "简短信息",
                    "type": "string"
                }
            }
        },
        "internal.IDResult-string": {
            "type": "object",
            "properties": {
                "id": {
                    "description": "数据库 ID",
                    "type": "string"
                }
            }
        },
        "records.postReq": {
            "type": "object",
            "required": [
                "app",
                "createTime",
                "duration",
                "name",
                "saveDays",
                "size",
                "stream"
            ],
            "properties": {
                "app": {
                    "description": "app",
                    "type": "string",
                    "maxLength": 64
                },
                "createTime": {
                    "description": "创建时间",
                    "type": "integer",
                    "minimum": 1
                },
                "duration": {
                    "description": "时长",
                    "type": "number",
                    "minimum": 1
                },
                "isRecording": {
                    "description": "是否在录像时间内",
                    "type": "boolean"
                },
                "name": {
                    "description": "minio 的标识",
                    "type": "string",
                    "maxLength": 40
                },
                "saveDays": {
                    "description": "保存天数",
                    "type": "integer",
                    "minimum": 0
                },
                "size": {
                    "description": "大小",
                    "type": "integer",
                    "minimum": 1
                },
                "stream": {
                    "description": "stream",
                    "type": "string",
                    "maxLength": 64
                }
            }
        },
        "util.GORMList-db_Record": {
            "type": "object",
            "properties": {
                "data": {
                    "description": "列表",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/db.Record"
                    }
                },
                "total": {
                    "description": "总数",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "接口文档",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
