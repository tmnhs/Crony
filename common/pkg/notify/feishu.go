package notify

var feiShuTemplateCard = `{
  "msg_type": "interactive",
  "card": {
    "config": {
      "wide_screen_mode": true
    },
    "header": {
      "title": {
        "tag": "plain_text",
        "content": "1 级报警 - 数据平台" 
      },
      "template": "red" 
    },
    "elements": [
				{
			  "fields": [
				{
				  "is_short": true,
				  "text": {
					"content": "**🕐 时间：**\ntimeSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": true,
				  "text": {
					"content": "**👤 值班：**\nuserSlot",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": false,
				  "text": {
					"content": "",
					"tag": "lark_md"
				  }
				},
				{
				  "is_short": false,
				  "text": {
					"content": "**报警信息:**\nmsgSlot",
					"tag": "lark_md"
				  }
				}
			  ],
			  "tag": "div"
			},
			{
			  "actions": [
				{
				  "tag": "button",
				  "text": {
					"content": "跟进处理",
					"tag": "plain_text"
				  },
				  "type": "primary",
				  "value": {
					"key1": "value1"
				  }
				}
			  ],
			  "tag": "action"
			},
			{
			  "tag": "hr"
			}
    ]
  }
}`
